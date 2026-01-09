package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	bindings "sweep/config/bindings"
	colors "sweep/config/colors"
	cursor "sweep/config/cursor"
	flags "sweep/config/flags"
	glyphs "sweep/config/glyphs"
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
	Cursor   cursor.Cursor     `json:"cursor"`
	Glyphs   glyphs.Glyphs     `json:"glyphs"`

	Mines  uint16 `json:"mines,omitempty"`
	Width  uint16 `json:"width,omitempty"`
	Height uint16 `json:"height,omitempty"`
}

type ConfigValidationError struct {
	errors []error
}

func (e *ConfigValidationError) Error() string {
	var msg strings.Builder
	msg.WriteString("Your configuration has errors")
	for index, err := range e.errors {
		msg.WriteString(fmt.Sprintf("%v. %v", index+1, err))
	}
	return msg.String()
}

type GoJsonSchemaConfigValidationError struct {
	err gojsonschema.ResultError
}

func (e *GoJsonSchemaConfigValidationError) Error() string {
	return e.err.String()
}

func (config *Config) Validate() (bool, []error) {
	configLoader := gojsonschema.NewGoLoader(config)

	schemaLoader := gojsonschema.NewGoLoader(schema)

	errors := []error{}

	result, err := gojsonschema.Validate(schemaLoader, configLoader)
	if err != nil {
		errors = append(errors, err)
	}

	if !result.Valid() {
		schemaErrors := make([]error, len(result.Errors()))
		for ix, err := range result.Errors() {
			schemaErrors[ix] = &GoJsonSchemaConfigValidationError{err}
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

	if isValid, glyphsErrors := config.Glyphs.Validate(); !isValid {
		errors = append(errors, glyphsErrors...)
	}

	return len(errors) == 0, errors
}

func (config *Config) Apply() {
	config.Cursor.Apply()
	config.Bindings.Apply()
	config.Glyphs.Apply()
	config.Colors.Apply()
	config.Flags.Apply()

	if val, ok := os.LookupEnv(envkeys.Preview); ok && val == "true" {
		fmt.Println(themepreview.RenderThemePreview())
		os.Exit(0)
	}

	// Ignoring errors cause they were accounted for during validation
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

type EmptyConfigLoadOptionsError struct{}

func (e *EmptyConfigLoadOptionsError) Error() string {
	return "empty config load options provided"
}

type ConfigReadFileError struct {
	readFileErr error
	configPath  string
}

func (e *ConfigReadFileError) Error() string {
	return fmt.Sprintf("could not read config \"%v\": does the file exist?", e.configPath)
}

type ConfigParsingError struct {
	unmarshalError error
	configPath     string
}

func (e *ConfigParsingError) Error() string {
	var errorMsg string
	if errors.Is(e.unmarshalError, &json.SyntaxError{}) {
		errorMsg = "invalid JSON syntax"
	} else {
		errorMsg = e.unmarshalError.Error()
	}
	return fmt.Sprintf("could not parse config \"%v\": %v", e.configPath, errorMsg)

}

type JsonStringLoadConfigError struct {
	err error
}

func (e *JsonStringLoadConfigError) Error() string {
	return fmt.Sprintf("could not parse config from string: %v", e.err)
}

func LoadConfig(options *loadConfigOpts) (*Config, error) {
	if options == nil {
		return nil, &EmptyConfigLoadOptionsError{}
	}

	flags.ApplyFromArgs()

	var config *Config

	if options.path != "" {
		configBin, err := os.ReadFile(options.path)
		log.SetFlags(log.Lmsgprefix)
		if err != nil {
			return nil, &ConfigReadFileError{err, options.path}
		}
		err = json.Unmarshal(configBin, &config)
		if err != nil {
			return nil, &ConfigParsingError{err, options.path}
		}
	} else if options.jsonString != "" {
		json, err := gojsonschema.NewStringLoader(options.jsonString).LoadJSON()
		if err != nil {
			return nil, &JsonStringLoadConfigError{err}
		}
		configJson := (json).(Config)
		config = &configJson
	} else {
		config = options.config
	}

	isValid, errors := config.Validate()

	if !isValid {
		return nil, &ConfigValidationError{errors}
	}

	config.Apply()

	return config, nil
}

func GetConfig() *Config {
	config, err := LoadConfig(&loadConfigOpts{path: paths.ConfigPath})

	if err != nil {
		log.Fatal(err)
	}

	return config
}
