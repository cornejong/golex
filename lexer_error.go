package golex

import (
	"fmt"
	"strings"
)

// LexerError represents an error that occurred during lexical analysis.
type LexerError struct {
	Message  string
	Position Position
	Cursor   int
	Snippet  string
}

// Error implements the error interface for LexerError
func (e *LexerError) Error() string {
	return fmt.Sprintf("Lexer error at line %d, column %d: %s\n%s", e.Position.Row, e.Position.Col, e.Message, e.formatSnippet())
}

// Utility to create a new LexerError with a snippet from the input.
func NewLexerError(message string, position Position, cursor int, input []rune) *LexerError {
	snippet := extractSnippet(input, cursor)
	return &LexerError{
		Message:  message,
		Position: position,
		Cursor:   cursor,
		Snippet:  snippet,
	}
}

// Extracts a snippet of the input around the error position for better context
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
func (e *LexerError) formatSnippet() string {
	caretPosition := strings.Repeat(" ", e.Position.Col-1) + "^"
	return fmt.Sprintf("    %s\n    %s", e.Snippet, caretPosition)
}
