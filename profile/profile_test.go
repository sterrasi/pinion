package profile

import (
	"fmt"
	"github.com/sterrasi/pinion/app"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestProfile_ProductionProfileValue(t *testing.T) {
	var prod Profile = Production

	s := fmt.Sprintf("%v", prod)

	assert.Equal(t, "production", s)
	assert.Equal(t, "production", prod.String())
}

func TestProfile_DevelopmentProfileValue(t *testing.T) {
	var dev Profile = Development

	s := fmt.Sprintf("%v", dev)

	assert.Equal(t, "development", s)
	assert.Equal(t, "development", dev.String())
}

func TestProfile_TestProfileValue(t *testing.T) {
	var test Profile = Test

	s := fmt.Sprintf("%v", test)

	assert.Equal(t, "test", s)
	assert.Equal(t, "test", test.String())
}

func TestProfile_CanOverrideProfile(t *testing.T) {
	t.Cleanup(func() {
		OverrideProfile(Production)
	})
	assert.Equal(t, GetActiveProfile(), Production)
	OverrideProfile(Development)
	assert.Equal(t, GetActiveProfile(), Development)
}

func TestProfile_CanOverrideProfileWithValue(t *testing.T) {
	t.Cleanup(func() {
		OverrideProfile(Production)
	})

	override := "  	DEvelopMent   "

	assert.Equal(t, GetActiveProfile(), Production)
	assert.NoError(t, SetOptionalProfileOverride(&override))
	assert.Equal(t, GetActiveProfile(), Development)
}

func TestProfile_fetchSetProfile(t *testing.T) {
	p, err := fetchProfileWithoutSideEffects(Development, "prODuctIon ")
	assert.NoError(t, err)
	assert.Equal(t, Production, p)
}

func TestProfile_fetchInvalidProfile(t *testing.T) {
	_, err := fetchProfileWithoutSideEffects(Development, "prODu")
	assert.Error(t, err)
	assert.Equal(t, err.(app.Error).Code(), app.SystemConfigurationErrorCode)
}

func TestProfile_fetchDefaultProfile(t *testing.T) {
	p, err := fetchProfileWithoutSideEffects(Development, "")
	assert.NoError(t, err)
	assert.Equal(t, Development, p)
}

func fetchProfileWithoutSideEffects(defaultProfile Profile, envVarValue string) (Profile, error) {
	orig := os.Getenv(envVar)
	_ = os.Setenv(envVar, envVarValue)
	p, err := fetchProfile(defaultProfile)
	_ = os.Setenv(envVar, orig)
	return p, err
}
