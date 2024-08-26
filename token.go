package golex

import (
	"fmt"
)

// ###################################################
// #                    Token
// ###################################################
type Token struct {
	Type     TokenType
	Literal  string
	Value    any
	Position Position
}

func (t *Token) AppendChar(char ...rune) {
	t.Literal += string(char)
}

func (t Token) Dump() {
	fmt.Printf("%s -> %-22s%-22s(%v)\n", t.Position.String(), t.Type.String(), t.Literal, t.Value)
}

func (t Token) Is(token Token) bool {
	if token.Literal != "" && t.Literal != token.Literal {
		return false
	}

	return token.Type == AnyTokenType || t.Type == token.Type
}

func (t Token) IsAnyOf(tokens ...Token) bool {
	for _, token := range tokens {
		if t.Is(token) {
			return true
		}
	}

	return false
}

func (t Token) TypeIs(tt TokenType) bool {
	return t.Type == tt
}

func (t Token) TypeIsAnyOf(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if t.Type == tokenType {
			return true
		}
	}

	return false
}

func (t Token) LiteralIs(literal string) bool {
	return t.Literal == literal
}

func (t Token) LiteralIsAnyOf(literals ...string) bool {
	for _, literal := range literals {
		if t.Literal == literal {
			return true
		}
	}

	return false
}

// ###################################################
// #                   TokenType
// ###################################################
type TokenType interface {
	String() string
}

// ###################################################
// #                   Position
// ###################################################
type Position struct {
	Row int
	Col int
}

func (p Position) String() string {
	return fmt.Sprintf("%3d:%4d", p.Row, p.Col)
}
