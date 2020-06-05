package game

import "github.com/gofrs/uuid"

// Actor defines all the different entities that has the feature of change the
// behaviour of the game status
// TODO pending to move this as an interface when we have bots
type Actor struct {
	ID       uuid.UUID
	Name     string
	Position Point
	Life     int
}

// SetActors will attach the given actor to the game engine
func SetActors(actors map[uuid.UUID]Actor) engineOpt {
	return func(e *Engine) error {
		e.Actors = actors
		return nil
	}
}
