package profile

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/sterrasi/pinion/app"
	"os"
)

const envVar = "ACTIVE_PROFILE"

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
var profiles = [3]Profile{Production, Development, Test}

// the Profile that the application is running under
var activeProfile Profile

// GetActiveProfile returns the Profile that the application is running under
func GetActiveProfile() Profile {
	return activeProfile
}

// SetOptionalProfileOverride allows for the active profile to be overridden in out of band cases like unit tests.
func SetOptionalProfileOverride(value *string) app.Error {
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

// init will attempt to fetch the active profile from the environment. If no active profile is
// specified then 'Production' will be used by default.
func init() {
	p, err := fetchProfile(Production)
	if err != nil {
		log.Error().Msg("Cannot start the application under an unknown profile")
		os.Exit(2)
	}
	activeProfile = p
}

// fetchProfile retrieves this application's Profile from the activeProfileEnvVar environment variable.
//   - If the Profile value is not recognizable then an error is returned
//   - If the Profile value is not provided then the defaultProfile is returned
func fetchProfile(defaultProfile Profile) (Profile, app.Error) {
	nml := app.Normalize(os.Getenv(envVar))

	// use the default profile if not specified
	if nml == "" {
		return defaultProfile, nil
	}

	// resolve the profile value against the array of known profiles
	for _, profile := range profiles {
		if profile.String() == nml {
			return profile, nil
		}
	}

	// profile was not recognized
	msg := fmt.Sprintf("invalid active profile value '%s'", os.Getenv(envVar))
	log.Error().Msg(msg)
	return 0, app.NewSysConfigError(msg)
}

// parseProfile will parse the provided string value into a Profile enum
func parseProfile(value string) (*Profile, app.Error) {
	nml := app.Normalize(value)
	if nml == "" {
		return nil, app.BuildIllegalArgumentError().Context("ParseProfile").
			Msg("Cannot parse blank string into a Profile")
	}

	// resolve the Profile value against the array of known profiles
	for _, ft := range profiles {
		if ft.String() == nml {
			return &ft, nil
		}
	}

	// profile was not recognized
	msg := fmt.Sprintf("invalid active profile value '%s'", nml)
	log.Error().Msg(msg)
	return nil, app.BuildIllegalArgumentError().Context("ParseProfile").Msg(msg)
}
