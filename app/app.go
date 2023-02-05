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

	// create the Configuration
	cfg, err := NewConfiguration(configPath)
	if err != nil {
		return nil, err
	}

	// fields
	// to turn on unstructured logging for development profile
	cfg.BuildBooleanField("unstructuredLogger").
		ArgName("ul").
		EnvVar("UNSTRUCTURED_LOGGER").
		ConfigName("Logging", "UseUnstructuredLogger").
		ShortDesc("Use text based(unstructured logger)").
		Default(false).
		Register()

	// logging level (corresponds to zerolog Level values)
	cfg.BuildStringField("logLevel").
		ArgName("ll").
		EnvVar("LOG_LEVEL").
		ConfigName("Logging", "Level").
		ShortDesc("Logging Level").
		Default("INFO").
		Register()

	// The active Profile to start the application under
	cfg.BuildStringField("activeProfile").
		ArgName("pr").
		EnvVar("ACTIVE_PROFILE").
		ConfigName("Logging", "Level").
		ShortDesc("Logging Level").
		Default("INFO").
		Register()

	// Load the config fields
	if err := cfg.LoadFields(os.Args); err != nil {
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

	// get log level from config
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
