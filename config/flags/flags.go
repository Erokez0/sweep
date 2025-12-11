package Flags

import "sweep/types"

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
