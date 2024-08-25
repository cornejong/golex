package golex

import (
	"fmt"
	"iter"
)

// Tokens represents a set of tokens
type Tokens []Token

// TokenCollection represents an iterable collection of tokens
type TokenCollection struct {
	tokens       Tokens
	tokensLength int
	cursor       int
}

func NewTokenCollection(tokens Tokens) TokenCollection {
	return TokenCollection{tokens: tokens, tokensLength: len(tokens), cursor: 0}
}

func (ti *TokenCollection) Iter() iter.Seq2[int, Token] {
	return func(yield func(int, Token) bool) {
		for ti.cursor = 0; ti.cursor < ti.tokensLength; ti.cursor++ {
			if !yield(ti.cursor, ti.tokens[ti.cursor]) {
				return
			}
		}
	}
}

// IncrementCursor increments the cursor by the amount
func (t *TokenCollection) IncrementCursor(amount int) {
	t.cursor += amount
}

func (t TokenCollection) CursorIsOutOfBounds() bool {
	return t.cursor >= t.tokensLength
}

func (t TokenCollection) ReachedEOF() bool {
	return t.TokenAtCursor().Type == TypeEof
}

// NextToken increments the cursor position by 1
// and returns the token at that position
func (t *TokenCollection) NextToken() Token {
	t.cursor += 1
	return t.TokenAtPosition(t.cursor)
}

// TokenAtCursor returns the token at the current cursor position
func (t *TokenCollection) TokenAtCursor() Token {
	return t.TokenAtPosition(t.cursor)
}

// TokenAtPosition returns the token at the absolute position
func (t *TokenCollection) TokenAtPosition(pos int) Token {
	// TODO: add out of bounds checking
	return t.tokens[pos]
}

// TokenAtRelativePosition returns the token at the position relative to the cursor
func (t *TokenCollection) TokenAtRelativePosition(pos int) Token {
	// TODO: add out of bounds checking
	return t.tokens[t.cursor+pos]
}

// CollectTokensBetweenParentheses collects all the tokens between
// the opening and closing parentheses starting at the offset start
// It assumes that the offset start contains the opening parenthesis.
// Nested parenthesis will be contained in the output tokens until
// the matching closing parenthesis to the start is found.
func (t *TokenCollection) CollectTokensBetweenParentheses() (Tokens, int, int, error) {
	return t.CollectTokensBetween(TypeOpenParen, TypeCloseParen)
}

// CollectTokensBetweenCurlyBraced collects all the tokens between
// the opening and closing curly braces starting at the offset start
// It assumes that the offset start contains the opening curly brace.
// Nested curly braces will be contained in the output tokens until
// the matching closing curly brace to the start is found.
func (t *TokenCollection) CollectTokensBetweenCurlyBraces() (Tokens, int, int, error) {
	return t.CollectTokensBetween(TypeOpenCurly, TypeCloseCurly)
}

// CollectTokensBetween collects all the tokens between the open
// and close type starting at the offset start
// it assumes that the offset start contains the opening token.
// Nested openers and closers will be contained in the output
// tokens until the matching closer is found.
func (t *TokenCollection) CollectTokensBetween(open TokenType, close TokenType) (Tokens, int, int, error) {
	collected := Tokens{}
	token := t.TokenAtCursor()

	if !token.Is(open) {
		return collected, -1, -1, fmt.Errorf("Token at start offset is not of opener type %s", open)
	}

	start := t.cursor
	end := start
	level := 1

	for !token.Is(TypeEof) {
		end = t.cursor
		t.cursor += 1
		token = t.tokens[t.cursor]

		if token.Is(TypeEof) {
			return collected, start, end, fmt.Errorf("Unexpected EndOfFile")
		}

		if token.Is(close) {
			level -= 1
			if level == 0 {
				break
			}
		}

		if token.Is(open) {
			level += 1
		}

		collected = append(collected, token)
	}

	return collected, start, end, nil
}

func (t *TokenCollection) CollectTokensDelimited(tokenType TokenType, delimiter TokenType) (Tokens, error) {
	tokens := Tokens{}

	token := t.TokenAtCursor()
	for !token.Is(TypeEof) {
		if !token.Is(tokenType) {
			return tokens, fmt.Errorf("expected %s but found %s", tokenType, token.Type)
		}

		tokens = append(tokens, token)

		if !t.TokenAtRelativePosition(1).Is(delimiter) {
			break
		}

		t.IncrementCursor(2) // Consume the delimiter and get the next token after it
		token = t.TokenAtCursor()
	}

	return tokens, nil
}

func (t *TokenCollection) CollectAnyTokensDelimited(delimiter TokenType) ([]Token, error) {
	tokens := []Token{}

	token := t.TokenAtCursor()
	for !token.Is(TypeEof) {
		tokens = append(tokens, token)

		if !t.TokenAtRelativePosition(1).Is(delimiter) {
			break
		}

		t.IncrementCursor(2) // Consume the delimiter and get the next token after it
		token = t.TokenAtCursor()
	}

	return tokens, nil
}

func (t *TokenCollection) CollectTokensUntil(delimiter TokenType) ([]Token, error) {
	tokens := []Token{}

	token := t.TokenAtCursor()
	for !token.Is(TypeEof) && !token.Is(delimiter) {
		tokens = append(tokens, token)

		if !t.TokenAtRelativePosition(1).Is(delimiter) {
			break
		}

		token = t.NextToken()
	}

	return tokens, nil
}
