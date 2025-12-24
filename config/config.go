package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"

	flags "sweep/config/flags"
	tilecontent "sweep/shared/consts/tile-content"
	types "sweep/shared/types"
	glyphs "sweep/shared/vars/glyphs"
	"sweep/shared/vars/regexes"
	styles "sweep/tui/styles"
	themepreview "sweep/tui/theme-preview"

	gojsonschema "github.com/xeipuuv/gojsonschema"
)

type Defaults struct {
	Width  uint16
	Height uint16
	Mines  uint16
}

type Color string

func (c Color) isSet() bool {
	return c != ""
}

type Colors map[string]string

type Action string

const (
	MOVE_CURSOR_UP    Action = "moveCursorUp"
	MOVE_CURSOR_DOWN  Action = "moveCursorDown"
	MOVE_CURSOR_LEFT  Action = "moveCursorLeft"
	MOVE_CURSOR_RIGHT Action = "moveCursorRight"
	OPEN_TILE         Action = "openTile"
	FLAG_TILE         Action = "flagTile"
)

type Bindings struct {
	MoveCursorUp    []string `json:"moveCursorUp"`
	MoveCursorDown  []string `json:"moveCursorDown"`
	MoveCursorLeft  []string `json:"moveCursorLeft"`
	MoveCursorRight []string `json:"moveCursorRight"`
	OpenTile        []string `json:"openTile"`
	FlagTile        []string `json:"flagTile"`
}

func (b *Bindings) IsMoveCursorUp(str string) bool {
	return slices.Contains(b.MoveCursorUp, str)
}
func (b *Bindings) IsMoveCursorDown(str string) bool {
	return slices.Contains(b.MoveCursorDown, str)
}
func (b *Bindings) IsMoveCursorLeft(str string) bool {
	return slices.Contains(b.MoveCursorLeft, str)
}
func (b *Bindings) IsMoveCursorRight(str string) bool {
	return slices.Contains(b.MoveCursorRight, str)
}
func (b *Bindings) IsOpenTile(str string) bool {
	return slices.Contains(b.OpenTile, str)
}
func (b *Bindings) IsFlagTile(str string) bool {
	return slices.Contains(b.FlagTile, str)
}

type Cursor struct {
	Color     Color  `json:"color"`
	LeftHalf  string `json:"left half"`
	RightHalf string `json:"right half"`
}

func (c *Cursor) validate() (bool, []string) {
	errors := []string{}
	if c.Color != "" && !regexes.ColorRegex.MatchString(string(c.Color)) {
		errors = append(errors, "(cursor.color) cursor color does not match ANSI nor HEX RGB")
	}
	if len(c.LeftHalf) > 1 {
		errors = append(errors, "(cursor.left half) cursor left half is longer than one character")
	}
	if len(c.RightHalf) > 1 {
		errors = append(errors, "(cursor.right half) cursor right half is longer than one character")
	}

	return len(errors) == 0, errors
}

type Config struct {
	Flags    []types.Flag `json:"flags"`
	Defaults Defaults     `json:"defaults"`
	Colors   Colors       `json:"colors"`
	Bindings Bindings     `json:"bindings"`

	Mines  uint16 `json:"mines,omitempty"`
	Width  uint16 `json:"width,omitempty"`
	Height uint16 `json:"height,omitempty"`

	Cursor Cursor `json:"cursor"`
}

func (c *Colors) validate() (bool, []string) {
	errors := []string{}
	for key, val := range *c {
		if val != "" && !regexes.ColorRegex.MatchString(val) {
			errors = append(errors, fmt.Sprintf("(colors.%v) %v does not match ANSI nor HEX RGB", key, val))
		}
		if _, err := tilecontent.FromString(key); err != nil {
			errors = append(errors, fmt.Sprintf("(colors) %v is not a valid option", key))
		}
	}
	return len(errors) == 0, errors
}

func (config *Config) applyColors() {
	for key, color := range config.Colors {
		tileContent, _ := tilecontent.FromString(key)
		styles.SetTileColor(tileContent, color)
	}
}

