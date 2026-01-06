package paths

import (
	"fmt"
	"log"
	"os"
	"runtime"
	consts "sweep/shared/consts/misc"
)

const (
	unixBasePath    = "%v/.config/" + consts.AppName + "/"
	windowsBasePath = "%v\\AppData\\Roaming\\" + consts.AppName + "\\"

	configName        = "config.json"
	configSchemaName  = "config.schema.json"
	defaultConfigName = "config.default.json"
)

var (
	ConfigPath        string
	ConfigSchemaPath  string
	DefaultConfigPath string
)

func init() {
	var basePath string

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not determine users home directory")
	}

	switch runtime.GOOS {
	case "linux", "darwin":
		basePath = fmt.Sprintf(unixBasePath, home)
	case "windows":
		basePath = fmt.Sprintf(windowsBasePath, home)
	}

	ConfigPath = basePath + configName
	ConfigSchemaPath = basePath + configSchemaName
	DefaultConfigPath = basePath + defaultConfigName
}
