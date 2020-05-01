package game

// MapElement describes which kind of object are inside a specific position
type MapElement int

const (
	// MapElementNone nothing on this point of the map
	MapElementNone MapElement = iota
	// MapElementWall identifies when there is a wall on this point of the map
	MapElementWall
	// MapElementSpawn identifies when there is someone on this point of the map
	MapElementSpawn
)

// GetMapElements goes through the game map, and return a description of each
// map element and his position
func (e *Engine) GetMapElements() map[MapElement][]Point {
	center := e.getMapCenter()
	elements := make(map[MapElement][]Point, 0)
	for mapY, row := range e.gameMap {
		for mapX, col := range row {
			mapElement := MapElementNone
			switch col {
			case '█':
				mapElement = MapElementWall
			case 'S':
				mapElement = MapElementSpawn
			}
			elements[mapElement] = append(elements[mapElement], Point{
				X: mapX - center.X,
				Y: mapY - center.Y,
			})
		}
	}
	return elements
}

// getMapDimensions will get the dimensions of the current map, in the form
// width + height
func (e *Engine) getMapDimensions() (int, int) {
	if len(e.gameMap) == 0 {
		return 0, 0
	}
	return len(e.gameMap[0]), len(e.gameMap)
}

func (e *Engine) getMapCenter() (c Point) {
	width, height := e.getMapDimensions()
	return Point{
		X: width / 2,
		Y: height / 2,
	}
}
