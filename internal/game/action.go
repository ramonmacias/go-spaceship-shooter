package game

import (
	"math/rand"
	"time"

	"github.com/gofrs/uuid"
)

// Action interface should be implemented for each of the possible actions a
type Action interface {
	Perform(e *Engine)
}

// Direction is used to represent Direction constants.
type Direction int

// Contains direction constants - DirectionStop will take no effect.
const (
	DirectionNone Direction = iota
	DirectionUp
	DirectionDown
	DirectionLeft
	DirectionRight
)

// RandomDirection will get a random direction
func RandomDirection() Direction {
	rn := rand.Intn(5)
	return Direction(rn)
}

// MoveAction keep the information about the actions launched by the user, such
// as arrow keys pressed, an action is a definition of a movement applied to an
// entity on the map and each movement can have a specific direction
type MoveAction struct {
	ActorID   uuid.UUID
	Direction Direction
	CreatedAt time.Time
}

// Perform will execute all the behaviour associated to the action given
func (m *MoveAction) Perform(e *Engine) {
	actor := e.Actors[m.ActorID]
	switch m.Direction {
	case DirectionUp:
		actor.Position.Y--
	case DirectionDown:
		actor.Position.Y++
	case DirectionRight:
		actor.Position.X++
	case DirectionLeft:
		actor.Position.X--
	}
	// Check if we collide with a wall
	if e.GameMap.IsWall(actor.Position) {
		return
	}
	e.Actors[m.ActorID] = actor
}

// BotMoveAction defines the concept of movement for the bots
// TODO merge this with the MoveAction, the Actor can be an interface and share
// a common mover interface
type BotMoveAction struct {
	BotID     uuid.UUID
	Direction Direction
	CreatedAt time.Time
}

// Perform will execute the behaviour linked to a bot move action
func (m *BotMoveAction) Perform(e *Engine) {
	bot := e.Bots[m.BotID]
	switch m.Direction {
	case DirectionUp:
		bot.Position.Y--
	case DirectionDown:
		bot.Position.Y++
	case DirectionRight:
		bot.Position.X++
	case DirectionLeft:
		bot.Position.X--
	}
	// Check if we collide with a wall
	if e.GameMap.IsWall(bot.Position) {
		return
	}
	e.Bots[m.BotID] = bot
}

// LaserAction keep the information about all the lasers actioned by the player
// or players, and his position and also how long
type LaserAction struct {
	LaserID   uuid.UUID
	Direction Direction
	CreatedAt time.Time
}

// Perform will execute the specific behaviour for a laser action, which is
// move until it collide with something, once we collide the perform will
// take the specific reactions
func (l *LaserAction) Perform(e *Engine) {
	go func(la *LaserAction, en *Engine) {
		las, _ := en.Lasers.Load(la.LaserID)
		laser := las.(Laser)
		ticker := time.NewTicker(18 * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				switch la.Direction {
				case DirectionUp:
					laser.Position.Y--
				case DirectionDown:
					laser.Position.Y++
				case DirectionRight:
					laser.Position.X++
				case DirectionLeft:
					laser.Position.X--
				}
				// Check collisions with wall and also for other actors
				if e.GameMap.IsWall(laser.Position) {
					e.Lasers.Delete(l.LaserID)
					return
				}
				if e.checkLaserCollisions(laser.Position) {
					e.Lasers.Delete(l.LaserID)
					return
				}
				// Update position to be printed
				e.Lasers.Store(l.LaserID, laser)
			default:
			}
		}
	}(l, e)
}
