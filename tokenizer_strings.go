package golex

type StringTokenizerCacheKey string

type StringTokenizer struct {
}

func NewStringTokenizer() StringTokenizer {
	return StringTokenizer{}
}

func (s StringTokenizer) CanTokenize(l *Lexer) bool {
	for _, enclosure := range l.StringEnclosures {
		if l.NextCharsAre([]rune(enclosure.Enclosure)) {
			l.state.Cache[StringTokenizerCacheKey("enclosure")] = enclosure
			return true
		}
	}

	return false
}

func (s StringTokenizer) Tokenize(l *Lexer) Token {
	var enclosure StringEnclosure
	if enclosureI, ok := l.state.Cache[StringTokenizerCacheKey("enclosure")]; ok {
		enclosure = enclosureI.(StringEnclosure)
	} else {
		if !s.CanTokenize(l) {
			return Token{Type: TypeInvalid, Position: l.GetPosition()}
		} else {
			return s.Tokenize(l)
		}
	}

	return enclosure.Tokenize(l)
}

// ###################################################
// #              StringEnclosures
// ###################################################

type StringEnclosure struct {
	Type      TokenType
	Enclosure string
	Escapable bool
}

func (se StringEnclosure) Tokenize(l *Lexer) Token {
	if len(se.Enclosure) > 1 {
		return se.TokenizeNotEscapableMultiChar(l)
	}

	if !se.Escapable {
		return se.TokenizeNotEscapableSingleChar(l)
	}

	return se.TokenizeEscapable(l)
}

func (se StringEnclosure) TokenizeEscapable(l *Lexer) Token {
	enclosureChar := []rune(se.Enclosure)[0]
	token := Token{Type: se.Type, Position: l.GetPosition()}
	start := l.GetCursor()

	token.AppendChar(l.CharAtCursor())
	l.IncrementCursor(1)

	nextEnclosureCharIsEscaped := false
	for !l.CursorIsOutOfBounds() && (l.CharAtCursor() != enclosureChar || nextEnclosureCharIsEscaped) {
		if nextEnclosureCharIsEscaped {
			nextEnclosureCharIsEscaped = false
		}

		if l.CharAtCursor() == '\\' && l.CharAtRelativePosition(1) == enclosureChar {
			nextEnclosureCharIsEscaped = true
		}

		token.AppendChar(l.CharAtCursor())
		l.IncrementCursor(1)
	}

	token.AppendChar(l.CharAtCursor())

	token.Value = l.GetSourceSubsString(start+1, l.GetCursor())

	return token
}

func (se StringEnclosure) TokenizeNotEscapableSingleChar(l *Lexer) Token {
	enclosureChar := []rune(se.Enclosure)[0]
	token := Token{Type: se.Type, Position: l.GetPosition()}
	start := l.GetCursor()

	token.AppendChar(l.CharAtCursor())
	l.IncrementCursor(1)

	for !l.CursorIsOutOfBounds() && l.CharAtCursor() != enclosureChar {
		token.AppendChar(l.CharAtCursor())
		l.IncrementCursor(1)
	}

	token.AppendChar(l.CharAtCursor())
	token.Value = l.GetSourceSubsString(start+1, l.GetCursor())

	return token
}

func (se StringEnclosure) TokenizeNotEscapableMultiChar(l *Lexer) Token {
	enclosureLen := len(se.Enclosure)
	token := Token{Type: se.Type, Position: l.GetPosition()}
	start := l.GetCursor()

	token.AppendChar([]rune(se.Enclosure)...)
	l.IncrementCursor(enclosureLen)

	for !l.CursorIsOutOfBounds() && !l.NextCharsAre([]rune(se.Enclosure)) {
		token.AppendChar(l.CharAtCursor())
		l.IncrementCursor(1)
	}

	token.AppendChar([]rune(se.Enclosure)...)

	token.Value = l.GetSourceSubsString(start+enclosureLen, l.GetCursor())
	l.IncrementCursor(enclosureLen - 1)

	return token
}

var (
	DoubleQuoteStringEnclosure StringEnclosure = StringEnclosure{
		Type:      TypeDoubleQuoteString,
		Enclosure: "\"",
		Escapable: true,
	}
	SingleQuoteStringEnclosure StringEnclosure = StringEnclosure{
		Type:      TypeSingleQuoteString,
		Enclosure: "'",
	}
	BacktickStringEnclosure StringEnclosure = StringEnclosure{
		Type:      TypeBacktickString,
		Enclosure: "`",
	}
	TripleBacktickStringEnclosure StringEnclosure = StringEnclosure{
		Type:      TypeTripleBacktickString,
		Enclosure: "```",
	}
)
