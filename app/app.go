package app

import (
	"github.com/rs/zerolog"
	"github.com/sterrasi/pinion/logger"
	"os"
)

type Application struct {
	name          string
	configuration *Configuration
	profile       Profile
}

// Create the Application.  This should be done after configuration fields are registered
func Create(configPath string, name string) (*Application, Error) {
	return CreateWithBuilder(configPath, name, nil)
}

// CreateWithBuilder creates the Application using the additional ConfigurationBuilderFn to add
// application specific configuration
func CreateWithBuilder(configPath string, name string, builderFn ConfigurationBuilderFn) (*Application, Error) {

	// create the Configuration
	cfg, err := NewConfiguration(configPath)
	if err != nil {
		return nil, err
	}

	registry := &FieldRegistry{
		fields: cfg.fields,
	}

	// fields
	// to turn on unstructured logging for development profile
	registry.CreateBooleanField("unstructuredLogger").
		ArgName("ul").
		EnvVar("UNSTRUCTURED_LOGGER").
		ConfigName("Logging", "UseUnstructuredLogger").
		ShortDesc("Use text based(unstructured logger)").
		Default(false).
		Register()

	// logging level (corresponds to zerolog Level values)
	registry.CreateStringField("logLevel").
		ArgName("ll").
		EnvVar("LOG_LEVEL").
		ConfigName("Logging", "Level").
		ShortDesc("Logging Level").
		Default("INFO").
		Register()

	// The active Profile to start the application under
	registry.CreateStringField("activeProfile").
		ArgName("pr").
		EnvVar("ACTIVE_PROFILE").
		ConfigName("Logging", "Level").
		ShortDesc("Logging Level").
		Default("INFO").
		Register()

	// Load the registry fields needed to initialize logging and the active profile
	cfg.fields = registry.fields
	if err = cfg.LoadFields(os.Args); err != nil {
		return nil, err
	}

	// Get the active profile
	pProfileVal, err := cfg.GetStringValue("activeProfile")
	if err != nil {
		return nil, err
	}
	pProfile, err := parseProfile(*pProfileVal)
	if err != nil {
		return nil, err
	}

	// configure the root logger
	err = configureRootLogger(cfg, *pProfile)
	if err != nil {
		return nil, err
	}

	// now that logging and the active profile are initialized go ahead and call the
	// builder function if specified
	if builderFn != nil {
		if err = builderFn(registry); err != nil {
			return nil, err
		}

		// Load any additional fields specified on the registry as a result of calling the builder function
		cfg.fields = registry.fields
		if err = cfg.LoadFields(os.Args); err != nil {
			return nil, err
		}
	}

	// create the application
	app := &Application{
		name:          name,
		configuration: cfg,
		profile:       *pProfile,
	}

	return app, nil
}

// configureRootLogger configure the root logger for the application
func configureRootLogger(cfg *Configuration, profile Profile) Error {

	// get log level from registry
	pLogLevelVal, err := cfg.GetStringValue("logLevel")
	if err != nil {
		return err
	}
	logLevel, er := zerolog.ParseLevel(*pLogLevelVal)
	if er != nil {
		return BuildSysConfigError().Str("logLevel", *pLogLevelVal).
			Cause(er).
			Msg("Error interpreting configured log level")
	}

	// determine if an unstructured log should be used
	unstructuredLogMd := cfg.GetValueMetadata("unstructuredLogger")
	useUnstructuredLog := unstructuredLogMd.Value.(bool)
	if unstructuredLogMd.Specifier == None && profile != Production {
		useUnstructuredLog = true
	}

	logger.ConfigureLogging(logLevel, useUnstructuredLog)
	return nil
}
