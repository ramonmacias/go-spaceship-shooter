package game_test

import (
	"testing"

	"github.com/ramonmacias/go-spaceship-shooter/internal/game"
	"github.com/stretchr/testify/assert"
)

func TestPointEqualMethod(t *testing.T) {
	tests := []struct {
		name        string
		firstPoint  game.Point
		secondPoint game.Point
		expected    bool
	}{
		{
			name: "Should be equals",
			firstPoint: game.Point{
				X: 1,
				Y: 1,
			},
			secondPoint: game.Point{
				X: 1,
				Y: 1,
			},
			expected: true,
		},
		{
			name: "Shouldn't be equals",
			firstPoint: game.Point{
				X: 1,
				Y: 1,
			},
			secondPoint: game.Point{
				X: 0,
				Y: 3,
			},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.firstPoint.Equal(tt.secondPoint))
		})
	}
}
