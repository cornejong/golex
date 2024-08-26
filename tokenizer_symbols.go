package golex

import (
	"slices"
	"strings"
)

type SymbolTokenizer struct{}

func (s SymbolTokenizer) CanTokenize(l *Lexer) bool {
	return strings.Contains(l.SymbolStartCharacterMap, string(l.CharAtCursor()))
}

func (s SymbolTokenizer) Tokenize(l *Lexer) (Token, error) {
	token := Token{Type: TypeSymbol, Position: l.GetPosition()}

	for !l.CursorIsOutOfBounds() {
		if !strings.Contains(l.SymbolContinueCharacterMap, string(l.CharAtCursor())) {
			break
		}

		token.AppendChar(l.CharAtCursor())
		l.IncrementCursor(1)
	}

	// TODO: Refactor this func to be lookahead based to remove this hacky backtrack
	l.IncrementCursor(-1) // reset to the character we couldn't tokenize

	if l.CheckForKeywords && slices.Contains(l.Keywords, token.Literal) {
		token.Type = TypeKeyword
	}

	return token, nil
}
