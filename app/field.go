package app

import (
	"strconv"
	"strings"
)

// FieldSpecifier is the enumerated source type that was responsible for specifying the Field's value
type FieldSpecifier uint8

const (
	None FieldSpecifier = iota
	EnvironmentVar
	CommandLine
	File
)

// String Stringer interface
func (s FieldSpecifier) String() string {
	switch s {
	case None:
		return "none"
	case EnvironmentVar:
		return "environment-var"
	case CommandLine:
		return "command-line"
	case File:
		return "config-file"
	default:
		return ""
	}
}

// ValueType is the datatype of the Fields value
type ValueType uint8

const (
	Int ValueType = iota
	Uint
	Float
	String
	Bool
)

// String Stringer interface for a ValueType
func (vt ValueType) String() string {
	switch vt {
	case Int:
		return "integer"
	case Uint:
		return "unsigned-integer"
	case Float:
		return "float"
	case String:
		return "string"
	case Bool:
		return "boolean"
	default:
		return ""
	}
}

// ToString converts the value to a string
func (vt ValueType) ToString(value any) string {
	switch value.(type) {
	case float64:
		return strconv.FormatFloat(value.(float64), 'f', 10, 64)
	case int:
		return strconv.Itoa(value.(int))
	case uint:
		return strconv.Itoa(value.(int))
	case bool:
		return strconv.FormatBool(value.(bool))
	case string:
		return value.(string)
	default:
		return ""
	}
}

// Field defines a general field for configuration
type Field struct {
	ShortDescription  string
	LongDescription   string
	Name              string
	ArgName           string
	EnvVar            string
	ConfigSectionName string
	ConfigFieldName   string
	DefaultValue      any
	Required          bool
	Type              ValueType
}

// ValueMetadata contains the resolved value of a Field as well as some extraction metadata
type ValueMetadata struct {
	Value     any
	Specifier FieldSpecifier
	Field     *Field
}

func newValue(field *Field, raw string, specifier FieldSpecifier) (*ValueMetadata, Error) {

	var formatted any
	var err error = nil
	switch field.Type {
	case Bool:
		formatted, err = strconv.ParseBool(raw)
	case Int:
		// int is 32 bit
		formatted, err = strconv.Atoi(raw)
	case Uint:
		u64, err := strconv.ParseUint(raw, 10, 32)
		if err == nil {
			formatted = uint(u64)
		}
	case String:
		formatted = strings.TrimSpace(raw)
	case Float:
		formatted, err = strconv.ParseFloat(raw, 64)

	default:
		return nil, NewInternalError("Unknown field type specified '%T'", field.Type)
	}
	if err != nil {
		return nil, BuildInternalError().Cause(err).
			Str("rawValue", raw).
			Str("specifier", specifier.String()).
			Msg("error formatting configuration field")

	}

	return &ValueMetadata{Value: formatted, Specifier: specifier, Field: field}, nil
}
