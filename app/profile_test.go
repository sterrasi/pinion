package app

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// the {production} Profile type should be associated with the "production" string value
func TestProfile_ProductionProfileValue(t *testing.T) {
	var prod = Production

	s := fmt.Sprintf("%v", prod)

	assert.Equal(t, "production", s)
	assert.Equal(t, "production", prod.String())
}

// the {development} Profile type should be associated with the "development" string value
func TestProfile_DevelopmentProfileValue(t *testing.T) {
	var dev Profile = Development

	s := fmt.Sprintf("%v", dev)

	assert.Equal(t, "development", s)
	assert.Equal(t, "development", dev.String())
}

// The {test} Profile type should be associated with the "test" string value
func TestProfile_TestProfileValue(t *testing.T) {
	var test Profile = Test

	s := fmt.Sprintf("%v", test)

	assert.Equal(t, "test", s)
	assert.Equal(t, "test", test.String())
}

// Should expect an error when an invalid ACTIVE_PROFILE is set
func TestProfile_fetchInvalidProfile(t *testing.T) {
	_, err := fetchProfileWithoutSideEffects(Development, "prODu")
	assert.Error(t, err)
}

// Should expect the default profile when the ACTIVE_PROFILE is blank
func TestProfile_fetchDefaultProfile(t *testing.T) {
	p, err := fetchProfileWithoutSideEffects(Development, "")
	assert.NoError(t, err)
	assert.Equal(t, Development, p)
}

// The active Profile should be able to be programmatically set
func TestProfile_CanOverrideProfile(t *testing.T) {

	assert.Equal(t, GetActiveProfile(), Production)
	t.Cleanup(func() {
		OverrideProfile(Production)
	})
	OverrideProfile(Development)
	assert.Equal(t, GetActiveProfile(), Development)
}

// The value that is parsed into the active Profile should be able to be set programmatically
// (instead of loading the environment variable)
func TestProfile_CanOverrideProfileWithValue(t *testing.T) {
	t.Cleanup(func() {
		OverrideProfile(Production)
	})

	override := "  	DEvelopMent   "

	assert.Equal(t, GetActiveProfile(), Production)
	assert.NoError(t, OverrideProfileValue(&override))
	assert.Equal(t, GetActiveProfile(), Development)
}

// Test parsing the active Profile from an environment variable
func TestProfile_fetchSetProfile(t *testing.T) {
	p, err := fetchProfileWithoutSideEffects(Development, "prODuctIon ")
	assert.NoError(t, err)
	assert.Equal(t, Production, p)
}

func fetchProfileWithoutSideEffects(defaultProfile Profile, envVarValue string) (Profile, error) {
	orig := os.Getenv(profileEnvVar)
	_ = os.Setenv(profileEnvVar, envVarValue)
	p, err := fetchProfile(defaultProfile)
	_ = os.Setenv(profileEnvVar, orig)
	if p != nil {
		return *p, nil
	}
	return defaultProfile, err
}
