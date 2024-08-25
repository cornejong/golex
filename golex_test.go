package golex

import (
	"fmt"
	"testing"
)

var source string = " func() { test = \"SomeStringValue\"; test = 1.2; test = 88 }"
var expected []Token = []Token{
	{Type: TypeKeyword, Literal: "func", Position: Position{Row: 1, Col: 2}},
	{Type: TypeOpenParen, Literal: "(", Position: Position{Row: 1, Col: 6}},
	{Type: TypeCloseParen, Literal: ")", Position: Position{Row: 1, Col: 7}},
	{Type: TypeOpenCurly, Literal: "{", Position: Position{Row: 1, Col: 9}},
	{Type: TypeSymbol, Literal: "test", Position: Position{Row: 1, Col: 11}},
	{Type: TypeAssign, Literal: "=", Position: Position{Row: 1, Col: 16}},
	{Type: TypeDoubleQuoteString, Literal: "\"SomeStringValue\"", Value: "SomeStringValue", Position: Position{Row: 1, Col: 18}},
	{Type: TypeSemicolon, Literal: ";", Position: Position{Row: 1, Col: 35}},
	{Type: TypeSymbol, Literal: "test", Position: Position{Row: 1, Col: 37}},
	{Type: TypeAssign, Literal: "=", Position: Position{Row: 1, Col: 42}},
	{Type: TypeFloat, Literal: "1.2", Value: 1.2, Position: Position{Row: 1, Col: 44}},
	{Type: TypeSemicolon, Literal: ";", Position: Position{Row: 1, Col: 47}},
	{Type: TypeSymbol, Literal: "test", Position: Position{Row: 1, Col: 49}},
	{Type: TypeAssign, Literal: "=", Position: Position{Row: 1, Col: 54}},
	{Type: TypeInteger, Literal: "88", Value: 88, Position: Position{Row: 1, Col: 56}},
	{Type: TypeCloseCurly, Literal: "}", Position: Position{Row: 1, Col: 59}},
	{Type: TypeEof, Literal: string(EOF), Position: Position{Row: 1, Col: 60}},
}

func getLexer() *Lexer {
	return NewLexer(
		// DebugPrintTokens(),
		WithKeywords("fun", "func", "def"),
	)
}

func TestGolexUsageToSlice(t *testing.T) {
	fmt.Println("TestGolexUsageToSlice...")

	lexer := getLexer()
	tokens := lexer.TokenizeToSlice(source)
	differ := &Differ{}
	differ.Compare(expected, tokens)
	if differ.HasDifference() {
		fmt.Println(differ)
		fmt.Printf("\n%d differences between expected and result\n", len(differ.Diffs))
		t.FailNow()
	}
}

func TestGolexUsageToChannel(t *testing.T) {
	fmt.Println("TestGolexUsageToChannel...")

	lexer := getLexer()
	tokensChannel := make(chan Token)
	tokens := []Token{}
	go lexer.TokenizeToChannel(source, tokensChannel)

	for token := range tokensChannel {
		tokens = append(tokens, token)
	}

	differ := &Differ{}
	differ.Compare(expected, tokens)
	if differ.HasDifference() {
		fmt.Println(differ)
		fmt.Printf("\n%d differences between expected and result\n", len(differ.Diffs))
		t.FailNow()
	}
}

func TestGolexUsageToCallback(t *testing.T) {
	fmt.Println("TestGolexUsageToCallback...")

	lexer := getLexer()
	tokens := []Token{}
	lexer.TokenizeToCallback(source, func(t Token) {
		tokens = append(tokens, t)
	})

	differ := &Differ{}
	differ.Compare(expected, tokens)
	if differ.HasDifference() {
		fmt.Println(differ)
		fmt.Printf("\n%d differences between expected and result\n", len(differ.Diffs))
		t.FailNow()
	}
}

func TestGolexUsageManual(t *testing.T) {
	fmt.Println("TestGolexUsageManual...")

	lexer := getLexer()
	lexer.TokenizeManual(source)
	tokens := []Token{}

	for !lexer.ReachedEOF() {
		tokens = append(tokens, lexer.NextToken())
	}

	differ := &Differ{}
	differ.Compare(expected, tokens)
	if differ.HasDifference() {
		fmt.Println(differ)
		fmt.Printf("\n%d differences between expected and result\n", len(differ.Diffs))
		t.FailNow()
	}
}

