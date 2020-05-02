package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mapTest1 = [][]rune{
		{'█', '█', '█', '█', '█', '█', '█', '█'},
		{'█', ' ', ' ', ' ', ' ', ' ', ' ', '█'},
		{'█', ' ', ' ', ' ', ' ', '█', ' ', '█'},
		{'█', ' ', ' ', 'S', ' ', ' ', ' ', '█'},
		{'█', '█', ' ', ' ', ' ', ' ', ' ', '█'},
		{'█', ' ', ' ', ' ', '█', ' ', ' ', '█'},
		{'█', ' ', ' ', ' ', ' ', ' ', 'S', '█'},
		{'█', '█', '█', '█', '█', '█', '█', '█'},
	}
	mapTest2 = [][]rune{
		{'█', '█', '█', '█', '█', '█', '█', '█'},
		{'█', ' ', ' ', ' ', ' ', ' ', ' ', '█'},
		{'█', ' ', ' ', ' ', ' ', '█', ' ', '█'},
		{'█', ' ', ' ', 'S', ' ', ' ', ' ', '█'},
		{'█', '█', '█', '█', '█', '█', '█', '█'},
	}
)

func TestGetDimensionsMapMethod(t *testing.T) {
	tests := []struct {
		name           string
		gameMap        Map
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "Should return 0 and 0 for [empty map]",
			gameMap:        [][]rune{},
			expectedWidth:  0,
			expectedHeight: 0,
		},
		{
			name:           "Should return same height and width [square map]",
			gameMap:        mapTest1,
			expectedWidth:  8,
			expectedHeight: 8,
		},
		{
			name:           "Should return different height and widht [rectangle map]",
			gameMap:        mapTest2,
			expectedWidth:  8,
			expectedHeight: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width, height := tt.gameMap.getMapDimensions()
			assert.Equal(t, tt.expectedWidth, width)
			assert.Equal(t, tt.expectedHeight, height)
		})
	}
}

func TestGetMapCenterMethod(t *testing.T) {
	tests := []struct {
		name     string
		gameMap  Map
		expected Point
	}{
		{
			name:    "Should return 0 and 0 for an [empty map]",
			gameMap: [][]rune{},
			expected: Point{
				X: 0,
				Y: 0,
			},
		},
		{
			name:    "Should return same X and Y for a [square map]",
			gameMap: mapTest1,
			expected: Point{
				X: 4,
				Y: 4,
			},
		},
		{
			name:    "Should return different X and Y for a [rectangle map]",
			gameMap: mapTest2,
			expected: Point{
				X: 4,
				Y: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			center := tt.gameMap.getMapCenter()
			assert.Equal(t, tt.expected, center)
		})
	}
}
