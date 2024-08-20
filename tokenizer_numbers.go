package golex

import (
	"strconv"
	"unicode"
)

type NumberTokenizer struct{}

func NewNumberTokenizer() NumberTokenizer {
	return NumberTokenizer{}
}

func (n NumberTokenizer) CanTokenize(l *Lexer) bool {
	return unicode.IsNumber(l.CharAtCursor()) || (l.CharAtCursor() == '-' && unicode.IsNumber(l.CharAtRelativePosition(1)))
}

func (n NumberTokenizer) Tokenize(l *Lexer) Token {
	token := Token{Type: TypeInteger, Position: l.GetPosition()}

	token.AppendChar(l.CharAtCursor())
	l.IncrementCursor(1)

	for !l.CursorIsOutOfBounds() && (unicode.IsNumber(l.CharAtCursor()) || l.CharAtCursor() == '.') {
		if l.CharAtCursor() == '.' {
			token.Type = TypeFloat
		}

		token.AppendChar(l.CharAtCursor())
		l.IncrementCursor(1)
	}

	l.IncrementCursor(-1)

	if token.Type == TypeInteger {
		token.Value, _ = strconv.Atoi(token.Literal)
	}

	if token.Type == TypeFloat {
		token.Value, _ = strconv.ParseFloat(token.Literal, 64)
	}

	return token
}
