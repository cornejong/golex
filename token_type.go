package golex

type TokenType interface {
	String() string
}

type Type string

func (t Type) String() string {
	return string(t)
}
