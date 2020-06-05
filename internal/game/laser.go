package game

import "github.com/gofrs/uuid"

// Origin is used to define from where the action is coming from
type Origin int

const (
	// OriginUnknown define a non known origin
	OriginUnknown Origin = iota
	// OriginPlayer is applied when the owner of the action is a player
	OriginPlayer
	// OriginBot is applied when the own of the action is a bot
	OriginBot
)

// Laser defines a laser shoot
type Laser struct {
	ID       uuid.UUID
	Position Point
	Origin   Origin
}

// checkCollisions will check if there is some bot on the given laser position
// if it is will reduce the life and remove the bot if needed and return true,
// otherwise do nothing and return false
func (e *Engine) checkLaserCollisions(laserPosition Point, origin Origin) (collide bool) {
	// TODO refacor this, each actor should has his own check collider
	switch origin {
	case OriginPlayer:
		botToDelete := uuid.Nil
		e.Bots.Range(func(key interface{}, value interface{}) bool {
			bot := value.(Bot)
			botID := key.(uuid.UUID)
			if bot.Position.Equal(laserPosition) {
				bot.Life--
				e.Bots.Store(botID, bot)
				if bot.Life == 0 {
					botToDelete = botID
				}
				collide = true
				return false
			}
			return true
		})
		// We can't remove a key value on a ranging map
		if botToDelete != uuid.Nil {
			e.Bots.Delete(botToDelete)
			e.LevelComplete = e.botCount() == 0
		}
	case OriginBot:
		for actorID, actor := range e.Actors {
			if actor.Position.Equal(laserPosition) {
				collide = true
				actor.Life--
				e.Actors[actorID] = actor
				if actor.Life <= 0 && !e.GameOver {
					e.GameOver = true
				}
			}
		}
	}
	return collide
}
