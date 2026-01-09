package flags

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	envkeys "sweep/shared/consts/env-keys"
	consts "sweep/shared/consts/misc"
	tilecontent "sweep/shared/consts/tile-content"
	types "sweep/shared/types"
	paths "sweep/shared/vars/paths"
	styles "sweep/tui/styles"
)

type Flags []types.Flag

const (
	ASCII       types.Flag = "--ascii"
	ASCII_SHORT types.Flag = "--A"

	THEME_PREVIEW       types.Flag = "--preview"
	THEME_PREVIEW_SHORT types.Flag = "--P"

	CONFIG       types.Flag = "--config-path"
	CONFIG_SHORT types.Flag = "--C"

	WIDTH       types.Flag = "--width"
	WIDTH_SHORT types.Flag = "--W"

	HEIGHT       types.Flag = "--height"
	HEIGHT_SHORT types.Flag = "--H"

	MINES       types.Flag = "--mines"
	MINES_SHORT types.Flag = "--M"

	FILL       types.Flag = "--fill"
	FILL_SHORT types.Flag = "--F"

	DEFAULT_CONFIG       types.Flag = "--default-config"
	DEFAULT_CONFIG_SHORT types.Flag = "--D"

	HELP types.Flag = "--help"
)

type NoArgumentProvidedFlagError struct {
	flag types.Flag
}

func (e *NoArgumentProvidedFlagError) Error() string {
	return fmt.Sprintf("argument for flag \"%v\" was not provided", e.flag)
}

func (e *NoArgumentProvidedFlagError) Is(target error) bool {
	return e.Error() == target.Error()
}

type MustBeUin16FlagError struct {
	flag types.Flag
}

func (e *MustBeUin16FlagError) Error() string {
	return fmt.Sprintf("argument for flag \"%v\" must be an unsigned 16 bit integer (0-65535)", e.flag)
}

func (e *MustBeUin16FlagError) Is(target error) bool {
	return e.Error() == target.Error()
}

func validateFlagUint16Argument(args []string, index int) error {
	flag := args[index]

	if index+1 >= len(args) {
		return &NoArgumentProvidedFlagError{flag}
	}
	val := args[index+1]
	_, err := strconv.ParseUint(val, 10, 16)

	if err != nil {
		return &MustBeUin16FlagError{flag}
	}

	return nil
}

func getFlagArgument(args []string, index int) string {
	return args[index+1]
}

type InvalidFlagError struct {
	flag string
}

func (e *InvalidFlagError) Error() string {
	return fmt.Sprintf("(flags) invalid flag - %v", e.flag)
}

func (e *InvalidFlagError) Is(target error) bool {
	return e.Error() == target.Error()
}

func (f Flags) Validate() (bool, []error) {
	skip := false
	args := os.Args[1:]
	flagList := append(args, f...)
	errors := []error{}

	for ix, arg := range flagList {
		if skip {
			skip = false
			continue
		}

		switch arg {
		case HEIGHT, HEIGHT_SHORT, WIDTH, WIDTH_SHORT, MINES, MINES_SHORT:
			skip = true

			if err := validateFlagUint16Argument(flagList, ix); err != nil {
				errors = append(errors, err)
			}
		case ASCII, ASCII_SHORT,
			FILL, FILL_SHORT, CONFIG, CONFIG_SHORT,
			THEME_PREVIEW, THEME_PREVIEW_SHORT,
			DEFAULT_CONFIG, DEFAULT_CONFIG_SHORT,
			HELP:

			continue
		default:
			errors = append(errors, &InvalidFlagError{arg})
		}
	}
	return len(errors) == 0, errors
}

type DefaultConfigReadError struct {
	readFileErr error
}

func (e *DefaultConfigReadError) Error() string {
	hint := "does the file exist?"
	if errors.Is(e.readFileErr, os.ErrNotExist) {
		hint = "does the file exist?"
	}
	if errors.Is(e.readFileErr, os.ErrPermission) {
		hint = "does the program have permissions?"
	}

	return fmt.Sprintf("could not read the default config file at \"%v\": %v", paths.DefaultConfigPath, hint)
}

func (e *DefaultConfigReadError) Is(target error) bool {
	return e.Error() == target.Error()
}

type ConfigWriteError struct {
}

func (e *ConfigWriteError) Error() string {
	hint := "do you have the right permissions?"

	return fmt.Sprintf("could not write the default config to \"%v\": %v", paths.ConfigPath, hint)
}

func (e *ConfigWriteError) Is(target error) bool {
	return e.Error() == target.Error()
}

func ResetConfig() {
	defaultConfig, err := os.ReadFile(paths.DefaultConfigPath)
	if err != nil {
		fmt.Print(&DefaultConfigReadError{})
	}

	err = os.WriteFile(paths.ConfigPath, defaultConfig, 0666)
	if err != nil {
		fmt.Print(&ConfigWriteError{})
	}
	fmt.Printf("the default config was copied to the main config file at %v", paths.ConfigPath)
	os.Exit(0)
}

func (f Flags) Apply() {
	skip := false
	args := os.Args[1:]
	flagList := append(args, f...)

	for ix, arg := range flagList {
		if skip {
			skip = false
			continue
		}

		switch arg {
		case HELP:
			fmt.Print(consts.HelpMessage)
			os.Exit(0)

		case THEME_PREVIEW, THEME_PREVIEW_SHORT:
			os.Setenv(envkeys.Preview, "true")

		case CONFIG, CONFIG_SHORT:
			fmt.Printf("%v\n", paths.ConfigPath)
			os.Exit(0)

		case HEIGHT, HEIGHT_SHORT:
			skip = true
			os.Setenv(envkeys.Height, getFlagArgument(args, ix))

		case WIDTH, WIDTH_SHORT:
			skip = true
			os.Setenv(envkeys.Width, getFlagArgument(args, ix))

		case MINES, MINES_SHORT:
			skip = true
			os.Setenv(envkeys.Mines, getFlagArgument(args, ix))

		case ASCII, ASCII_SHORT:
			tilecontent.SetGlyph(tilecontent.Mine, "M")
			tilecontent.SetGlyph(tilecontent.Flag, "M")
			tilecontent.SetGlyph(tilecontent.WrongFlag, "M")
			tilecontent.SetGlyph(tilecontent.Empty, " ")

			tilecontent.SetGlyph(tilecontent.Zero, "0")
			tilecontent.SetGlyph(tilecontent.One, "1")
			tilecontent.SetGlyph(tilecontent.Two, "2")
			tilecontent.SetGlyph(tilecontent.Three, "3")
			tilecontent.SetGlyph(tilecontent.Four, "4")
			tilecontent.SetGlyph(tilecontent.Five, "5")
			tilecontent.SetGlyph(tilecontent.Six, "6")
			tilecontent.SetGlyph(tilecontent.Seven, "7")
			tilecontent.SetGlyph(tilecontent.Eight, "8")

		case FILL, FILL_SHORT:
			styles.SetFill(true)

		case DEFAULT_CONFIG, DEFAULT_CONFIG_SHORT:
			ResetConfig()
		}
	}
}
