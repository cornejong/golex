package golex

import (
	"fmt"
	"iter"
	"slices"
	"unicode"
)

var (
	EOF rune = rune(byte(0x03))

	// defaultSymbolCharacterMap         string = "a-zA-Z0-9_"
	defaultSymbolContinueCharacterMapExpanded string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	defaultSymbolStartCharacterMapExpanded    string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
)

type State struct {
	Content       []rune
	ContentLength int
	Cursor        int

	PositionCursor int
	Position       Position

	LineIndexes      []int
	LineIndexesCount int
	CurrentToken     *Token
	LookaheadCache   LookaheadCache
}

type LookaheadCache struct {
	tokens []Token
	count  int
}

func (lc *LookaheadCache) ContainsItems() bool { return lc.count > 0 }
func (lc *LookaheadCache) ItemCount() int      { return lc.count }

func (lc *LookaheadCache) AddItem(token Token) bool {
	if slices.Contains(lc.tokens, token) {
		return false // We already cached this token
	}

	lc.tokens = append(lc.tokens, token)
	lc.count += 1
	return true
}

// TODO: this does not check for out of bounds stuff..
// TODO: Probably not what we want to do...
func (lc *LookaheadCache) PluckItem() Token {
	token := lc.tokens[0]

	lc.tokens = lc.tokens[1:]
	lc.count -= 1

	return token
}

func (lc *LookaheadCache) GetFirstItem() Token {
	return lc.tokens[0]
}

func (lc *LookaheadCache) GetItem(pos int) Token {
	return lc.tokens[pos]
}

func NewState(content string) State {
	c := []rune(content)

	lineIndexes := []int{}
	for i, char := range c {
		if char == '\n' {
			lineIndexes = append(lineIndexes, i)
		}
	}

	return State{
		Content:          append(c, EOF),
		ContentLength:    len(c),
		LineIndexes:      lineIndexes,
		LineIndexesCount: len(lineIndexes),
		PositionCursor:   0,
		Position:         Position{Col: 1, Row: 1, Cursor: 0},
		CurrentToken: &Token{
			Type:     TypeSof,
			Position: Position{},
		},
	}
}

type Lexer struct {
	state State

	CommentTokenizer CommentTokenizer
	LiteralTokenizer LiteralTokenizer
	NumberTokenizer  NumberTokenizer
	BooleanTokenizer BooleanTokenizer
	SymbolTokenizer  SymbolTokenizer
	StringTokenizer  StringTokenizer

	tokenizers        map[TokenizerType]Tokenizer
	tokenizationOrder []TokenizerType

	// options
	LiteralTokens    []LiteralToken
	StringEnclosures []StringEnclosure
	CommentSyntaxes  []CommentSyntax
	Keywords         []string
	IgnoreTokens     []TokenType

	IgnoreWhitespace           bool
	IgnoreComments             bool
	UseBuiltinTypes            bool
	CheckForKeywords           bool
	SymbolStartCharacterMap    string
	SymbolContinueCharacterMap string
	DebugPrintTokens           bool
	OmitTokenPosition          bool
}

