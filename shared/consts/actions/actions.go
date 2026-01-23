package actions

import (
	"fmt"
	"math"
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
	Quantifier uint16
}

type QuantifierParseError struct {
	Quantifier string
}

func (e *QuantifierParseError) Error() string {
	return fmt.Sprintf("could not parse Quantifier \"%v\"", e.Quantifier)
}
func (e *QuantifierParseError) Is(target error) bool {
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

func getKeysFromKeyStrokes(keyStrokes string) string {
	lastDigitIx := -1

	symbols := strings.Split(keyStrokes, "")
loop:
	for ix, symbol := range symbols {
		switch symbol {
		case "0", "1", "2",
			"3", "4", "5",
			"6", "7", "8", "9":
			lastDigitIx = ix
		default:
			break loop
		}
	}

	return strings.Join(symbols[lastDigitIx+1:], "")
}

func getQuantifierFromKeyStrokes(keyStrokes string, keys string) (uint16, error) {
	quantifier := uint16(1)
	quantifierPart := strings.Replace(keyStrokes, keys, "", 1)

	if len(quantifierPart) > 6 {
		return 0, &QuantifierParseError{quantifierPart}
	}

	if quantifierPart != "" {
		quantifier64, err := strconv.ParseUint(quantifierPart, 10, 16)
		if err != nil || quantifier64 > math.MaxUint16 {
			return 0, &QuantifierParseError{quantifierPart}
		}
		quantifier = uint16(quantifier64)
	}

	return quantifier, nil
}

func AnyBindingStartWith(keyStrokes string) bool {
	for keyPress, actionType := range bindingsMap {
		if strings.HasPrefix(keyPress, keyStrokes) {
			return true
		}
		switch actionType {
		case MoveCursorDown, MoveCursorLeft,
			MoveCursorRight, MoveCursorUp,
			MoveCursorToBottomRow,
			MoveCursorToTopRow, MoveCursorToLastColumn:

			keys := getKeysFromKeyStrokes(keyStrokes)
			quantifier, err := getQuantifierFromKeyStrokes(keyStrokes, keys)
			if err != nil {
				return false
			}

			if strconv.FormatUint(uint64(quantifier), 10) == keys {
				return true
			}

			if keys == "" {
				return true
			}

			if strings.HasPrefix(keyPress, keys) {
				return true
			}
		}
	}

	return false
}

func GetAction(keyStrokes string) (*Action, error) {
	kind, ok := bindingsMap[keyStrokes]
	if ok {
		return &Action{
			Kind:       kind,
			Quantifier: 1,
		}, nil
	}

	keys := getKeysFromKeyStrokes(keyStrokes)
	kind, ok = bindingsMap[keys]
	if !ok {
		return nil, &InvalidBindError{keys}
	}
	quantifier, err := getQuantifierFromKeyStrokes(keyStrokes, keys)
	if err != nil {
		return nil, err
	}

	return &Action{
		Kind:       kind,
		Quantifier: quantifier,
	}, nil
}
