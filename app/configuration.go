package app

import (
	"errors"
	"github.com/sterrasi/pinion/logger"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

// Configuration encapsulates the configuration for an application
type Configuration struct {
	base      *ini.File
	overrides *ini.File
	cliArgs   *CLIArgs
	fields    map[string]*Field
	values    map[string]*ValueMetadata
}

// NewConfigurationFromContents creates a Configuration from the given ini contents
func NewConfigurationFromContents(contents string) (*Configuration, Error) {

	baseCfg, err := ini.Load([]byte(contents))
	if err != nil {
		return nil, BuildSysConfigError().Cause(err).
			Str("contents", contents).
			Msg("Error loading ini file from contents")
	}

	cfg := &Configuration{
		base:   baseCfg,
		fields: make(map[string]*Field),
		values: make(map[string]*ValueMetadata)}

	return cfg, nil
}

// NewConfiguration creates a Configuration from the given ini file path
func NewConfiguration(path string) (*Configuration, Error) {

	// make sure the file exists
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil, BuildSysConfigError().Cause(err).
			Str("path", path).
			Msg("Config file does not exist")
	}

	// load the ini file
	baseCfg, err := ini.Load(path)
	if err != nil {
		return nil, BuildSysConfigError().Cause(err).
			Str("path", path).
			Msg("Error loading base config file")
	}

	// if there is an active profile configuration then load it
	var profileCfg *ini.File
	envPath := getEnvironmentConfigPath(path)
	if _, err = os.Stat(envPath); err == nil {
		profileCfg, err = ini.Load(envPath)
		if err != nil {
			return nil, BuildSysConfigError().Cause(err).
				Str("path", envPath).
				Msg("Error loading environment config file")
		}
	}

	// set the profile level registry
	cfg := &Configuration{
		base:      baseCfg,
		overrides: profileCfg,
		fields:    make(map[string]*Field),
		values:    make(map[string]*ValueMetadata)}

	return cfg, nil
}

func (c *Configuration) GetIntValue(fieldName string) (*int, Error) {
	val, err := c.getValue(fieldName, Int)
	if err != nil {
		return nil, err
	}
	n := val.(int)
	return &n, nil
}

func (c *Configuration) GetUintValue(fieldName string) (*uint, Error) {
	val, err := c.getValue(fieldName, Uint)
	if err != nil {
		return nil, err
	}
	u := val.(uint)
	return &u, nil
}

func (c *Configuration) GetFloatValue(fieldName string) (*float64, Error) {
	val, err := c.getValue(fieldName, Float)
	if err != nil {
		return nil, err
	}
	f := val.(float64)
	return &f, nil
}

func (c *Configuration) GetStringValue(fieldName string) (*string, Error) {
	val, err := c.getValue(fieldName, String)
	if err != nil {
		return nil, err
	}
	s := val.(string)
	return &s, nil
}

func (c *Configuration) GetBoolValue(fieldName string) (*bool, Error) {
	val, err := c.getValue(fieldName, Bool)
	if err != nil {
		return nil, err
	}
	b := val.(bool)
	return &b, nil
}

func (c *Configuration) getValue(fieldName string, expectedType ValueType) (any, Error) {
	md, present := c.values[fieldName]
	if !present {
		return nil, BuildNotFoundError().
			Str("fieldName", fieldName).
			Str("valueType", expectedType.String()).
			Msgf("%s field not found", expectedType.String())
	}
	if md.Field.Type != expectedType {
		return nil, BuildIllegalStateError().
			Str("fieldName", fieldName).
			Str("fieldValueType", md.Field.Type.String()).
			Str("expectedValueType", expectedType.String()).
			Msg("Field value type does not match expected value type")
	}

	return md.Value, nil
}

// LoadFields loads the values from the registered fields into the configuration
func (c *Configuration) LoadFields(cliArgs []string) Error {

	// if no fields exist then do nothing
	if len(c.fields) == 0 {
		logger.Debug().Msg("No fields were registered")
		return nil
	}

	// clear the field values
	for k := range c.values {
		delete(c.values, k)
	}

	// parse the args into the field values
	args, err := parseArgs(cliArgs, c.fieldsByArgName())
	if err != nil {
		return err
	}
	c.cliArgs = args

	// resolve each field
	for k, f := range c.fields {

		logger.Trace().Str("field", k).Msg("loading field")
		val, err := c.loadField(f, args)
		if err != nil {
			return err
		}
		c.values[f.Name] = val
	}

	return nil
}

// GetValueMetadata returns the metadata obtained when parsing a Field with the associated fieldName
func (c *Configuration) GetValueMetadata(fieldName string) *ValueMetadata {
	return c.values[fieldName]
}

// loadField returns the string value of the given Key
func (c *Configuration) loadField(field *Field, args *CLIArgs) (*ValueMetadata, Error) {

	// check for an environment variable
	if field.EnvVar != "" {
		envVarVal := strings.TrimSpace(os.Getenv(field.EnvVar))
		if envVarVal != "" {
			v, err := newValue(field, envVarVal, EnvironmentVar)
			if err != nil {
				return nil, err
			}
			return v, nil
		}
	}

	// check for a command line argument
	argValue, present := args.fieldValues[field.Name]
	if present {
		v, err := newValue(field, argValue, CommandLine)
		if err != nil {
			return nil, err
		}
		return v, nil
	}

	// fetch from the ini file
	pv, err := c.getIniValue(field)
	if err != nil {
		return nil, err
	}

	specifier := File

	// check for a default value
	if pv == nil {
		specifier = None
		// if the field was marked required then error
		if field.Required {
			return nil, BuildSysConfigError().
				Str("fieldName", field.Name).
				Msg("No value specified for required field")
		}
		// if a default was provided then set it.. else the value is optional and nil
		if field.DefaultValue != nil {
			tmp := field.Type.ToString(field.DefaultValue)
			pv = &tmp
		}
	}

	v, err := newValue(field, *pv, specifier)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// getIniValue returns the string value of the given Key
func (c *Configuration) getIniValue(field *Field) (*string, Error) {

	// if the section exists
	pv, err := getValueFromIni(c.base, field)
	if err != nil {
		return nil, err
	}
	if pv == nil && c.overrides != nil {
		return getValueFromIni(c.overrides, field)
	}
	return pv, nil
}

func (c *Configuration) fieldsByArgName() map[string]*Field {
	fba := make(map[string]*Field)
	for _, f := range c.fields {
		fba[f.ArgName] = f
	}
	return fba
}

func getValueFromIni(file *ini.File, field *Field) (*string, Error) {

	if file.HasSection(field.ConfigSectionName) {
		sect, err := file.GetSection(field.ConfigSectionName)
		if err != nil {
			return nil, BuildSysConfigError().Cause(err).
				Str("sectionName", field.ConfigSectionName).
				Msg("Error retrieving base ini section")
		}
		if sect.HasKey(field.ConfigFieldName) {
			key, err := sect.GetKey(field.ConfigFieldName)
			if err != nil {
				return nil, BuildSysConfigError().Cause(err).
					Str("sectionName", field.ConfigSectionName).
					Str("fieldName", field.ConfigFieldName).
					Msg("Error retrieving field from base ini section")
			}

			value := key.Value()
			if value != "" {
				return &value, nil
			}
		}
	}
	return nil, nil
}

func getEnvironmentConfigPath(basePath string) string {
	ext := filepath.Ext(basePath)
	return strings.TrimSuffix(basePath, ext) + "." + GetActiveProfile().String() + ext
}
