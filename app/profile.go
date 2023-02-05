package app

import (
	"errors"
	"fmt"
	"github.com/sterrasi/pinion"
	"os"
)

const profileEnvVar = "ACTIVE_PROFILE"

// Profile enum
type Profile uint8

const (
	Production Profile = iota + 1
	Development
	Test
)

// String stringer interface
func (p Profile) String() string {
	switch p {
	case Production:
		return "production"
	case Development:
		return "development"
	case Test:
		return "test"
	default:
		return ""
	}
}

// array of all profiles
var allProfiles = [3]Profile{Production, Development, Test}

// the Profile that the application is running under
var activeProfile = Production

// GetActiveProfile returns the Profile that the application is running under
func GetActiveProfile() Profile {
	return activeProfile
}

// OverrideProfileValue allows for the active profile to be overridden in out of band cases like unit tests.
func OverrideProfileValue(value *string) error {
	if value == nil {
		return nil
	}
	pr, err := parseProfile(*value)
	if err != nil {
		return err
	}
	OverrideProfile(*pr)
	return nil
}

// OverrideProfile allows for active profile to be overridden
func OverrideProfile(p Profile) {
	activeProfile = p
}

// LoadProfile will attempt to fetch the active profile from the environment. If no active profile is
// specified then the defaultProfile will be used.
func LoadProfile(defaultProfile Profile) error {
	p, err := fetchProfile(defaultProfile)
	if err != nil {
		return BuildSysConfigError().Cause(err).
			Msg("Cannot start the application under an unknown profile")
	}
	activeProfile = *p
	return nil
}

// fetchProfile retrieves this application's Profile from the activeProfileEnvVar environment variable.
//   - If the Profile value is not recognizable then an error is returned
//   - If the Profile value is not provided then the defaultProfile is returned
func fetchProfile(defaultProfile Profile) (*Profile, error) {
	nml := pinion.Normalize(os.Getenv(profileEnvVar))

	// use the default profile if not specified
	if nml == "" {
		return &defaultProfile, nil
	}

	// resolve the profile value against the array of known profiles
	for _, profile := range allProfiles {
		if profile.String() == nml {
			return &profile, nil
		}
	}

	// profile was not recognized
	msg := fmt.Sprintf("invalid active profile value '%s'", os.Getenv(profileEnvVar))
	return nil, errors.New(msg)
}

// parseProfile will parse the provided string value into a Profile enum
func parseProfile(value string) (*Profile, Error) {
	nml := pinion.Normalize(value)
	if nml == "" {
		return nil, BuildIllegalArgumentError().Context("ParseProfile").
			Msg("Cannot parse blank string into a Profile")
	}

	// resolve the Profile value against the array of known profiles
	for _, ft := range allProfiles {
		if ft.String() == nml {
			return &ft, nil
		}
	}

	// profile was not recognized
	return nil, BuildIllegalArgumentError().Context("ParseProfile").
		Msgf("invalid active profile value '%s'", value)
}
