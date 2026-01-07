package flags

import (
	"fmt"
	"os"
	"strconv"

	envkeys "sweep/shared/consts/env-keys"
	consts "sweep/shared/consts/misc"
	types "sweep/shared/types"
	glyphs "sweep/shared/vars/glyphs"
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

func validateFlagUint16Argument(args []string, index int) (bool, []string) {
	errors := []string{}
	if index+1 >= len(args) {
		return false, []string{fmt.Sprintf("argument for flag \"%v\" was not provided", args[index])}
	}
	val := args[index+1]
	_, err := strconv.ParseUint(val, 10, 16)

	if err != nil {
		errors = append(errors, fmt.Sprintf("argument for flag \"%v\" must be a unsigned 16 bit integer (0-65535)", args[index]))
	}

	return len(errors) == 0, errors
}

func getFlagArgument(args []string, index int) string {
	return args[index+1]
}

func (f Flags) Validate() (bool, []string) {
	skip := false
	args := os.Args[1:]
	flagList := append(args, f...)
	errors := []string{}

	for ix, arg := range flagList {
		if skip {
			skip = false
			continue
		}

		switch arg {
		case HEIGHT, HEIGHT_SHORT, WIDTH, WIDTH_SHORT, MINES, MINES_SHORT:
			skip = true
			if isValid, flagErrors := validateFlagUint16Argument(args, ix); !isValid {
				errors = append(errors, flagErrors...)
			}
		case ASCII, ASCII_SHORT,
			FILL, FILL_SHORT, CONFIG, CONFIG_SHORT,
			THEME_PREVIEW, THEME_PREVIEW_SHORT,
			DEFAULT_CONFIG, DEFAULT_CONFIG_SHORT,
			HELP:

			continue
		default:
			errors = append(errors, fmt.Sprintf("invalid flag \"%v\"", arg))
		}
	}
	return len(errors) == 0, errors
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
			glyphs.Mine = "M"

			glyphs.Flag = "F"

			glyphs.WrongFlag = "W"

		case FILL, FILL_SHORT:
			styles.SetFill(true)

		case DEFAULT_CONFIG, DEFAULT_CONFIG_SHORT:
			defaultConfig, err := os.ReadFile(paths.DefaultConfigPath)

			if err != nil {
				fmt.Printf("could not read the default config file at %v\nDoes the file exist?", paths.DefaultConfigPath)
				os.Exit(1)
			}

			err = os.WriteFile(paths.ConfigPath, defaultConfig, 0666)
			if err != nil {
				fmt.Printf("could not write the default config to %v\nDo you have the right permissions?", paths.ConfigPath)
				os.Exit(1)
			}
			fmt.Printf("the default config was copied to the main config file at %v", paths.ConfigPath)
			os.Exit(0)
		}
	}
}
