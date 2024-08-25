package golex

import (
	"errors"
	"fmt"
	"slices"
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

	CachedPositionCursor int
	CachedPosition       Position

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
		Content:              append(c, EOF),
		ContentLength:        len(c),
		LineIndexes:          lineIndexes,
		LineIndexesCount:     len(lineIndexes),
		CachedPositionCursor: 0,
		CachedPosition:       Position{Col: 1, Row: 1},
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
		l.IncrementCursor(1)
	}
}

func (l *Lexer) Lookahead(count int) Token {
	if l.state.LookaheadCache.ItemCount() >= count {
		return l.state.LookaheadCache.GetItem(count - 1)
	}

	state := l.GetState()
	var token Token

	for i := 0; i < count; i++ {
		token = l.nextToken()
		state.LookaheadCache.AddItem(token)

		if token.Is(TypeEof) {
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
	// check if we have anything in the lookahead cache
	if l.state.LookaheadCache.ContainsItems() {
		return l.state.LookaheadCache.PluckItem()
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

func (l *Lexer) GetPosition() Position {
	if l.OmitTokenPosition {
		return Position{}
	}

	if l.state.Cursor == l.state.CachedPositionCursor {
		return l.state.CachedPosition
	}

	if l.state.LineIndexesCount == 0 {
		l.state.CachedPosition.Col += l.state.Cursor - l.state.CachedPositionCursor
		l.state.CachedPositionCursor = l.state.Cursor
		return l.state.CachedPosition
	}

	nextLineStart := l.state.LineIndexes[0]
	if l.state.Cursor <= nextLineStart {
		l.state.CachedPosition.Col += l.state.Cursor - l.state.CachedPositionCursor
		l.state.CachedPositionCursor = l.state.Cursor
		return l.state.CachedPosition
	}

	l.state.CachedPosition.Row += 1
	l.state.CachedPosition.Col = l.state.Cursor - nextLineStart
	l.state.CachedPositionCursor = l.state.Cursor

	// Remove the consumed entry
	l.state.LineIndexes = l.state.LineIndexes[1:]
	l.state.LineIndexesCount -= 1

	return l.state.CachedPosition
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

// NextCharsAre checks if the next chars from the cursor on match the provided chars
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
// Parsing Utilities
// ---------------------------------------------------------------

func (l *Lexer) CollectTokensBetweenParentheses() (Tokens, int, int, error) {
	return l.CollectTokensBetween(TypeOpenParen, TypeCloseParen)
}

func (l *Lexer) CollectTokensBetweenCurlyBraces() (Tokens, int, int, error) {
	return l.CollectTokensBetween(TypeOpenCurly, TypeCloseCurly)
}

func (l *Lexer) CollectTokensBetween(open TokenType, close TokenType) (Tokens, int, int, error) {
	tokens := Tokens{}
	token := *l.state.CurrentToken

	if !token.Is(open) {
		return tokens, -1, -1, fmt.Errorf("Current token is not of opener type %s", open)
	}

	start := l.GetCursor()
	end := start
	level := 1

	for !token.Is(TypeEof) {
		end = l.GetCursor()
		token = l.NextToken()

		if token.Is(TypeEof) {
			return tokens, start, end, fmt.Errorf("Unexpected EndOfFile")
		}

		if token.Is(close) {
			level -= 1
			if level == 0 {
				break
			}
		}

		if token.Is(open) {
			level += 1
		}

		tokens = append(tokens, token)
	}

	return tokens, start, end, nil
}

func (l *Lexer) GetTokensDelimited(tokenType TokenType, delimiter TokenType) (Tokens, error) {
	tokens := Tokens{}

	token := *l.state.CurrentToken
	for !token.Is(TypeEof) {
		if !token.Is(tokenType) {
			return tokens, fmt.Errorf("expected %s but found %s", tokenType, token.Type)
		}

		tokens = append(tokens, token)

		if !l.Lookahead(1).Is(delimiter) {
			break
		}

		token = l.NextToken() // Just consume the delimiter
		token = l.NextToken()
	}

	return tokens, nil
}

func (l *Lexer) GetAnyTokenDelimited(delimiter TokenType) ([]Token, error) {
	tokens := []Token{}

	token := *l.state.CurrentToken
	for !token.Is(TypeEof) {
		tokens = append(tokens, token)

		if !l.Lookahead(1).Is(delimiter) {
			break
		}

		token = l.NextToken() // Just consume the delimiter
		token = l.NextToken()
	}

	return tokens, nil
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
