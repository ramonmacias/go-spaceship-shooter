package game

import (
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
	DirectionUp Direction = iota
	DirectionDown
	DirectionLeft
	DirectionRight
	DirectionStop
)

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
		actor.Position.X--
	case DirectionLeft:
		actor.Position.X++
	}
	// Check if we collide with a wall
	if e.GameMap.IsWall(actor.Position) {
		return
	}
	e.Actors[m.ActorID] = actor
}
