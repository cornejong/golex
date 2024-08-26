package golex

import (
	"fmt"
	"strings"
)

// LexerError represents an error that occurred during lexical analysis.
type Error struct {
	Message  string
	Position Position
	Cursor   int
	Snippet  string
}

// Error implements the error interface for LexerError
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s\n%s", e.Position.String(), e.Message, e.formatSnippet())
}

func NewError(message string, position Position, cursor int, input []rune) *Error {
	snippet := extractSnippet(input, cursor)
	return &Error{
		Message:  message,
		Position: position,
		Cursor:   cursor,
		Snippet:  snippet,
	}
}

func extractSnippet(input []rune, cursor int) string {
	const contextLength = 20
	start := cursor - contextLength
	if start < 0 {
		start = 0
	}
	end := cursor + contextLength
	if end > len(input) {
		end = len(input)
	}
	return string(input[start:end])
}

// Formats the snippet with a caret (^) to indicate the exact error location.
func (e *Error) formatSnippet() string {
	caretPosition := strings.Repeat(" ", e.Position.Col-1) + "^"
	return fmt.Sprintf("    %s\n    %s", e.Snippet, caretPosition)
}
