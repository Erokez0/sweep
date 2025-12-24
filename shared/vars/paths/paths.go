package paths

import (
	"fmt"
	"os"
	"runtime"
)

const (
	unixBasePath    = "/home/%v/.config/sweep/"
	windowsBasePath = "c:/ProgramData/sweep/"

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
	switch runtime.GOOS {
	case "linux", "darwin":
		basePath = fmt.Sprintf(unixBasePath, os.Getenv("USER"))
	case "windows":
		basePath = windowsBasePath
	}

	ConfigPath = basePath + configName
	ConfigSchemaPath = basePath + configSchemaName
	DefaultConfigPath = basePath + defaultConfigName
}
