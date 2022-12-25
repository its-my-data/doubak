package util

import (
	"strings"
	"unicode"
)

// MergeSpaces merges all space characters into one regular space.
func MergeSpaces(str *string) string {
	firstSpace := true // Tracking consecutive spaces.

	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// If the rune is the first space, then replace it with a normal space.
			if firstSpace {
				firstSpace = false
				return ' '
			}

			// Not the first space in this section.
			return -1
		}

		// Reset the space counter.
		firstSpace = true

		// Keep it in the string.
		return r
	}, strings.TrimSpace(*str))
}

// StripAllSpaces removes all space characters.
func StripAllSpaces(str *string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			// if the character is a space, drop it
			return -1
		}
		// else keep it in the string
		return r
	}, *str)
}
