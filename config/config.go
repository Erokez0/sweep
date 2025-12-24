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

var (
	schema = gojsonschema.NewReferenceLoader("file:///home/erokez/Desktop/code/sweep/config.schema.json")
)

func (config *Config) validate() (bool, []string) {
	configLoader := gojsonschema.NewGoLoader(config)
	errors := []string{}

	result, err := gojsonschema.Validate(schema, configLoader)
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

func LoadConfig(configPath string) (*Config, error) {
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

	return config, nil
}

func GetConfig() *Config {
	config := new(Config)

	var err error
	config, err = LoadConfig("/home/erokez/Desktop/code/sweep/config.json")
	if err == nil {
		return config
	}

	return config
}
