package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"

	flags "sweep/config/flags"
	glyphs "sweep/shared/vars/glyphs"
	"sweep/tui/styles"
	themepreview "sweep/tui/theme-preview"
	types "sweep/shared/types"

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

type Colors struct {
	Zero  Color `json:"0"`
	One   Color `json:"1"`
	Two   Color `json:"2"`
	Three Color `json:"3"`
	Four  Color `json:"4"`
	Five  Color `json:"5"`
	Six   Color `json:"6"`
	Seven Color `json:"7"`
	Eight Color `json:"8"`
	Mine Color `json:"mine"`
	Flag Color `json:"flag"`
	WrongFlag Color `json:"wrong flag"`
	Empty Color `json:"empty"`
}

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

type Config struct {
	Flags    []types.Flag `json:"flags"`
	Defaults Defaults     `json:"defaults"`
	Colors   Colors       `json:"colors"`
	Bindings Bindings     `json:"bindings"`

	Mines  uint16 `json:"mines,omitempty"`
	Width  uint16 `json:"width,omitempty"`
	Height uint16 `json:"height,omitempty"`
}

func (config *Config) setColors() {
	if config.Colors == *new(Colors) {
		return
	}
	if config.Colors.Zero.isSet() {
		styles.SetColor("0", string(config.Colors.Zero))
	}
	if config.Colors.One.isSet() {
		styles.SetColor("1", string(config.Colors.One))
	}
	if config.Colors.Two.isSet() {
		styles.SetColor("2", string(config.Colors.Two))
	}
	if config.Colors.Three.isSet() {
		styles.SetColor("3", string(config.Colors.Three))
	}
	if config.Colors.Four.isSet() {
		styles.SetColor("4", string(config.Colors.Four))
	}
	if config.Colors.Five.isSet() {
		styles.SetColor("5", string(config.Colors.Five))
	}
	if config.Colors.Six.isSet() {
		styles.SetColor("6", string(config.Colors.Six))
	}
	if config.Colors.Seven.isSet() {
		styles.SetColor("7", string(config.Colors.Seven))
	}
	if config.Colors.Eight.isSet() {
		styles.SetColor("8", string(config.Colors.Eight))
	}

	if config.Colors.Mine.isSet() {
		styles.SetColor("mine", string(config.Colors.Mine))
	}
	if config.Colors.Flag.isSet() {
		styles.SetColor("flag", string(config.Colors.Flag))
	}
	if config.Colors.WrongFlag.isSet() {
		styles.SetColor("wrong flag", string(config.Colors.WrongFlag))
	}
	if config.Colors.Empty.isSet() {
		styles.SetColor("empty", string(config.Colors.Empty))
	}
}

var (
	schema = gojsonschema.NewReferenceLoader("file:///home/erokez/Desktop/code/sweep/config.schema.json")
)

func validate(source string) (bool, []error) {
	config := gojsonschema.NewReferenceLoader(fmt.Sprintf("file://%v", source))

	result, err := gojsonschema.Validate(schema, config)
	if err != nil {
		log.Fatal(err)
	}
	if result.Valid() {
		return true, nil
	}
	fmt.Println(result.Errors())
	return false, nil
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

func (c *Config) setFlags() {
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
		c.setColors()
		print(themepreview.RenderThemePreview())
		os.Exit(0)
	}
}

func LoadConfig(configPath string) (*Config, error) {
	configBin, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := new(Config)
	err = json.Unmarshal(configBin, config)
	if err != nil {
		return nil, err
	}

	isValid, errors := validate(configPath)

	if !isValid {
		fmt.Println("Your config file has errors")
		for k, v := range errors {
			fmt.Printf("%v - %v", k, v)
		}
	}

	config.setFlags()
	config.setColors()
	return config, nil
}

func GetConfig() *Config {
	config := new(Config)

	var err error
	config, err = LoadConfig("/home/erokez/Desktop/code/sweep/config.json")
	if err == nil {
		return config
	}

	config, err = LoadConfig("../../config.default.json")
	if err != nil {

	}
	return config
}
