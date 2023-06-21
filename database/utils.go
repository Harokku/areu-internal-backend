package database

import "strings"

// Cut slices s around the first instance of sep,
// returning the text before and after sep.
// Modified from strings.Cut to be cas insensitive and cut from the end of the string
// The found result reports whether sep appears in s.
// If sep does not appear in s, cut returns s, "", false.
func iCut(s, sep string) (before, after string, found bool) {
	if i := strings.LastIndex(strings.ToLower(s), strings.ToLower(sep)); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}
