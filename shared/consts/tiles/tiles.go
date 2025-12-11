package tiles

import (
	types "sweep/shared/types"
)
const (
	ClosedMine types.Tile = iota
	FlaggedMine
	OpenMine
	ClosedSafe
	OpenSafe
	FlaggedSafe
	OutOfBounds
)