package game

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mapTest = [][]rune{
		{'█', '█', '█', '█', '█', '█', '█', '█'},
		{'█', ' ', ' ', ' ', ' ', ' ', ' ', '█'},
		{'█', ' ', ' ', ' ', ' ', '█', ' ', '█'},
		{'█', ' ', ' ', 'S', ' ', ' ', ' ', '█'},
		{'█', '█', ' ', ' ', ' ', ' ', ' ', '█'},
		{'█', ' ', ' ', ' ', '█', ' ', ' ', '█'},
		{'█', ' ', ' ', ' ', ' ', ' ', 'S', '█'},
		{'█', '█', '█', '█', '█', '█', '█', '█'},
	}
)

func TestSetBots(t *testing.T) {
	tests := []struct {
		name string
		args struct {
			strategies []BotStrategy
			engine     *Engine
		}
		expected error
	}{
		{
			name: "Should fail with non matching bots sizing",
			args: struct {
				strategies []BotStrategy
				engine     *Engine
			}{
				strategies: []BotStrategy{OnlyMovementStrategy},
				engine: &Engine{
					GameMap: mapTest,
				},
			},
			expected: fmt.Errorf("Expected 2 bots but received 1"),
		},
		{
			name: "Should not fail",
			args: struct {
				strategies []BotStrategy
				engine     *Engine
			}{
				strategies: []BotStrategy{OnlyMovementStrategy, OnlyMovementStrategy},
				engine: &Engine{
					GameMap: mapTest,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, SetBots(tt.args.strategies)(tt.args.engine))
		})
	}
}
