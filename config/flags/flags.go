package flags

import (
	"fmt"
	"os"
	"strconv"
	"sweep/config/colors"
	types "sweep/shared/types"
	glyphs "sweep/shared/vars/glyphs"
	styles "sweep/tui/styles"
	themepreview "sweep/tui/theme-preview"
)

type Flags []types.Flag

type BasicConfig struct {
	Mines uint16
	Height uint16
	Width uint16
}
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
)

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

func (f Flags) Apply(colors *colors.Colors) BasicConfig {
	skip := false
	args := os.Args[1:]
	flagList := append(args, f...)
	var basicConfig BasicConfig

	preview := false

	for ix, arg := range flagList {
		if skip {
			skip = false
			continue
		}

		switch arg {
		case THEME_PREVIEW, THEME_PREVIEW_SHORT:
			preview = true

		case CONFIG, CONFIG_SHORT:
			dir, _ := os.Getwd()
			fmt.Printf("%v/config.json\n", dir)
			os.Exit(0)

		case HEIGHT, HEIGHT_SHORT:
			skip = true
			val, err := getFlagIntArgument(args, ix)
			if err != nil {
				fmt.Printf("%v - %v", arg, err)
			} else {
				basicConfig.Height = val
			}

		case WIDTH, WIDTH_SHORT:
			skip = true
			val, err := getFlagIntArgument(args, ix)
			if err != nil {
				fmt.Printf("%v - %v", arg, err)
			} else {
				basicConfig.Width = val
			}

		case MINES, MINES_SHORT:
			skip = true
			val, err := getFlagIntArgument(args, ix)
			if err != nil {
				fmt.Printf("%v - %v", arg, err)
			} else {
				basicConfig.Mines = val
			}

		case ASCII, ASCII_SHORT:
			glyphs.MINE = "M"

			glyphs.FLAG = "F"

			glyphs.WRONG_FLAG = "W"

		case FILL, FILL_SHORT:
			styles.SetFill(true)

		}
	}
	if preview {
		colors.Apply()
		print(themepreview.RenderThemePreview())
		os.Exit(0)
	}
	return basicConfig
}
