package golex

import "fmt"

type Position struct {
	Row int
	Col int
}

func (p Position) String() string {
	return fmt.Sprintf("%3d:%4d", p.Row, p.Col)
}
