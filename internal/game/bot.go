package game

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

// BotStrategy holds the specific strategy in terms of movement and shootin that
// the bot is going to follow
type BotStrategy int

const (
	// NoMovementStrategy defines the bot will be just stopped
	NoMovementStrategy BotStrategy = iota
	// OnlyMovementStrategy defines the bot will be in movement but no shooting
	OnlyMovementStrategy
	// OnlyShootingStrategy defines the bot will be shooting all the time but without
	// movement
	OnlyShootingStrategy
	// ShootAndMoveStrategy define the bot will be in movement and shooting
	ShootAndMoveStrategy
)

// Bot represents the basic information needed for handle all the AI enemies on
// the game
type Bot struct {
	ID       uuid.UUID
	Life     int
	Position Point
	Strategy BotStrategy
}

// SetBots will receive an slice of bot strategies, this slice should match in
// size with the expected numbers of spawn positions
func SetBots(strategies []BotStrategy) engineOpt {
	return func(e *Engine) error {
		spawnElements := e.GameMap.GetMapElements()[MapElementSpawn]
		if len(strategies) != len(spawnElements) {
			return fmt.Errorf("Expected %d bots but received %d", len(spawnElements), len(strategies))
		}
		for index, spawnPosition := range spawnElements {
			botID := uuid.Must(uuid.NewV4())
			e.Bots.Store(botID, Bot{
				ID:       botID,
				Life:     4,
				Position: spawnPosition,
				Strategy: strategies[index],
			})
		}
		return nil
	}
}

// botCount is a workaround I should do because there is not implemented yet
// this len(sync.Map). For more info visit this issue https://github.com/golang/go/issues/20680
// maybe exposing an atomic int will be enough, but I don't have access to it
// this is why I need to range over the bots, of course is very inefficient
func (e *Engine) botCount() (length int) {
	e.Bots.Range(func(key interface{}, value interface{}) bool {
		length++
		return true
	})
	return length
}

// startBots will apply all the strategies linked to each bot
func (e *Engine) startBots() {
	e.Bots.Range(func(key interface{}, value interface{}) bool {
		bot := value.(Bot)
		go bot.Strategy.perform(e, bot)
		return true
	})
}

// perform will execute the behaviour linked to the given strategy
func (s BotStrategy) perform(e *Engine, bot Bot) {
	switch s {
	case OnlyMovementStrategy:
		ticker := time.NewTicker(200 * time.Millisecond)
		for {
			if _, exists := e.Bots.Load(bot.ID); !exists {
				return
			}
			select {
			case <-ticker.C:
				e.ActionChan <- &BotMoveAction{
					BotID:     bot.ID,
					Direction: RandomDirection(),
					CreatedAt: time.Now(),
				}
			default:
			}
		}
	case OnlyShootingStrategy:
		ticker := time.NewTicker(300 * time.Millisecond)
		for {
			b, exists := e.Bots.Load(bot.ID)
			if !exists {
				return
			}
			bot := b.(Bot)
			select {
			case <-ticker.C:
				laserID := uuid.Must(uuid.NewV4())
				e.Lasers.Store(laserID, Laser{
					ID:       laserID,
					Position: bot.Position,
					Origin:   OriginBot,
				})
				e.ActionChan <- &LaserAction{
					LaserID:   laserID,
					Direction: RandomDirection(),
					CreatedAt: time.Now(),
				}
			default:
			}
		}
	case ShootAndMoveStrategy:
		movementTicker := time.NewTicker(200 * time.Millisecond)
		shootingTicker := time.NewTicker(900 * time.Millisecond)
		for {
			b, exists := e.Bots.Load(bot.ID)
			if !exists {
				return
			}
			bot := b.(Bot)
			select {
			case <-movementTicker.C:
				e.ActionChan <- &BotMoveAction{
					BotID:     bot.ID,
					Direction: RandomDirection(),
					CreatedAt: time.Now(),
				}
			case <-shootingTicker.C:
				laserID := uuid.Must(uuid.NewV4())
				e.Lasers.Store(laserID, Laser{
					ID:       laserID,
					Position: bot.Position,
					Origin:   OriginBot,
				})
				e.ActionChan <- &LaserAction{
					LaserID:   laserID,
					Direction: RandomDirection(),
					CreatedAt: time.Now(),
				}
			default:
			}
		}
	}
}
