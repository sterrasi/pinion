package app

import "strings"

// Normalize will convert the given string to lower case and trim it's whitespace.
func Normalize(value string) string {
	if value == "" {
		return ""
	}
	nml := strings.ToLower(value)
	nml = strings.TrimSpace(nml)
	if nml == "" {
		return ""
	}
	return nml
}
