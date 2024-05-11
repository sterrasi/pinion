package pinion

import (
	"golang.org/x/exp/constraints"
	"math"
	"strings"
)

type Number interface {
	constraints.Integer | constraints.Float
}

// Normalize will convert the given string to lower case and trim it's whitespace.
func Normalize(value string) string {
	if value == "" {
		return ""
	}
	nml := strings.ToLower(value)
	return strings.TrimSpace(nml)
}

// GetMapValues returns a slice of the specified maps values
func GetMapValues[K comparable, V comparable](m map[K]V) []V {
	r := make([]V, 0, len(m))
	for _, val := range m {
		r = append(r, val)
	}
	return r
}

// ToFloat64 casts the given value to a float64 and returns a pointer to it. If the value is not numeric then
// nil is returned
func ToFloat64(value any) *float64 {
	var result float64

	switch value.(type) {
	case float32:
		result = float64(value.(float32))
	case float64:
		result = value.(float64)
	case int8:
		result = float64(value.(int8))
	case int16:
		result = float64(value.(int16))
	case int32:
		result = float64(value.(int32))
	case int64:
		result = float64(value.(int64))
	case uint8:
		result = float64(value.(uint8))
	case uint16:
		result = float64(value.(uint16))
	case uint32:
		result = float64(value.(uint32))
	case uint64:
		result = float64(value.(uint64))
	case uintptr:
		result = float64(value.(uintptr))

	default:
		return nil
	}
	return &result
}

// IsWithinRange returns true if the two values fall within the given rangeVal
func IsWithinRange(v1 any, v2 any, rangeVal float64) bool {
	f1 := ToFloat64(v1)
	f2 := ToFloat64(v2)

	if f1 == nil || f2 == nil {
		return false
	}
	return math.Abs(*f1-*f2) < rangeVal
}
