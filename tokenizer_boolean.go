package golex

type BooleanTokenizer struct{}

func (b BooleanTokenizer) CanTokenize(l *Lexer) bool {
	return l.NextCharsAre([]rune("true")) || l.NextCharsAre([]rune("false"))
}

func (b BooleanTokenizer) Tokenize(l *Lexer) (Token, error) {
	token := Token{Type: TypeBool, Position: l.GetPosition()}

	if l.NextCharsAre([]rune("true")) {
		token.Literal = "true"
		token.Value = true
		l.IncrementCursor(3)
		return token, nil
	}

	if l.NextCharsAre([]rune("false")) {
		token.Literal = "false"
		token.Value = false
		l.IncrementCursor(4)
		return token, nil
	}

	// Should be unreachable if CanTokenize is called first to check...
	token.Type = TypeInvalid
	token.Literal = string(l.CharAtCursor())

	return token, NewError("Untokenizable boolean", token.Position, l.GetCursor(), l.state.Content)
}