func NewLexer(options ...LexerOptionFunc) *Lexer {
	lexer := &Lexer{
		tokenizers: map[TokenizerType]Tokenizer{},
		tokenizationOrder: []TokenizerType{
			TypeCommentTokenizer,
			TypeNumberTokenizer,
			TypeLiteralTokenizer,
			TypeStringTokenizer,
			TypeBooleanTokenizer,
			TypeSymbolTokenizer,
		},

		LiteralTokens:    SortLiteralTokens(buildInLiteralTokens),
		StringEnclosures: []StringEnclosure{SingleQuoteStringEnclosure, DoubleQuoteStringEnclosure},
		CommentSyntaxes:  []CommentSyntax{SlashSingleLineCommentSyntax, SlashMultilineCommentSyntax},

		DebugPrintTokens:           false,
		IgnoreWhitespace:           true,
		IgnoreComments:             false,
		UseBuiltinTypes:            false,
		SymbolStartCharacterMap:    defaultSymbolStartCharacterMapExpanded,
		SymbolContinueCharacterMap: defaultSymbolContinueCharacterMapExpanded,
	}

	// Comment Tokenizer
	lexer.CommentTokenizer = CommentTokenizer{}
	lexer.tokenizers[TypeCommentTokenizer] = &lexer.CommentTokenizer

	// Literal tokenizer
	lexer.LiteralTokenizer = LiteralTokenizer{}
	lexer.tokenizers[TypeLiteralTokenizer] = &lexer.LiteralTokenizer

	// Number tokenizer
	lexer.NumberTokenizer = NumberTokenizer{}
	lexer.tokenizers[TypeNumberTokenizer] = &lexer.NumberTokenizer

	// Boolean tokenizer
	lexer.BooleanTokenizer = BooleanTokenizer{}
	lexer.tokenizers[TypeBooleanTokenizer] = &lexer.BooleanTokenizer

	// String Tokenizer
	lexer.StringTokenizer = StringTokenizer{}
	lexer.tokenizers[TypeStringTokenizer] = &lexer.StringTokenizer

	// Symbol tokenizer
	lexer.SymbolTokenizer = SymbolTokenizer{}
	lexer.tokenizers[TypeSymbolTokenizer] = &lexer.SymbolTokenizer

	for _, opt := range options {
		opt(lexer)
	}

	return lexer
}

func (l *Lexer) RemoveTokenizer(tokenizerType TokenizerType) {
	delete(l.tokenizers, tokenizerType)
}

