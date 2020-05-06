package game_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/ramonmacias/go-spaceship-shooter/internal/game"
	"github.com/stretchr/testify/assert"
)

var mapTest = [][]rune{
	{'█', '█', '█', '█', '█', '█', '█', '█'},
	{'█', ' ', ' ', ' ', ' ', ' ', ' ', '█'},
	{'█', ' ', ' ', ' ', ' ', '█', ' ', '█'},
	{'█', ' ', ' ', 'S', ' ', ' ', ' ', '█'},
	{'█', '█', ' ', ' ', ' ', ' ', ' ', '█'},
	{'█', ' ', ' ', ' ', '█', ' ', ' ', '█'},
	{'█', ' ', ' ', ' ', ' ', ' ', 'S', '█'},
	{'█', '█', '█', '█', '█', '█', '█', '█'},
}

func TestMoveActionPerform(t *testing.T) {
	// Setup basics scenario
	var err error
	actors := make(map[uuid.UUID]game.Actor)
	actor := game.Actor{
		Name: "TestActor",
		// Initial postion is on the center of the map
		Position: game.Point{
			X: 0,
			Y: 0,
		},
	}
	actor.ID, err = uuid.NewV4()
	assert.Nil(t, err)
	actors[actor.ID] = actor
	e := &game.Engine{
		GameMap: mapTest,
		Actors:  actors,
	}

	movements := []struct {
		name     string
		action   game.MoveAction
		expected game.Point
	}{
		{
			name: "Should move up",
			action: game.MoveAction{
				ActorID:   actor.ID,
				Direction: game.DirectionUp,
				CreatedAt: time.Now(),
			},
			expected: game.Point{
				X: 0,
				Y: -1,
			},
		},
		{
			name: "Should move left",
			action: game.MoveAction{
				ActorID:   actor.ID,
				Direction: game.DirectionLeft,
				CreatedAt: time.Now(),
			},
			expected: game.Point{
				X: -1,
				Y: -1,
			},
		},
		{
			name: "Should move down",
			action: game.MoveAction{
				ActorID:   actor.ID,
				Direction: game.DirectionDown,
				CreatedAt: time.Now(),
			},
			expected: game.Point{
				X: -1,
				Y: 0,
			},
		},
		{
			name: "Should move right",
			action: game.MoveAction{
				ActorID:   actor.ID,
				Direction: game.DirectionRight,
				CreatedAt: time.Now(),
			},
			expected: game.Point{
				X: 0,
				Y: 0,
			},
		},
		{
			name: "Shouldn't be able to move down",
			action: game.MoveAction{
				ActorID:   actor.ID,
				Direction: game.DirectionDown,
				CreatedAt: time.Now(),
			},
			expected: game.Point{
				X: 0,
				Y: 0,
			},
		},
	}
	for _, tt := range movements {
		tt.action.Perform(e)
		if !assert.True(t, tt.expected.Equal(e.Actors[actor.ID].Position)) {
			t.Error("Failure on the step", tt.name, tt.expected, e.Actors[actor.ID].Position)
		}
	}
}
