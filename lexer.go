package golex

import (
	"errors"
	"strings"
	"unicode"
)

var (
	EOF rune = rune(byte(0x03))

	// defaultSymbolCharacterMap         string = "a-zA-Z0-9_"
	defaultSymbolContinueCharacterMapExpanded string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	defaultSymbolStartCharacterMapExpanded    string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"

	ErrNoEOFTokenType     error = errors.New("no eof token type specified. Specify an EOF token type or use build-in types")
	ErrNoInvalidTokenType error = errors.New("no eof token type specified. Specify an EOF token type or use build-in types")
)

type State struct {
	Content       []rune
	ContentLength int
	Cursor        int
	Line          int
	LineStart     int
	CurrentToken  *Token
	Cache         map[any]any
}

func NewState(content string) State {
	c := []rune(content)
	return State{
		Content:       append(c, EOF),
		ContentLength: len(c),
		Cache:         make(map[any]any),
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
	lexer.CommentTokenizer = NewCommentTokenizer()
	lexer.tokenizers[TypeCommentTokenizer] = &lexer.CommentTokenizer

	// Literal tokenizer
	lexer.LiteralTokenizer = NewLiteralTokenizer()
	lexer.tokenizers[TypeLiteralTokenizer] = &lexer.LiteralTokenizer

	// Number tokenizer
	lexer.NumberTokenizer = NewNumberTokenizer()
	lexer.tokenizers[TypeNumberTokenizer] = &lexer.NumberTokenizer

	// Boolean tokenizer
	lexer.BooleanTokenizer = NewBooleanTokenizer()
	lexer.tokenizers[TypeBooleanTokenizer] = &lexer.BooleanTokenizer

	// String Tokenizer
	lexer.StringTokenizer = NewStringTokenizer()
	lexer.tokenizers[TypeStringTokenizer] = &lexer.StringTokenizer

	// Symbol tokenizer
	lexer.SymbolTokenizer = NewSymbolTokenizer()
	lexer.tokenizers[TypeSymbolTokenizer] = &lexer.SymbolTokenizer

	for _, opt := range options {
		opt(lexer)
	}

	return lexer
}

func (l *Lexer) AddTokenizer(tokenizerType TokenizerType, tokenizer Tokenizer) {
	l.tokenizers[tokenizerType] = tokenizer
}

func (l *Lexer) RemoveTokenizer(tokenizerType TokenizerType) {
	delete(l.tokenizers, tokenizerType)
}

func (l *Lexer) TokenizeToSlice(content string) []Token {
	l.state = NewState(content)

	tokens := []Token{}
	for l.state.CurrentToken.Type != TypeEof {
		tokens = append(tokens, l.NextToken())
	}

	return tokens
}

func (l *Lexer) TokenizeToChannel(content string, tokens chan Token) {
	l.state = NewState(content)
	for l.state.CurrentToken.Type != TypeEof {
		tokens <- l.NextToken()
	}

	close(tokens)
}

func (l *Lexer) TokenizeToCallback(content string, callback func(Token)) {
	l.state = NewState(content)
	for l.state.CurrentToken.Type != TypeEof {
		callback(l.NextToken())
	}
}

func (l *Lexer) TokenizeManual(content string) {
	l.state = NewState(content)
}

func (l *Lexer) SkipWhitespace() {
	for !l.CursorIsOutOfBounds() && unicode.IsSpace(l.CharAtCursor()) {
		char := l.CharAtCursor()
		l.IncrementCursor(1)

		if char == '\n' {
			l.state.Line += 1
			l.state.LineStart = l.state.Cursor
		}
	}
}

func (l *Lexer) Lookahead(count int) Token {
	state := l.GetState()

	var token Token
	for i := 0; i < count; i++ {
		token = l.nextToken()
		if strings.ContainsRune(token.Literal, EOF) {
			l.SetState(state)
			return token
		}
	}

	l.SetState(state)
	return token
}

func (l *Lexer) NextToken() Token {
	token := l.nextToken()

	if l.DebugPrintTokens {
		token.Dump()
	}

	return token
}

func (l *Lexer) nextToken() Token {
	if l.IgnoreWhitespace {
		l.SkipWhitespace()
	}

	if l.CursorIsOutOfBounds() {
		l.state.CurrentToken = &Token{
			Type:     TypeEof,
			Position: l.GetPosition(),
		}

		return *l.state.CurrentToken
	}

	token := Token{
		Type:     TypeInvalid,
		Position: l.GetPosition(),
	}

	for _, tokenizerType := range l.tokenizationOrder {
		tokenizer, ok := l.tokenizers[tokenizerType]
		if !ok {
			continue
		}

		if tokenizer.CanTokenize(l) {
			token = l.tokenizers[tokenizerType].Tokenize(l)
			break
		}
	}

	l.state.CurrentToken = &token
	l.IncrementCursor(1)

	return token
}

func (l *Lexer) GetSourceSubsString(start int, end int) string {
	return string(l.state.Content[start:end])
}

func (l *Lexer) GetPosition() Position {
	return Position{Row: l.state.Line + 1, Col: l.state.Cursor - l.state.LineStart + 1}
}

func (l *Lexer) GetState() State {
	return l.state
}

func (l *Lexer) SetState(state State) {
	l.state = state
}

func (l *Lexer) CharAtCursor() rune {
	return l.CharAtPosition(l.state.Cursor)
}

func (l *Lexer) CharAtRelativePosition(pos int) rune {
	return l.CharAtPosition(l.state.Cursor + pos)
}

func (l *Lexer) CharAtPosition(pos int) rune {
	if l.state.ContentLength <= pos {
		return EOF
	}
	return l.state.Content[pos]
}

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

func (l *Lexer) GetLine() int {
	return l.state.Line
}
func (l *Lexer) SetLine(line int) {
	l.state.Line = line
}

func (l *Lexer) GetLineStart() int {
	return l.state.LineStart
}
func (l *Lexer) SetLineStart(lineStart int) {
	l.state.LineStart = lineStart
}