func (l *Lexer) TokenizeToSlice(content string) ([]Token, error) {
	tokens := []Token{}
	for token, err := range l.Iterate(content) {
		if err != nil {
			return tokens, err
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (l *Lexer) TokenizeManual(content string) {
	l.state = NewState(content)
}

func (l *Lexer) Iterate(content string) iter.Seq2[Token, error] {
	l.state = NewState(content)

	return func(yield func(Token, error) bool) {
		for !l.ReachedEOF() {
			if !yield(l.NextToken()) {
				return
			}
		}
	}
}

func (l *Lexer) SkipWhitespace() {
	for !l.CursorIsOutOfBounds() && unicode.IsSpace(l.CharAtCursor()) {
		l.IncrementCursor(1)
	}
}

// Lookahead returns the token at count offset from the cursor without consuming it
func (l *Lexer) Lookahead(offset int) Token {
	if l.state.LookaheadCache.ItemCount() >= offset {
		return l.state.LookaheadCache.GetItem(offset - 1)
	}

	state := l.GetState()
	var token Token

	var i int
	for i = 0; i < offset; i++ {
		token, _ = l.nextToken()
		state.LookaheadCache.AddItem(token)

		if token.TypeIs(TypeEof) {
			l.SetState(state)
			return token
		}
	}

	l.SetState(state)
	return token
}

// LookaheadIterator returns an iterator that iterates over
// the count number of tokens without consuming them
func (l *Lexer) LookaheadIterator(count int) iter.Seq[Token] {
	l.Lookahead(count)

	return func(yield func(Token) bool) {
		for i := 0; i < count; i++ {
			if !yield(l.state.LookaheadCache.tokens[i]) {
				return
			}
		}

		return
	}
}

func (l *Lexer) NextToken() (Token, error) {
	token, err := l.nextToken()

	if l.DebugPrintTokens {
		token.Dump()
	}

	return token, err
}

func (l *Lexer) nextToken() (Token, error) {
	// check if we have anything in the lookahead cache
	if l.state.LookaheadCache.ContainsItems() {
		return l.state.LookaheadCache.PluckItem(), nil
	}

	if l.IgnoreWhitespace {
		l.SkipWhitespace()
	}

	if l.CursorIsOutOfBounds() {
		l.state.CurrentToken = &Token{
			Type:     TypeEof,
			Literal:  string(EOF),
			Position: l.GetPosition(),
		}

		return *l.state.CurrentToken, nil
	}

	var err error
	token := Token{
		Type:     TypeInvalid,
		Literal:  string(l.CharAtCursor()),
		Position: l.GetPosition(),
	}

	for _, tokenizerType := range l.tokenizationOrder {
		tokenizer, ok := l.tokenizers[tokenizerType]
		if !ok {
			continue
		}

		if tokenizer.CanTokenize(l) {
			token, err = l.tokenizers[tokenizerType].Tokenize(l)
			break
		}
	}

	l.state.CurrentToken = &token
	l.IncrementCursor(1)

	if token.TypeIs(TypeInvalid) && err == nil {
		err = NewError(fmt.Sprintf("Invalid character '%c'", l.CharAtCursor()), token.Position, l.state.Content)
	}

	return token, err
}

func (l *Lexer) GetPosition() Position {
	if l.OmitTokenPosition {
		return Position{}
	}

	if l.state.Cursor == l.state.Position.Cursor {
		return l.state.Position
	}

	l.state.Position.Cursor = l.state.Cursor

	if l.state.LineIndexesCount == 0 {
		l.state.Position.Col += l.state.Cursor - l.state.PositionCursor
		l.state.PositionCursor = l.state.Cursor
		return l.state.Position
	}

	nextLineStart := l.state.LineIndexes[0]
	if l.state.Cursor <= nextLineStart {
		l.state.Position.Col += l.state.Cursor - l.state.PositionCursor
		l.state.PositionCursor = l.state.Cursor
		return l.state.Position
	}

	l.state.Position.Row += 1
	l.state.Position.Col = l.state.Cursor - nextLineStart
	l.state.PositionCursor = l.state.Cursor

	// Remove the consumed entry
	l.state.LineIndexes = l.state.LineIndexes[1:]
	l.state.LineIndexesCount -= 1

	return l.state.Position
}

func (l Lexer) GetCurrentLine() (int, int) {
	// TODO: remove consumed lines
	// TODO: this way we can always return the first entry

	lastStart := 0
	lastIndex := 0
	for _, start := range l.state.LineIndexes {
		if start >= l.state.Cursor {
			return lastIndex, lastStart
		}

		lastStart = start + 1
		lastIndex += 1
	}

	return len(l.state.LineIndexes), lastStart
}

// CharAtCursor returns the rune at the current cursor position
func (l *Lexer) CharAtCursor() rune {
	return l.CharAtPosition(l.state.Cursor)
}

// CharAtRelativePosition returns the rune at the relative position to the cursor
func (l *Lexer) CharAtRelativePosition(pos int) rune {
	return l.CharAtPosition(l.state.Cursor + pos)
}

// CharAtPosition returns the rune at the provided absolute position
func (l *Lexer) CharAtPosition(pos int) rune {
	if l.state.ContentLength <= pos {
		return EOF
	}
	return l.state.Content[pos]
}

// NextCharsAre checks if the next chars from the cursor on match the provided chars without consuming them
func (l *Lexer) NextCharsAre(chars []rune) bool {
	len := len(chars)
	if len == 0 {
		return true
	}

	if l.state.Cursor+len-1 >= l.state.ContentLength {
		return false
	}

	for i := 0; i < len; i++ {
		if chars[i] != l.CharAtRelativePosition(i) {
			return false
		}
	}

	return true
}

// ---------------------------------------------------------------
// Helpers / Getter/Setters
// ---------------------------------------------------------------

func (l *Lexer) GetSourceSubsString(start int, end int) string {
	return string(l.state.Content[start:end])
}

func (l *Lexer) GetState() State {
	return l.state
}

func (l *Lexer) SetState(state State) {
	l.state = state
}

func (l *Lexer) GetCursor() int {
	return l.state.Cursor
}
func (l *Lexer) SetCursor(cursor int) {
	l.state.Cursor = cursor
}

func (l *Lexer) IncrementCursor(amount int) {
	l.state.Cursor += amount
}

func (l *Lexer) CursorIsOutOfBounds() bool {
	return l.state.Cursor >= l.state.ContentLength
}

func (l Lexer) ReachedEOF() bool {
	return l.state.CurrentToken.Type == TypeEof
}

func (l Lexer) CurrentToken() Token {
	return *l.state.CurrentToken
}