func TestMultiLinePosition(t *testing.T) {
	fmt.Println("TestMultiLinePosition...")

	lexer := getLexer()
	lexer.TokenizeManual("a = true;\nb = false")
	tokens := []Token{}

	for !lexer.ReachedEOF() {
		tokens = append(tokens, lexer.NextToken())
	}

	expect := []Token{
		{Type: TypeSymbol, Literal: "a", Position: Position{Col: 1, Row: 1}},
		{Type: TypeAssign, Literal: "=", Position: Position{Col: 3, Row: 1}},
		{Type: TypeBool, Literal: "true", Value: true, Position: Position{Col: 5, Row: 1}},
		{Type: TypeSemicolon, Literal: ";", Position: Position{Col: 9, Row: 1}},
		{Type: TypeSymbol, Literal: "b", Position: Position{Col: 1, Row: 2}},
		{Type: TypeAssign, Literal: "=", Position: Position{Col: 3, Row: 2}},
		{Type: TypeBool, Literal: "false", Value: false, Position: Position{Col: 5, Row: 2}},
		{Type: TypeEof, Literal: string(EOF), Position: Position{Col: 10, Row: 2}},
	}

	differ := &Differ{}
	differ.Compare(expect, tokens)
	if differ.HasDifference() {
		fmt.Println(differ)
		fmt.Printf("\n%d differences between expected and result\n", len(differ.Diffs))
		t.FailNow()
	}
}

func TestRetainWhitespace(t *testing.T) {
	fmt.Println("TestRetainWhitespace...")

	lexer := getLexer()
	RetainWhitespace()(lexer)
	lexer.TokenizeManual("a = true;\nb = false")
	tokens := []Token{}

	for !lexer.ReachedEOF() {
		tokens = append(tokens, lexer.NextToken())
	}

	expect := []Token{
		{Type: TypeSymbol, Literal: "a", Position: Position{Col: 1, Row: 1}},
		{Type: TypeSpace, Literal: " ", Position: Position{Col: 2, Row: 1}},
		{Type: TypeAssign, Literal: "=", Position: Position{Col: 3, Row: 1}},
		{Type: TypeSpace, Literal: " ", Position: Position{Col: 4, Row: 1}},
		{Type: TypeBool, Literal: "true", Value: true, Position: Position{Col: 5, Row: 1}},
		{Type: TypeSemicolon, Literal: ";", Position: Position{Col: 9, Row: 1}},
		{Type: TypeNewline, Literal: "\n", Position: Position{Col: 10, Row: 1}},
		{Type: TypeSymbol, Literal: "b", Position: Position{Col: 1, Row: 2}},
		{Type: TypeSpace, Literal: " ", Position: Position{Col: 2, Row: 2}},
		{Type: TypeAssign, Literal: "=", Position: Position{Col: 3, Row: 2}},
		{Type: TypeSpace, Literal: " ", Position: Position{Col: 4, Row: 2}},
		{Type: TypeBool, Literal: "false", Value: false, Position: Position{Col: 5, Row: 2}},
		{Type: TypeEof, Literal: string(EOF), Position: Position{Col: 10, Row: 2}},
	}

	differ := &Differ{}
	differ.Compare(expect, tokens)
	if differ.HasDifference() {
		fmt.Println(differ)
		fmt.Printf("\n%d differences between expected and result\n", len(differ.Diffs))
		t.FailNow()
	}
}

func TestTripleBacktickString(t *testing.T) {
	fmt.Println("TestTripleBacktickString...")

	lexer := getLexer()
	WithStringEnclosure(TripleBacktickStringEnclosure)(lexer)

	lexer.TokenizeManual("```a string```")
	tokens := []Token{}

	for !lexer.ReachedEOF() {
		tokens = append(tokens, lexer.NextToken())
	}

	expect := []Token{
		{Type: TypeTripleBacktickString, Literal: "```a string```", Value: "a string", Position: Position{Col: 1, Row: 1}},
		{Type: TypeEof, Literal: string(EOF), Position: Position{Row: 1, Col: 15}},
	}

	differ := &Differ{}
	differ.Compare(expect, tokens)
	if differ.HasDifference() {
		fmt.Println(differ)
		fmt.Printf("\n%d differences between expected and result\n", len(differ.Diffs))
		t.FailNow()
	}
}

