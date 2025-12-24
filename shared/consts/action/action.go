package actions

type Action string

const (
	MoveCursorUp    Action = "move cursor up"
	MoveCursorDown  Action = "move cursor down"
	MoveCursorLeft  Action = "move cursor left"
	MoveCursorRight Action = "move cursor right"
	OpenTile        Action = "open tile"
	FlagTile        Action = "flag tile"
)

func IsAction(str string) bool {
	switch Action(str) {
	case MoveCursorDown, MoveCursorLeft,
			MoveCursorRight, MoveCursorUp,
			OpenTile, FlagTile:
		return true
	default:
		return false
	}
}
