package types

type Tile byte

const (
	ClosedBomb Tile = iota
	FlaggedBomb
	OpenBomb
	ClosedSafe
	OpenSafe
	FlaggedSafe
	OutOfBounds
)

type Position struct {
	X uint16
	Y uint16
}

func (p Position) GetCoords() (uint16, uint16) {
	return p.X, p.Y
}

type IGameEngine interface {
	FlagToggleTile(Position)
	OpenTile(Position)
	GetTile(Position) Tile
	IsFinished() bool
	GetField() [][]Tile
	SetFieldSize(uint16, uint16)
	SetBombCount(uint16)
	SetBombs(Position)
	CountNeighbouringBombs(Position) byte

	GetWidth() uint16
	GetHeight() uint16
}

type Flag = string

