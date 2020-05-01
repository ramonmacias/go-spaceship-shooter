package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var mapTest1 = [][]rune{
	{'█', '█', '█', '█', '█', '█', '█', '█'},
	{'█', ' ', ' ', ' ', ' ', ' ', ' ', '█'},
	{'█', ' ', ' ', ' ', ' ', '█', ' ', '█'},
	{'█', ' ', ' ', 'S', ' ', ' ', ' ', '█'},
	{'█', '█', ' ', ' ', ' ', ' ', ' ', '█'},
	{'█', ' ', ' ', ' ', '█', ' ', ' ', '█'},
	{'█', ' ', ' ', ' ', ' ', ' ', 'S', '█'},
	{'█', '█', '█', '█', '█', '█', '█', '█'},
}

func TestGetDimensionsMap(t *testing.T) {
	tests := []struct {
		name           string
		gameMap        [][]rune
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "Should return 0 and 0 for empty map",
			gameMap:        [][]rune{},
			expectedWidth:  0,
			expectedHeight: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Engine{
				gameMap: tt.gameMap,
			}
			width, height := e.getMapDimensions()
			assert.Equal(t, tt.expectedWidth, width)
			assert.Equal(t, tt.expectedHeight, height)
		})
	}
}
