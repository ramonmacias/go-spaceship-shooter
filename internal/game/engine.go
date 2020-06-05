package game

import (
	"log"
	"sync"

	"github.com/gofrs/uuid"
)

// engineOpt is a function used while the engine creation in order to setup
// all the basics for start the game
type engineOpt func(e *Engine) error

// Engine type will keep all the main information related with the game
type Engine struct {
	// Actors keep the information and link about all the interactors of the
	// game
	Actors map[uuid.UUID]Actor
	// GameMap keep link to the current map is playing
	GameMap Map
	// ActionChan is a buffered channel used for comunication between view and the
	// engine
	ActionChan chan Action
	// Score keep the info related with points and actors
	Score map[uuid.UUID]int
	// RoundWinner keep the id for the winner
	RoundWinner uuid.UUID
	// LevelComplete is the flag that determines when the level is complete
	LevelComplete bool
	// GameOver is the flag that determines when the player dies
	GameOver bool
	// Lasers keep the information about each lasers on the map
	Lasers sync.Map
	// Bots keep the information about bots on the map
	Bots sync.Map
}

// NewEngine function will build a new engine with the applied engine options
func NewEngine(opts ...engineOpt) *Engine {
	e := &Engine{
		ActionChan: make(chan Action, 100),
		Score:      make(map[uuid.UUID]int),
	}
	for _, fn := range opts {
		if err := fn(e); err != nil {
			log.Panic("Error while applying engineOpt", err)
		}
	}
	return e
}

// Start will setup the basics for run the game
func (e *Engine) Start() {
	go e.actionsListener()
	e.startBots()
}

// actionsListener will be listening for all the events received from the action
// channel and will apply them
func (e *Engine) actionsListener() {
	for {
		action := <-e.ActionChan
		action.Perform(e)
	}
}

// updateScores will update the global scores.
// TODO refactor this, now we only have one player, but in case we have multiple
// we will need to receive at least the actorID
func (e *Engine) updateScores() {
	for actorID := range e.Actors {
		e.Score[actorID] += 10
	}
}
