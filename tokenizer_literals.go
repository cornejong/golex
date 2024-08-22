package golex

import (
	"slices"
)

var (
	cachedLiteralToken *Token
)

type LiteralTokenizerCacheKey string

type LiteralTokenizer struct{}

func (t LiteralTokenizer) CanTokenize(l *Lexer) bool {
	pos := l.GetPosition()
	for _, literal := range l.LiteralTokens {
		if l.NextCharsAre([]rune(literal.Literal)) {
			cachedLiteralToken = &Token{
				Type:     literal.Type,
				Literal:  literal.Literal,
				Position: pos,
			}
			return true
		}
	}

	return false
}

func (t LiteralTokenizer) Tokenize(l *Lexer) Token {
	if cachedLiteralToken != nil {
		token := *cachedLiteralToken
		cachedLiteralToken = nil
		return token
	}

	return Token{
		Type:     TypeInvalid,
		Literal:  string(l.CharAtCursor()),
		Position: l.GetPosition(),
	}
}

// ###################################################
// #              Utils
// ###################################################

func SortLiteralTokens(tokens []LiteralToken) []LiteralToken {
	slices.SortFunc[[]LiteralToken, LiteralToken](tokens, func(a, b LiteralToken) int {
		aLen := len(a.Literal)
		bLen := len(b.Literal)

		if aLen < bLen {
			return 1
		}

		if aLen > bLen {
			return -1
		}

		return 0
	})

	return tokens
}

type LiteralToken struct {
	Type    TokenType
	Literal string
}
