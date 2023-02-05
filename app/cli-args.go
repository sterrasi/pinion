package app

import (
	"strings"
)

// CLIArgs stores parsed command line args
type CLIArgs struct {
	fieldValues   map[string]string
	anonymousArgs []string
}

// parses POSIX-like arguments: -arg1 arg1Value -arg2 arg2Value anonymousArg
func parseArgs(args []string, fieldsByArgName map[string]*Field) (*CLIArgs, Error) {

	cliArgs := &CLIArgs{
		fieldValues:   make(map[string]string),
		anonymousArgs: make([]string, 0, 0),
	}
	var field *Field

	for n := 1; n < len(args); n++ {

		var arg string
		isArgName := strings.HasPrefix(args[n], "-")
		if isArgName {
			// strip the dash off
			arg = args[n][1:len(args[n])]
		} else {
			arg = args[n]
		}

		// if a field was parsed prior then assign it to a new argument
		if field != nil {
			cliArgs.fieldValues[field.Name] = arg
			field = nil
			continue
		}

		// check for an anonymous arg
		if !isArgName {
			cliArgs.anonymousArgs = append(cliArgs.anonymousArgs, arg)
			continue
		}

		match, present := fieldsByArgName[arg]
		if !present {
			return nil, BuildSysConfigError().Str("argument", arg).
				Msg("Unknown command line argument")
		}

		// check for a flag that does not require another argument
		if match.Type == Bool {
			cliArgs.fieldValues[arg] = "true"
			continue
		} else {
			field = match
		}
	}
	return cliArgs, nil
}
