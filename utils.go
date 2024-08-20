package golex

import (
	"fmt"
	"strings"
)

func expandCharacterPattern(pattern string) (string, error) {
	var validChars strings.Builder
	length := len(pattern)

	for i := 0; i < length; i++ {
		if i+2 < length && pattern[i+1] == '-' {
			start := pattern[i]
			end := pattern[i+2]

			// Validate that 'start' is less than 'end' and both are letters or digits
			if start > end || !(isSameType(start, end)) {
				return "", fmt.Errorf("invalid pattern range: %c-%c", start, end)
			}

			// Expand the range from start to end
			for c := start; c <= end; c++ {
				validChars.WriteRune(rune(c))
			}
			i += 2 // Skip over the next two characters as they are part of the range
		} else {
			// If it's not a range, just add the character
			validChars.WriteByte(pattern[i])
		}
	}

	return validChars.String(), nil
}

func isSameType(a, b byte) bool {
	return (isLetter(a) && isLetter(b)) || (isDigit(a) && isDigit(b))
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
