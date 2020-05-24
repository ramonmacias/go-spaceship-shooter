package game

import (
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

// SetBots will go through the game map and will create a new bot on each
// spawn position
func SetBots() engineOpt {
	return func(e *Engine) error {
		for index, spawnPosition := range e.GameMap.GetMapElements()[MapElementSpawn] {
			botID := uuid.Must(uuid.NewV4())
			if index == 0 {
				e.Bots.Store(botID, Bot{
					ID:       botID,
					Life:     4,
					Position: spawnPosition,
					Strategy: OnlyShootingStrategy,
				})
			} else {
				e.Bots.Store(botID, Bot{
					ID:       botID,
					Life:     4,
					Position: spawnPosition,
					Strategy: OnlyMovementStrategy,
				})
			}
		}
		return nil
	}
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
	}
}
