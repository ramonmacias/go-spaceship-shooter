package game

import "github.com/gofrs/uuid"

// Bot represents the basic information needed for handle all the AI enemies on
// the game
type Bot struct {
	ID       uuid.UUID
	Life     int
	Position Point
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
			}
		}
		return nil
	}
}
