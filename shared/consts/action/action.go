package actions

import (
	"fmt"
	"strconv"
	"strings"
)

type ActionType string

const (
	MoveCursorUp    ActionType = "move cursor up"
	MoveCursorDown  ActionType = "move cursor down"
	MoveCursorLeft  ActionType = "move cursor left"
	MoveCursorRight ActionType = "move cursor right"
	OpenTile        ActionType = "open tile"
	FlagTile        ActionType = "flag tile"

	MoveCursorToTopRow      ActionType = "move cursor to top row"
	MoveCursorToBottomRow   ActionType = "move cursor to bottom row"
	MoveCursorToFirstColumn ActionType = "move cursor to first column"
	MoveCursorToLastColumn  ActionType = "move cursor to last column"
)

var bindingsMap map[string]ActionType = map[string]ActionType{}

func (a ActionType) SetBinding(binding string) {
	bindingsMap[binding] = a
}

func IsAction(str string) bool {
	switch ActionType(str) {
	case MoveCursorDown, MoveCursorLeft,
		MoveCursorRight, MoveCursorUp,
		OpenTile, FlagTile,
		MoveCursorToBottomRow, MoveCursorToFirstColumn,
		MoveCursorToLastColumn, MoveCursorToTopRow:
		return true
	default:
		return false
	}
}

type Action struct {
	Kind       ActionType
	Multiplier byte
}

type MultiplierParseError struct {
	multiplier   string
	strconvError error
}

func (e *MultiplierParseError) Error() string {
	return fmt.Sprintf("could not parse multiplier \"%v\": %v", e.multiplier, e.Error())
}
func (e *MultiplierParseError) Is(target error) bool {
	return e.Error() == target.Error()
}

type InvalidBindError struct {
	bind string
}

func (e *InvalidBindError) Error() string {
	return fmt.Sprintf("no matching action for bind \"%v\"", e.bind)
}
func (e *InvalidBindError) Is(target error) bool {
	return e.Error() == target.Error()
}

var errorCount int = 0

func GetAction(keyStrokes string) (*Action, error) {
	var firstKeyIx uint

	symbols := strings.Split(keyStrokes, "")
loop:
	for ix, symbol := range symbols {
		switch symbol {
		case "0", "1", "2",
			"3", "4", "5",
			"6", "7", "8", "9":
			continue
		default:
			firstKeyIx = uint(ix)
			break loop
		}
	}

	keys := strings.Join(symbols[firstKeyIx:], "")
	kind, ok := bindingsMap[keys]
	if !ok {
		return nil, &InvalidBindError{keys}
	}

	multiplier := byte(1)
	var nums string
	if firstKeyIx != 0 {
		nums = strings.Join(symbols[:firstKeyIx], "")
	}
	if len(nums) != 0 {
		multiplier64, err := strconv.ParseUint(nums, 10, 64)
		if err != nil {
			return nil, &MultiplierParseError{nums, err}
		}
		multiplier = byte(multiplier64)
	}

	return &Action{
		Kind:       kind,
		Multiplier: multiplier,
	}, nil
}
