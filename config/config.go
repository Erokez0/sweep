package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	bindings "sweep/config/bindings"
	colors "sweep/config/colors"
	cursor "sweep/config/cursor"
	flags "sweep/config/flags"
	envkeys "sweep/shared/consts/env-keys"
	paths "sweep/shared/vars/paths"
	themepreview "sweep/tui/theme-preview"

	gojsonschema "github.com/xeipuuv/gojsonschema"
)

var schema *any

func init() {
	schema = loadSchema(paths.ConfigSchemaPath)
}

type Defaults struct {
	Width  uint16
	Height uint16
	Mines  uint16
}

type Config struct {
	Flags    flags.Flags       `json:"flags,omitempty"`
	Defaults Defaults          `json:"defaults"`
	Colors   colors.Colors     `json:"colors,omitempty"`
	Bindings bindings.Bindings `json:"bindings,omitempty"`

	Mines  uint16 `json:"mines,omitempty"`
	Width  uint16 `json:"width,omitempty"`
	Height uint16 `json:"height,omitempty"`

	Cursor cursor.Cursor `json:"cursor"`
}

func (config *Config) Validate() (bool, []string) {
	configLoader := gojsonschema.NewGoLoader(config)

	schemaLoader := gojsonschema.NewGoLoader(schema)

	errors := []string{}
	result, err := gojsonschema.Validate(schemaLoader, configLoader)
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

	if isValid, colorsErrors := config.Colors.Validate(); !isValid {
		errors = append(errors, colorsErrors...)
	}

	if isValid, cursorErrors := config.Cursor.Validate(); !isValid {
		errors = append(errors, cursorErrors...)
	}

	if isValid, flagErrors := config.Flags.Validate(); !isValid {
		errors = append(errors, flagErrors...)
	}

	if isValid, bindingsErrors := config.Bindings.Validate(); !isValid {
		errors = append(errors, bindingsErrors...)
	}

	return len(errors) == 0, errors
}

func (config *Config) Apply() {
	config.Flags.Apply()
	config.Colors.Apply()
	config.Cursor.Apply()
	config.Bindings.Apply()

	if val, ok := os.LookupEnv(envkeys.Preview); ok && val == "true" {
		fmt.Println(themepreview.RenderThemePreview())
		os.Exit(0)
	}
	// Ignoring errors cause they were
	if val, ok := os.LookupEnv(envkeys.Height); ok {
		parsed, _ := strconv.ParseUint(val, 10, 16)
		config.Height = uint16(parsed)
	}
	if val, ok := os.LookupEnv(envkeys.Width); ok {
		parsed, _ := strconv.ParseUint(val, 10, 16)
		config.Width = uint16(parsed)
	}
	if val, ok := os.LookupEnv(envkeys.Mines); ok {
		parsed, _ := strconv.ParseUint(val, 10, 16)
		config.Mines = uint16(parsed)
	}
}

func loadSchema(schemaPath string) *any {
	schemaBin, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("Could not read config schema %v\nDoes the file exist?", schemaPath)
	}

	schema := new(any)
	err = json.Unmarshal(schemaBin, schema)
	if err != nil {
		log.Fatalf("Could not parse config schema %v\n", schemaPath)
	}

	return schema
}

type loadConfigOpts struct {
	path       string
	config     *Config
	jsonString string
}

func LoadConfig(options *loadConfigOpts) *Config {
	if options == nil {
		log.Fatalf("no path of config struct provided, can not load config")
	}

	flags.Flags{}.Apply()

	var config *Config

	if options.path != "" {
		configBin, err := os.ReadFile(options.path)
		log.SetFlags(log.Lmsgprefix)
		if err != nil {
			log.Fatalf("Could not read config %v\nDoes the file exist?", options.path)
		}
		err = json.Unmarshal(configBin, &config)
		if err != nil {
		log.Fatalf("Could not parse config %v\nInvalid config reference\n", options.path)
		}

	} else if options.jsonString != "" {
		json, err := gojsonschema.NewStringLoader(options.jsonString).LoadJSON()
		if err != nil {
			log.Fatalf("Could not parse config\n%v", err)
		}
		configJson := (json).(Config)
		config = &configJson

	} else {
		config = options.config
	}

	isValid, errors := config.Validate()

	if !isValid {
		fmt.Println("Your config has errors")
		for k, v := range errors {
			fmt.Printf("%v. %v\n", k+1, v)
		}
		os.Exit(1)
	}

	config.Apply()

	return config
}

func GetConfig() *Config {
	config := new(Config)

	config = LoadConfig(&loadConfigOpts{path: paths.ConfigPath})

	return config
}
