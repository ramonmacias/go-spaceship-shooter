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
	mapTest3 = [][]rune{
		{'█', '█', '█'},
		{'█', 'S', '█'},
		{'█', ' ', '█'},
		{'█', '█', '█'},
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

func TestSizesGetMapElementsMethod(t *testing.T) {
	tests := []struct {
		name     string
		gameMap  Map
		expected map[MapElement]int
	}{
		{
			name:    "Should return and empty elements on a [empty map]",
			gameMap: [][]rune{},
			expected: map[MapElement]int{
				MapElementNone:  0,
				MapElementWall:  0,
				MapElementSpawn: 0,
			},
		},
		{
			name:    "Should return a specific number of elements on a [square map]",
			gameMap: mapTest1,
			expected: map[MapElement]int{
				MapElementWall:  31,
				MapElementNone:  31,
				MapElementSpawn: 2,
			},
		},
		{
			name:    "Should return a specific number of elements on a [rectangle map]",
			gameMap: mapTest2,
			expected: map[MapElement]int{
				MapElementWall:  23,
				MapElementNone:  16,
				MapElementSpawn: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements := tt.gameMap.GetMapElements()
			assert.Equal(t, tt.expected[MapElementNone], len(elements[MapElementNone]))
			assert.Equal(t, tt.expected[MapElementWall], len(elements[MapElementWall]))
			assert.Equal(t, tt.expected[MapElementSpawn], len(elements[MapElementSpawn]))
		})
	}
}

func TestIsWallMethod(t *testing.T) {
	tests := []struct {
		name     string
		gameMap  Map
		position Point
		expected bool
	}{
		{
			name:    "Shouldn't find a wall on this position",
			gameMap: mapTest3,
			position: Point{
				X: 0,
				Y: 0,
			},
			expected: false,
		},
		{
			name:    "Should find a wall on this position",
			gameMap: mapTest3,
			position: Point{
				X: -1,
				Y: -1,
			},
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.gameMap.IsWall(tt.position))
		})
	}
}
