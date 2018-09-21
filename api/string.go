package api

import (
	"fmt"
	"strings"
)

// Escape the given values, wrapping it in quotes and escaping and
// quotes so that Kakoune parses it as a single argument.
func Escape(v ...interface{}) string {
	// TODO(leeola): figure out the fastest way to print the v...
	// as if Sprintln did it, but WITHOUT the newline at the end.
	//
	// Apparently fmt.Sprint() and fmt.Sprintln() have different
	// behavior, Sprint puts spaces between arguments only if they're
	// not strings.. so i can't use Sprint() ...sadface.
	s := fmt.Sprintln(v...)
	l := len(s)
	s = s[:l-1]
	return escapeString(s)
}

// escapeString escapes the given string.
func escapeString(s string) string {
	surroundWithQuotes := strings.IndexByte(s, ' ') != -1
	s = strings.Replace(s, `"`, `""`, -1)

	if !surroundWithQuotes {
		return s
	}

	return `"` + s + `"`
}

func EscapeString(s string) string {
	return escapeString(s)
}