func TestBooleanTokenizer(t *testing.T) {
	fmt.Println("TestBooleanTokenizer...")

	lexer := getLexer()
	lexer.TokenizeManual("a = true; b = false")
	tokens := []Token{}

	for !lexer.ReachedEOF() {
		tokens = append(tokens, lexer.NextToken())
	}

	expect := []Token{
		{Type: TypeSymbol, Literal: "a", Position: Position{Col: 1, Row: 1}},
		{Type: TypeAssign, Literal: "=", Position: Position{Col: 3, Row: 1}},
		{Type: TypeBool, Literal: "true", Value: true, Position: Position{Col: 5, Row: 1}},
		{Type: TypeSemicolon, Literal: ";", Position: Position{Col: 9, Row: 1}},
		{Type: TypeSymbol, Literal: "b", Position: Position{Col: 11, Row: 1}},
		{Type: TypeAssign, Literal: "=", Position: Position{Col: 13, Row: 1}},
		{Type: TypeBool, Literal: "false", Value: false, Position: Position{Col: 15, Row: 1}},
		{Type: TypeEof, Literal: string(EOF), Position: Position{Col: 20, Row: 1}},
	}

	differ := &Differ{}
	differ.Compare(expect, tokens)
	if differ.HasDifference() {
		fmt.Println(differ)
		fmt.Printf("\n%d differences between expected and result\n", len(differ.Diffs))
		t.FailNow()
	}
}

func TestLexerLookaheadCache(t *testing.T) {
	fmt.Println("TestLexerLookaheadCache...")

	lexer := getLexer()
	lexer.TokenizeManual(source)

	lexer.Lookahead(2)

	if lexer.state.LookaheadCache.count != 2 {
		t.Errorf("Expected the lookahead cache to contain 2 items. %d items were found.", lexer.state.LookaheadCache.count)
	}

	lexer.NextToken()

	if lexer.state.LookaheadCache.count != 1 {
		t.Errorf("Expected the lookahead cache to contain 1 item. %d items were found.", lexer.state.LookaheadCache.count)
	}

	lexer.Lookahead(2)

	if lexer.state.LookaheadCache.count != 2 {
		t.Errorf("Expected the lookahead cache to contain 2 items. %d items were found.", lexer.state.LookaheadCache.count)
	}

	lexer.NextToken()

	if lexer.state.LookaheadCache.count != 1 {
		t.Errorf("Expected the lookahead cache to contain 1 item. %d items were found.", lexer.state.LookaheadCache.count)
	}

	lexer.NextToken()

	if lexer.state.LookaheadCache.count != 0 {
		t.Errorf("Expected the lookahead cache to contain 0 items. %d items were found.", lexer.state.LookaheadCache.count)
	}
}

// ###################################################
// #              Test Individual Stuff
// ###################################################
func TestCharacterPatternExpand(t *testing.T) {
	fmt.Println("TestCharacterPatternExpand...")

	expect := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	expanded, err := expandCharacterPattern("a-zA-Z0-9_")
	if err != nil {
		panic(err)
	}

	if expanded != expect {
		fmt.Printf("Expected '%s' but got '%s'\n", expect, expanded)
		t.FailNow()
	}
}

func TestTokenizationInsertOrderBefore(t *testing.T) {
	fmt.Println("TestTokenizationInsertOrderBefore...")

	tokenizers := map[TokenizerType]Tokenizer{}
	order := []TokenizerType{
		TypeCommentTokenizer,
		TypeNumberTokenizer,
		TypeLiteralTokenizer,
		TypeStringTokenizer,
	}

	insert := InsertBefore(TypeNumberTokenizer, TypeSymbolTokenizer, SymbolTokenizer{})
	_, newOrder := insert.Insert(tokenizers, order)

	expectOrder := []TokenizerType{
		TypeCommentTokenizer,
		TypeSymbolTokenizer,
		TypeNumberTokenizer,
		TypeLiteralTokenizer,
		TypeStringTokenizer,
	}

	differ := &Differ{}
	differ.Compare(expectOrder, newOrder)
	if differ.HasDifference() {
		fmt.Println(differ)
		fmt.Printf("\n%d differences between expected and result\n", len(differ.Diffs))
		t.FailNow()
	}

}

func TestTokenizationInsertOrderAfter(t *testing.T) {
	fmt.Println("TestTokenizationInsertOrderAfter...")

	tokenizers := map[TokenizerType]Tokenizer{}
	order := []TokenizerType{
		TypeCommentTokenizer,
		TypeNumberTokenizer,
		TypeLiteralTokenizer,
		TypeStringTokenizer,
	}

	insert := InsertAfter(TypeNumberTokenizer, TypeSymbolTokenizer, SymbolTokenizer{})
	_, newOrder := insert.Insert(tokenizers, order)

	expectOrder := []TokenizerType{
		TypeCommentTokenizer,
		TypeNumberTokenizer,
		TypeSymbolTokenizer,
		TypeLiteralTokenizer,
		TypeStringTokenizer,
	}

	differ := &Differ{}
	differ.Compare(expectOrder, newOrder)
	if differ.HasDifference() {
		fmt.Println(differ)
		fmt.Printf("\n%d differences between expected and result\n", len(differ.Diffs))
		t.FailNow()
	}

}
