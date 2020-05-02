package game

import (
	"sync"
	"time"

	"github.com/gofrs/uuid"
)

// Engine type will keep all the main information related with the game
type Engine struct {
	Entities        map[uuid.UUID]Identifier
	GameMap         Map
	Mu              sync.RWMutex
	ChangeChan      chan Change
	ActionChan      chan Action
	lastAction      map[string]time.Time
	Score           map[uuid.UUID]int
	NewRoundAt      time.Time
	RoundWinner     uuid.UUID
	WaitForRound    bool
	IsAuthoritative bool
	spawnPointIndex int
}

// Action interface should be implemented for each of the possible actions a
type Action interface {
	Perform(e Engine)
}

// Identifier is an entity that provides an ID method.
type Identifier interface {
	ID() uuid.UUID
}

// Change is sent by the game engine in response to Actions.
// TODO review this part
type Change interface{}

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
// as arrow keys pressed
type MoveAction struct {
	Direction Direction
	ID        uuid.UUID
	CreatedAt time.Time
}

// MoveChange this is a result change from an user action
type MoveChange struct {
	Change
	Entity    Identifier
	Direction Direction
	Position  Point
}

// AddEntityChange change for resulting behaviour in other entities
type AddEntityChange struct {
	Change
	Entity Identifier
}

// Point represents a specific position in a mapp
type Point struct {
	X int
	Y int
}
