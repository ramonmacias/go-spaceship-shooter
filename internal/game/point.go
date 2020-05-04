package game

// Point represents a specific position in a mapp
type Point struct {
	X int
	Y int
}

// Equal returns wether the positions represented by p and p2 are equal
func (p Point) Equal(p2 Point) bool {
	return p.X == p2.X && p.Y == p2.Y
}
