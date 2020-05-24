package game

import (
	"log"
	"sync"
	"time"

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
	// Bots keep the information about bots on the map
	Bots sync.Map
}

// NewEngine function will build a new engine with the applied engine options
func NewEngine(opts ...engineOpt) *Engine {
	e := &Engine{
		ActionChan: make(chan Action, 100),
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

// SetMap will attach the given map to the game engine
func SetMap(m Map) engineOpt {
	return func(e *Engine) error {
		e.GameMap = m
		return nil
	}
}

// SetActors will attach the given actor to the game engine
func SetActors(actors map[uuid.UUID]Actor) engineOpt {
	return func(e *Engine) error {
		e.Actors = actors
		return nil
	}
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

// checkCollisions will check if there is some bot on the given laser position
// if it is will reduce the life and remove the bot if needed and return true,
// otherwise do nothing and return false
func (e *Engine) checkLaserCollisions(laserPosition Point, origin Origin) (collide bool) {
	// TODO refacor this, each actor should his own check collider
	if origin != OriginBot {
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
		}
	}
	return collide
}

// Actor defines all the different entities that has the feature of change the
// behaviour of the game status
// TODO pending to move this as an interface when we have bots
type Actor struct {
	ID       uuid.UUID
	Name     string
	Position Point
}
