package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	bindings "sweep/config/bindings"
	colors "sweep/config/colors"
	cursor "sweep/config/cursor"
	flags "sweep/config/flags"
	"sweep/shared/vars/paths"

	gojsonschema "github.com/xeipuuv/gojsonschema"
)

type Defaults struct {
	Width  uint16
	Height uint16
	Mines  uint16
}



type Config struct {
	Flags    flags.Flags       `json:"flags"`
	Defaults Defaults          `json:"defaults"`
	Colors   colors.Colors     `json:"colors"`
	Bindings bindings.Bindings `json:"bindings"`

	Mines  uint16 `json:"mines,omitempty"`
	Width  uint16 `json:"width,omitempty"`
	Height uint16 `json:"height,omitempty"`

	Cursor  cursor.Cursor `json:"cursor"`
}

func (config *Config) validate() (bool, []string) {
	configLoader := gojsonschema.NewGoLoader(config)
	
	schema := loadSchema(paths.ConfigSchemaPath)
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

	if isValid, cursorErrors := config.Bindings.Validate(); !isValid {
		errors = append(errors, cursorErrors...)
	}

	return len(errors) == 0, errors
}



func (config *Config) Apply() {
	basic := config.Flags.Apply(&config.Colors)
	config.Mines = basic.Mines
	config.Height = basic.Height
	config.Width = basic.Width

	config.Colors.Apply()
	config.Cursor.Apply()
}

func loadSchema(schemaPath string) (*any) {
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

func LoadConfig(configPath string) (*Config) {
	configBin, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Could not read config %v\nDoes the file exist?", configPath)
	}

	config := new(Config)
	err = json.Unmarshal(configBin, config)
	if err != nil {
		log.Fatalf("Could not parse config %v\n", configPath)
	}

	isValid, errors := config.validate()

	if !isValid {
		fmt.Println("Your config file has errors")
		for k, v := range errors {
			fmt.Printf("%v. %v\n", k+1, v)
		}
		os.Exit(1)
	}

	config.Apply()

	return config
}

// TODO add fallback to default config
func GetConfig() *Config {
	config := new(Config)

	// var err error
	config = LoadConfig(paths.ConfigPath)

	// config, err = LoadConfig(paths.DefaultConfigPath)
	// if err != nil {
		// return config
	// }

	return config
}
