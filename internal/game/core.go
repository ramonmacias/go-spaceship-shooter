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
	ActionChan      chan Action
	lastAction      map[string]time.Time
	Score           map[uuid.UUID]int
	NewRoundAt      time.Time
	RoundWinner     uuid.UUID
	WaitForRound    bool
	IsAuthoritative bool
	spawnPointIndex int
	// Lasers keep the information about each lasers on the map
	Lasers sync.Map
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

// Laser defines a laser shoot
type Laser struct {
	ID       uuid.UUID
	Position Point
}
