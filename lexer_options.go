package golex

import "slices"

type LexerOptionFunc func(*Lexer)

func DebugPrintTokens() LexerOptionFunc {
	return LexerOptionFunc(func(l *Lexer) {
		l.DebugPrintTokens = true
	})
}

func OmitTokenPosition() LexerOptionFunc {
	return LexerOptionFunc(func(l *Lexer) {
		l.OmitTokenPosition = true
	})
}

func IgnoreTokens(types ...TokenType) LexerOptionFunc {
	return LexerOptionFunc(func(l *Lexer) {
		l.IgnoreTokens = append(l.IgnoreTokens, types...)
	})
}

func RetainWhitespace() LexerOptionFunc {
	return LexerOptionFunc(func(l *Lexer) {
		l.IgnoreWhitespace = false
	})
}

func WithKeywords(keywords ...string) LexerOptionFunc {
	return LexerOptionFunc(func(l *Lexer) {
		l.Keywords = append(l.Keywords, keywords...)
		if len(l.Keywords) > 0 {
			l.CheckForKeywords = true
		}
	})
}

func SymbolCharacterMap(startCharMap, continueCharMap string) LexerOptionFunc {
	expandedStartMap, err := expandCharacterPattern(startCharMap)
	if err != nil {
		// TODO: Maybe not panic this...
		panic(err)
	}

	expandedContinueMap, err := expandCharacterPattern(continueCharMap)
	if err != nil {
		// TODO: Maybe not panic this...
		panic(err)
	}

	return LexerOptionFunc(func(l *Lexer) {
		l.SymbolStartCharacterMap = expandedStartMap
		l.SymbolContinueCharacterMap = expandedContinueMap
	})
}

func WithTokenizer(inserter TokenizerInserter) LexerOptionFunc {
	return LexerOptionFunc(func(l *Lexer) {
		l.tokenizers, l.tokenizationOrder = inserter.Insert(l.tokenizers, l.tokenizationOrder)
	})
}

func WithLiteralTokens(literalTokens ...LiteralToken) LexerOptionFunc {
	return LexerOptionFunc(func(l *Lexer) {
		l.LiteralTokens = SortLiteralTokens(append(l.LiteralTokens, literalTokens...))
	})
}

func WithoutLiteralTokens(literalTokens ...TokenType) LexerOptionFunc {
	return LexerOptionFunc(func(l *Lexer) {
		literals := []LiteralToken{}

		for _, t := range l.LiteralTokens {
			if !slices.Contains(literalTokens, t.Type) {
				literals = append(literals, t)
			}
		}

		l.LiteralTokens = literals
	})
}

func WithCommentSyntax(syntaxes ...CommentSyntax) LexerOptionFunc {
	return LexerOptionFunc(func(l *Lexer) {
		l.CommentSyntaxes = append(l.CommentSyntaxes, syntaxes...)
	})
}

func WithoutCommentSyntax(syntaxes ...CommentSyntax) LexerOptionFunc {
	return LexerOptionFunc(func(l *Lexer) {
		commentSyntax := []CommentSyntax{}

		for _, s := range l.CommentSyntaxes {
			if !slices.Contains(syntaxes, s) {
				commentSyntax = append(commentSyntax, s)
			}
		}

		l.CommentSyntaxes = commentSyntax
	})
}

func WithStringEnclosure(enclosures ...StringEnclosure) LexerOptionFunc {
	return LexerOptionFunc(func(l *Lexer) {
		l.StringEnclosures = append(l.StringEnclosures, enclosures...)
	})
}

func WithoutStringEnclosure(enclosures ...string) LexerOptionFunc {
	return LexerOptionFunc(func(l *Lexer) {
		stringEnclosures := []StringEnclosure{}

		for _, e := range l.StringEnclosures {
			if !slices.Contains(enclosures, e.Enclosure) {
				stringEnclosures = append(stringEnclosures, e)
			}
		}

		l.StringEnclosures = stringEnclosures
	})
}
