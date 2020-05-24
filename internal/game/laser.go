package game

import "github.com/gofrs/uuid"

// Origin is used to define from where the action is coming from
type Origin int

const (
	// OriginUnknown define a non known origin
	OriginUnknown Origin = iota
	// OriginPlayer is applied when the owner of the action is a player
	OriginPlayer
	// OriginBot is applied when the own of the action is a bot
	OriginBot
)

// Laser defines a laser shoot
type Laser struct {
	ID       uuid.UUID
	Position Point
	Origin   Origin
}
