package golex

type BooleanTokenizer struct{}

func NewBooleanTokenizer() BooleanTokenizer {
	return BooleanTokenizer{}
}

func (b BooleanTokenizer) CanTokenize(l *Lexer) bool {
	return l.NextCharsAre([]rune("true")) || l.NextCharsAre([]rune("false"))
}

func (b BooleanTokenizer) Tokenize(l *Lexer) Token {
	token := Token{Type: TypeBool, Position: l.GetPosition()}

	if l.NextCharsAre([]rune("true")) {
		token.Literal = "true"
		token.Value = true
		l.IncrementCursor(3)
		return token
	}

	if l.NextCharsAre([]rune("false")) {
		token.Literal = "false"
		token.Value = false
		l.IncrementCursor(4)
		return token
	}

	token.Type = TypeInvalid
	token.Literal = string(l.CharAtCursor())

	return token
}
