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
		e.Bots = make(map[uuid.UUID]Bot)
		for _, spawnPosition := range e.GameMap.GetMapElements()[MapElementSpawn] {
			botID := uuid.Must(uuid.NewV4())
			e.Bots[botID] = Bot{
				ID:       botID,
				Life:     4,
				Position: spawnPosition,
				Strategy: OnlyMovementStrategy,
			}
		}
		return nil
	}
}

// startBots will apply all the strategies linked to each bot
func (e *Engine) startBots() {
	for _, bot := range e.Bots {
		go bot.Strategy.perform(e, bot)
	}
}

// perform will execute the behaviour linked to the given strategy
func (s BotStrategy) perform(e *Engine, bot Bot) {
	switch s {
	case OnlyMovementStrategy:
		ticker := time.NewTicker(200 * time.Millisecond)
		for {
			e.Mu.Lock()
			if _, ok := e.Bots[bot.ID]; !ok {
				return
			}
			e.Mu.Unlock()
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
	}
}
