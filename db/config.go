package db

import "github.com/sterrasi/pinion/app"

const dbSectionName = "Database"

const dbHostFieldName = "dbHost"
const dbNameFieldName = "dbName"
const dbSchemaFieldName = "dbSchema"
const dbUserFieldName = "dbUser"
const dbPasswordFieldName = "dbPassword"
const maxIdleConnectionsFieldName = "maxIdleConnections"
const maxOpenConnectionsFieldName = "maxOpenConnections"

// DbConfig contains values required to connect to a database
type DbConfig struct {
	DbName             string
	Host               string
	User               string
	Schema             string
	Password           string
	MaxIdleConnections uint
	MaxOpenConnections uint
}

// RegisterConfig will register the config field definitions needed for connecting to a database
func RegisterConfig(reg app.FieldRegistry) {

	// database host (Required)
	// ex. localhost:5432
	reg.CreateStringField(dbHostFieldName).
		ArgName("db-host").
		EnvVar("DB_HOST").
		ConfigName(dbSectionName, "Host").
		ShortDesc("Database Host URI").
		Required().
		Register()

	// database name (Required)
	// ex. test_db
	reg.CreateStringField(dbNameFieldName).
		ArgName("db-name").
		EnvVar("DB_NAME").
		ConfigName(dbSectionName, "Name").
		ShortDesc("Database Name").
		Required().
		Register()

	// database schema (Required)
	reg.CreateStringField(dbSchemaFieldName).
		ArgName("db-schema").
		EnvVar("DB_SCHEMA").
		ConfigName(dbSectionName, "Schema").
		ShortDesc("Database Schema").
		Required().
		Register()

	// database user (Required)
	reg.CreateStringField(dbUserFieldName).
		ArgName("db-user").
		EnvVar("DB_USER").
		ConfigName(dbSectionName, "Schema").
		ShortDesc("Database Schema").
		Required().
		Register()

	// database password [encrypted] (Required)
	reg.CreateStringField(dbPasswordFieldName).
		ArgName("db-password").
		ConfigName(dbSectionName, "Password").
		ShortDesc("Database Password").
		Required().
		Register()

	// max idle connections
	reg.CreateUintField(maxIdleConnectionsFieldName).
		ArgName("max-idle-connections").
		ConfigName(dbSectionName, "MaxIdleConnections").
		ShortDesc("Max number of idle database connections").
		Default(30).
		Register()

	// max open connections
	reg.CreateUintField(maxOpenConnectionsFieldName).
		ArgName("max-open-connections").
		ConfigName(dbSectionName, "MaxOpenConnections").
		ShortDesc("Max number of open database connections").
		Default(20).
		Register()
}

// NewDbConfig creates a DbConfig from the given parsed app.Configuration
func NewDbConfig(cfg app.Configuration) (*DbConfig, app.Error) {

	dbHost, err := cfg.GetStringValue(dbHostFieldName)
	if err != nil {
		return nil, err
	}

	dbName, err := cfg.GetStringValue(dbNameFieldName)
	if err != nil {
		return nil, err
	}

	dbSchema, err := cfg.GetStringValue(dbSchemaFieldName)
	if err != nil {
		return nil, err
	}

	dbUser, err := cfg.GetStringValue(dbUserFieldName)
	if err != nil {
		return nil, err
	}

	dbPassword, err := cfg.GetStringValue(dbPasswordFieldName)
	if err != nil {
		return nil, err
	}

	maxIdleConnections, err := cfg.GetUintValue(maxIdleConnectionsFieldName)
	if err != nil {
		return nil, err
	}

	maxOpenConnections, err := cfg.GetUintValue(maxOpenConnectionsFieldName)
	if err != nil {
		return nil, err
	}

	return &DbConfig{
		DbName:             *dbName,
		Host:               *dbHost,
		Schema:             *dbSchema,
		User:               *dbUser,
		Password:           *dbPassword,
		MaxIdleConnections: *maxIdleConnections,
		MaxOpenConnections: *maxOpenConnections,
	}, nil
}
