package golex

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type NumberTokenizer struct{}

func (n NumberTokenizer) CanTokenize(l *Lexer) bool {
	return unicode.IsNumber(l.CharAtCursor()) || (l.CharAtCursor() == '-' && unicode.IsNumber(l.CharAtRelativePosition(1)))
}

func (n NumberTokenizer) Tokenize(l *Lexer) (Token, error) {
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

	if token.Type == TypeFloat {
		if strings.HasSuffix(token.Literal, ".") {
			return token, NewError(fmt.Sprintf("Malformed float '%s'. Missing Decimal places.", token.Literal), token.Position, l.state.Content)
		}

		decimalSeparatorCount := strings.Count(token.Literal, ".")
		if decimalSeparatorCount > 1 {
			return token, NewError(fmt.Sprintf("Malformed float '%s'. To many decimal separators. Expect 1 but got %d", token.Literal, decimalSeparatorCount), token.Position, l.state.Content)
		}

		// TODO: Make a lexer option to enable number parsing errors
		// TODO: For now just ignore them, moslty a convinence feature..
		token.Value, _ = strconv.ParseFloat(token.Literal, 64)
	}

	// TODO: Make a lexer option to enable number parsing errors
	// TODO: For now just ignore them, moslty a convinence feature..
	if token.Type == TypeInteger {
		token.Value, _ = strconv.Atoi(token.Literal)
	}

	return token, nil
}
