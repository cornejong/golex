package golex

import "fmt"

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
