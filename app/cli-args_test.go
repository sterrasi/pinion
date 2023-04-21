package app

import "testing"

func TestParseArgs(t *testing.T) {

	cfg := createConfiguration(t)

	registry := &FieldRegistry{
		fields: cfg.fields,
	}
	host := registerStringField(registry)
	port := registerIntegerField(registry)
	if err := cfg.LoadFields([]string{"appName", "-h", "app.com", "-p", "6000"}); err != nil {
		t.Fatalf("Error loading fields: %s", err.Error())
	}

	assertMetadata(t, cfg, host, "app.com", CommandLine)
	assertMetadata(t, cfg, port, 6000, CommandLine)
}
