package golex

import "fmt"

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

func (t Token) Is(tt TokenType) bool {
	return t.Type == tt
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