var (
	schema = gojsonschema.NewReferenceLoader("file:///home/erokez/Desktop/code/sweep/config.schema.json")
)

func (config *Config) validate() (bool, []string) {
	configLoader := gojsonschema.NewGoLoader(config)
	errors := []string{}

	result, err := gojsonschema.Validate(schema, configLoader)
	if err != nil {
		log.Fatalf("%v\nProbable cause: config file does not exist", err.Error())
	}

	if !result.Valid() {
		schemaErrors := make([]string, len(result.Errors()))
		for ix, error := range result.Errors() {
			schemaErrors[ix] = error.String()
		}
		errors = append(errors, schemaErrors...)
	}

	if isValid, colorsErrors := config.Colors.validate(); !isValid {
		errors = append(errors, colorsErrors...)
	}

	if isValid, cursorErrors := config.Cursor.validate(); !isValid {
		errors = append(errors, cursorErrors...)
	}

	return len(errors) == 0, errors
}

func getFlagIntArgument(args []string, index int) (uint16, error) {
	val := args[index+1]
	if len(args) <= index {
		return 0, fmt.Errorf("argument was not provided")
	}
	numVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("argument is not an integer")
	}
	return uint16(numVal), nil
}

func (c *Config) applyFlags() {
	skip := false
	args := os.Args[1:]
	flagList := append(args, c.Flags...)

	preview := false

	for ix, arg := range flagList {
		if skip {
			skip = false
			continue
		}

		switch arg {
		case flags.THEME_PREVIEW, flags.THEME_PREVIEW_SHORT:
			preview = true

		case flags.CONFIG, flags.CONFIG_SHORT:
			dir, _ := os.Getwd()
			fmt.Printf("%v/config.json\n", dir)
			os.Exit(0)

		case flags.HEIGHT, flags.HEIGHT_SHORT:
			skip = true
			val, err := getFlagIntArgument(args, ix)
			if err != nil {
				fmt.Printf("%v - %v", arg, err)
			} else {
				c.Height = val
			}

		case flags.WIDTH, flags.WIDTH_SHORT:
			skip = true
			val, err := getFlagIntArgument(args, ix)
			if err != nil {
				fmt.Printf("%v - %v", arg, err)
			} else {
				c.Width = val
			}

		case flags.MINES, flags.MINES_SHORT:
			skip = true
			val, err := getFlagIntArgument(args, ix)
			if err != nil {
				fmt.Printf("%v - %v", arg, err)
			} else {
				c.Mines = val
			}

		case flags.ASCII, flags.ASCII_SHORT:
			glyphs.MINE = "M"

			glyphs.FLAG = "F"

			glyphs.WRONG_FLAG = "W"

		case flags.FILL, flags.FILL_SHORT:
			styles.SetFill(true)

		}
	}
	if preview {
		c.applyColors()
		print(themepreview.RenderThemePreview())
		os.Exit(0)
	}
}

func (config *Config) applyCursorStyle() {
	if config.Cursor.Color.isSet() {
		styles.SetCursorColor(string(config.Cursor.Color))
	}
	if config.Cursor.LeftHalf != "" {
		glyphs.CursorLeftHalf = config.Cursor.LeftHalf
	}
	if config.Cursor.RightHalf != "" {
		glyphs.CursorRightHalf = config.Cursor.RightHalf
	}
}

func LoadConfig(configPath string) (*Config, error) {
	configBin, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Could not read config %v\nDoes the file exist?", configPath)
	}

	config := new(Config)
	err = json.Unmarshal(configBin, config)
	if err != nil {
		log.Fatalf("Could not parse config %v\nCheck its integriy", configPath)
	}

	isValid, errors := config.validate()

	if !isValid {
		fmt.Println("Your config file has errors")
		for k, v := range errors {
			fmt.Printf("%v. %v", k+1, v)
		}
		os.Exit(1)
	}

	config.applyFlags()
	config.applyColors()
	config.applyCursorStyle()

	return config, nil
}

func GetConfig() *Config {
	config := new(Config)

	var err error
	config, err = LoadConfig("/home/erokez/Desktop/code/sweep/config.json")
	if err == nil {
		return config
	}

	return config
}
