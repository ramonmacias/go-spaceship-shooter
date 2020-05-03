package game

// Map is a matrix representation of game map
type Map [][]rune

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
// map element and his position, seems we are doing distances and not coordinates
// TODO check this coordinat/position
func (m Map) GetMapElements() map[MapElement][]Point {
	center := m.getMapCenter()
	elements := make(map[MapElement][]Point, 0)
	for mapY, row := range m {
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

// IsWall will check if on the given position exists a wall
func (m Map) IsWall(p Point) bool {
	for _, position := range m.GetMapElements()[MapElementWall] {
		if position == p {
			return true
		}
	}
	return false
}

// getMapDimensions will get the dimensions of the current map, in the form
// width + height
func (m Map) getMapDimensions() (int, int) {
	if len(m) == 0 {
		return 0, 0
	}
	return len(m[0]), len(m)
}

func (m Map) getMapCenter() (c Point) {
	width, height := m.getMapDimensions()
	return Point{
		X: width / 2,
		Y: height / 2,
	}
}
