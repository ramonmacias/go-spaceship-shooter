package game

import (
	"sync"
	"time"

	"github.com/gofrs/uuid"
)

// Engine type will keep all the main information related with the game
type Engine struct {
	// Actors keep the information and link about all the interactors of the
	// game
	Actors map[uuid.UUID]Actor
	// GameMap keep link to the current map is playing
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

// Start will setup the basics for run the game
func (e *Engine) Start() {
	go e.actionsListener()
}

// actionsListener will be listening for all the events received from the action
// channel and will apply them
func (e *Engine) actionsListener() {
	for {
		action := <-e.ActionChan
		e.Mu.Lock()
		action.Perform(e)
		e.Mu.Unlock()
	}
}

// Actor defines all the different entities that has the feature of change the
// behaviour of the game status
// TODO pending to move this as an interface when we have bots
type Actor struct {
	ID       uuid.UUID
	Name     string
	Position Point
}

// Identifier is an entity that provides an ID method.
type Identifier interface {
	ID() uuid.UUID
}

// Change is sent by the game engine in response to Actions.
// TODO review this part
type Change interface{}

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
