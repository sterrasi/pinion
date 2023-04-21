package app

import (
	"github.com/sterrasi/pinion"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// Expect a SystemConfigurationErrorCode when creating a Configuration with a bad path
func TestNewConfigurationWithBadPath(t *testing.T) {

	path := "./doesNotExist.ini"
	_, err := NewConfiguration(path)
	if err == nil {
		t.Fatalf("expecting NewConfiguration to fail on bad config file path")
	}

	assert.Equal(t, err.Code(), SystemConfigurationErrorCode)
	assert.Equal(t, err.GetMetadataValue("path"), path)
}

// Validates the mappings for values between a FieldBuilder and Field
func TestBuilderMappingsForField(t *testing.T) {

	cfg := createConfiguration(t)
	f := registerIntegerField(createRegistry(cfg))

	assert.Equal(t, f.ArgName, "p")
	assert.Equal(t, f.EnvVar, "PORT")
	assert.Equal(t, f.ConfigSectionName, "Server")
	assert.Equal(t, f.ConfigFieldName, "Port")
	assert.Equal(t, f.ShortDescription, "Http port")
	assert.Equal(t, f.LongDescription, "Http server port")
	assert.Equal(t, f.DefaultValue, 3000)
}

// Validate the ValueMetadata and accessor method for an integer field pulled from a registry file
func TestParseIntFieldFromConfig(t *testing.T) {

	cfg := createConfiguration(t)
	port := registerIntegerField(createRegistry(cfg))
	if err := cfg.LoadFields([]string{}); err != nil {
		t.Fatalf("Error loading fields: %s", err.Error())
	}

	assertMetadata(t, cfg, port, 4000, File)
}

// Validate the ValueMetadata and accessor method for a uint field pulled from a registry file
func TestParseUintFieldFromConfig(t *testing.T) {

	cfg := createConfiguration(t)
	poolSize := registerUintField(createRegistry(cfg))
	if err := cfg.LoadFields([]string{}); err != nil {
		t.Fatalf("Error loading fields: %s", err.Error())
	}

	assertMetadata(t, cfg, poolSize, uint(20), File)
}

// Validate the ValueMetadata and accessor method for a float field pulled from a registry file
func TestParseFloatFieldFromConfig(t *testing.T) {

	cfg := createConfiguration(t)
	mrFloaty := registerFloatField(createRegistry(cfg))
	if err := cfg.LoadFields([]string{}); err != nil {
		t.Fatalf("Error loading fields: %s", err.Error())
	}

	assertMetadata(t, cfg, mrFloaty, float64(3.54321), File)
}

// Validate the ValueMetadata and accessor method for a string field pulled from a registry file
func TestParseStringFieldFromConfig(t *testing.T) {

	cfg := createConfiguration(t)
	host := registerStringField(createRegistry(cfg))
	if err := cfg.LoadFields([]string{}); err != nil {
		t.Fatalf("Error loading fields: %s", err.Error())
	}

	assertMetadata(t, cfg, host, "localhost", File)
}

// Validate the ValueMetadata and accessor method for a bool field pulled from a registry file
func TestParseBoolFieldFromConfig(t *testing.T) {

	cfg := createConfiguration(t)
	verbosity := registerBoolField(createRegistry(cfg))
	if err := cfg.LoadFields([]string{}); err != nil {
		t.Fatalf("Error loading fields: %s", err.Error())
	}

	assertMetadata(t, cfg, verbosity, true, File)
}

// Make sure the Default value gets picked up when the field is not defined anywhere else
func TestIntFieldDefaultValue(t *testing.T) {
	cfg := createConfiguration(t)
	unknown := registerUnknownIntegerFieldWithDefault(createRegistry(cfg), 3)
	if err := cfg.LoadFields([]string{}); err != nil {
		t.Fatalf("Error loading fields: %s", err.Error())
	}

	assertMetadata(t, cfg, unknown, 3, None)
}

// Make sure the LoadFields fails when a RequiredField cannot be found
func TestRequiredIntFieldWithNoValue(t *testing.T) {
	cfg := createConfiguration(t)
	numCats := registerRequiredUnknownIntegerField(createRegistry(cfg))
	err := cfg.LoadFields([]string{})
	if err == nil {
		t.Fatalf("expecting LoadFields to fail on 'cannot find required field'")
	}
	assert.Equal(t, err.Code(), SystemConfigurationErrorCode)
	assert.Equal(t, err.GetMetadataValue("fieldName"), numCats.Name)
}

// Test environment variable override
func TestEnvironmentVariableOverride(t *testing.T) {

	cfg := createConfiguration(t)
	port := registerIntegerField(createRegistry(cfg))

	// set the PORT env variable
	orig := os.Getenv("PORT")
	_ = os.Setenv("PORT", "8080")
	defer func() {
		os.Setenv("PORT", orig)
	}()

	// use cli arg as well to demonstrate that the env var takes higher precedence
	if err := cfg.LoadFields([]string{"appName", "-p", "6000"}); err != nil {
		t.Fatalf("Error loading fields: %s", err.Error())
	}

	assertMetadata(t, cfg, port, 8080, EnvironmentVar)
}

// Test cli arg override
// Test environment variable override
func TestCliArgOverride(t *testing.T) {

	cfg := createConfiguration(t)
	port := registerIntegerField(createRegistry(cfg))

	// use cli arg as well to demonstrate that the env var takes higher precedence
	if err := cfg.LoadFields([]string{"appName", "-p", "6000"}); err != nil {
		t.Fatalf("Error loading fields: %s", err.Error())
	}

	assertMetadata(t, cfg, port, 6000, CommandLine)
}

func assertMetadata(t *testing.T, cfg *Configuration, f *Field, expectedVal any, specifier FieldSpecifier) {
	md := cfg.GetValueMetadata(f.Name)
	assert.Equal(t, md.Field, f)

	if md.Field.Type == Float {

		if !pinion.IsWithinRange(md.Value, expectedVal, 0.1) {
			t.Fatalf("Float values are not equal: md.Value=%f, expected=%f", md.Value.(float64),
				expectedVal.(float64))
		}

	} else {
		assert.Equal(t, md.Value, expectedVal)
	}

	assert.Equal(t, md.Specifier, specifier)

	switch f.Type {
	case Int:
		np, err := cfg.GetIntValue(f.Name)
		if err != nil {
			t.Fatalf("Error asserting metadata for int field %s : %s", f.Name, err.Error())
		}
		assert.Equal(t, expectedVal, *np)

	case Uint:
		up, err := cfg.GetUintValue(f.Name)
		if err != nil {
			t.Fatalf("Error asserting metadata for uint field %s : %s", f.Name, err.Error())
		}
		assert.Equal(t, expectedVal, *up)

	case Float:
		fp, err := cfg.GetFloatValue(f.Name)
		if err != nil {
			t.Fatalf("Error asserting metadata for float field %s : %s", f.Name, err.Error())
		}
		if !pinion.IsWithinRange(md.Value, expectedVal, 0.1) {
			t.Fatalf("Float values are not equal: md.Value=%f, expected=%f", *fp,
				expectedVal.(float64))
		}

	case Bool:
		bv, err := cfg.GetBoolValue(f.Name)
		if err != nil {
			t.Fatalf("Error asserting metadata for boolean field %s : %s", f.Name, err.Error())
		}
		assert.Equal(t, expectedVal, *bv)

	case String:
		sv, err := cfg.GetStringValue(f.Name)
		if err != nil {
			t.Fatalf("Error asserting metadata for string field %s : %s", f.Name, err.Error())
		}
		assert.Equal(t, expectedVal, *sv)
	}
}

func createConfiguration(t *testing.T) *Configuration {
	cfg, err := NewConfiguration("./test/application.ini")
	if err != nil {
		t.Fatalf("Error initializing configuration: %s", err.Error())
	}
	return cfg
}

// Registers an integer field with default value that does not exist in the configuration
func registerUnknownIntegerFieldWithDefault(reg *FieldRegistry, defaultValue int) *Field {

	f := reg.CreateIntField("numCats").
		ArgName("cats").
		EnvVar("NUM_CATS").
		ConfigName("Animals", "NumCats").
		Descriptions("Number of Cats", "Number of Cats").
		Default(defaultValue).
		Register()

	return f
}

// Registers a required integer field with no default value that does not exist in the configuration
func registerRequiredUnknownIntegerField(reg *FieldRegistry) *Field {
	f := reg.CreateIntField("numCats").
		ArgName("cats").
		EnvVar("NUM_CATS").
		ConfigName("Animals", "NumCats").
		Descriptions("Number of Cats", "Number of Cats").
		Required().
		Register()

	return f
}

// Registers and returns the 'port' field
func registerIntegerField(reg *FieldRegistry) *Field {
	f := reg.CreateIntField("port").
		ArgName("p").
		EnvVar("PORT").
		ConfigName("Server", "Port").
		Descriptions("Http port", "Http server port").
		Default(3000).
		Register()

	return f
}

// Registers and returns the 'poolSize' field
func registerUintField(reg *FieldRegistry) *Field {
	f := reg.CreateUintField("poolSize").
		ArgName("ps").
		EnvVar("POOL_SIZE").
		ConfigName("Database", "PoolSize").
		Descriptions("Connection Pool Size", "Database Connection Pool Size").
		Default(10).
		Register()

	return f
}

// Registers and returns the 'floaty' field
func registerFloatField(reg *FieldRegistry) *Field {
	f := reg.CreateFloatField("floaty").
		ArgName("f").
		EnvVar("FLOATY").
		ConfigName("Float", "MrFloaty").
		Descriptions("Made up float value", "Made up float value").
		Default(10.12345).
		Register()

	return f
}

// Registers and returns the 'host' field
func registerStringField(reg *FieldRegistry) *Field {
	f := reg.CreateStringField("host").
		ArgName("h").
		EnvVar("HOST").
		ConfigName("Server", "Host").
		Descriptions("Server Host", "Server Host Address").
		Default("myHost.com").
		Register()
	return f
}

// Registers and returns the 'useHttps' field
func registerBoolField(reg *FieldRegistry) *Field {
	f := reg.CreateBooleanField("verbose").
		ArgName("v").
		EnvVar("VERBOSE").
		ConfigName("Server", "Verbose").
		Descriptions("Is Verbose", "Is Verbose").
		Default(false).
		Register()
	return f
}

func createRegistry(c *Configuration) *FieldRegistry {
	return &FieldRegistry{
		fields: c.fields,
	}
}
