package consts

const (
	AppAsciiLogo = `███████╗██╗    ██╗███████╗███████╗██████╗ 
██╔════╝██║    ██║██╔════╝██╔════╝██╔══██╗
███████╗██║ █╗ ██║█████╗  █████╗  ██████╔╝
╚════██║██║███╗██║██╔══╝  ██╔══╝  ██╔═══╝ 
███████║╚███╔███╔╝███████╗███████╗██║     
╚══════╝ ╚══╝╚══╝ ╚══════╝╚══════╝╚═╝`
	AppName     = "sweep"
	HelpMessage = "Usage " + AppName + ` [OPTION] ...
List of options:
  --help                  display help and exit
  --D, --default-config   copy the default configuration file to the 
                            destination of the main configuration file and exit
  --C, --config-path        print out the canonical path of the 
                              configuration file and exit
  --P, --preview            print out a preview of your color theme and exit

  --A, --ascii              replace glyphs for tiles (like flags and mines) 
                              to ASCII friendly characters
  --F, --fill               vary how the color affect the tiles. If the option 
                              is provided, then the colors will fill the background of 
                              the tiles, otherwise they will fill the foreground

  --M, --mines[ uint16]     sets the desired amount of mines to the field
                              if other field arguments are set
  --W, --width[ uint16]     sets the desired field width in columns 
                              if other field arguments are set
  --H, --height[ uint16]    sets the desired field height in rows 
                              if other field arguments are set
`
)
