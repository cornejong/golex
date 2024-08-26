package golex

import (
	"fmt"
	"iter"
)

// ---------------------------------------------------------------
// Lookahead Helpers
// ---------------------------------------------------------------

// NextTokenIs checks if the next token is the same as
// the provided token without consuming the token
func (l *Lexer) NextTokenIs(token Token) bool {
	return l.Lookahead(1).Is(token)
}

// NextTokenIsAnyOf checks if the next token is of any of
// the provided tokens without consuming the token
func (l *Lexer) NextTokenIsAnyOf(tokens ...Token) bool {
	return l.Lookahead(1).IsAnyOf(tokens...)
}

// NextTokenSequenceIs checks if the next sequence of tokens
// in the lexer matches the provided token sequence
// without consuming the tokens
func (l *Lexer) NextTokenSequenceIs(tokens ...Token) bool {
	i := 0
	for token := range l.LookaheadIterator(len(tokens)) {
		if !token.Is(tokens[i]) {
			return false
		}

		i += 1
	}

	return true
}

// ---------------------------------------------------------------
// Iterators
// ---------------------------------------------------------------

// IterateTokensBetweenParentheses returns an iterator that iterates over all
// the tokens between the opening and closing parentheses starting at the offset start
// It assumes that the offset start contains the opening parenthesis.
// Nested parenthesis will be contained in the output tokens until
// the matching closing parenthesis to the start is found.
// In addition it returns the start and end cursor position for the iterated portion
func (l *Lexer) IterateTokensBetweenParentheses() (iter.Seq2[Token, error], *int, *int, error) {
	return l.IterateTokensBetween(TypeOpenParen, TypeCloseParen)
}

// IterateTokensBetweenCurlyBraced returns an iterator that iterates over all
// the tokens between the opening and closing curly braces starting at the offset start
// It assumes that the offset start contains the opening curly brace.
// Nested curly braces will be contained in the output tokens until
// the matching closing curly brace to the start is found.
// In addition it returns the start and end cursor position for the iterated portion
func (l *Lexer) IterateTokensBetweenCurlyBraces() (iter.Seq2[Token, error], *int, *int, error) {
	return l.IterateTokensBetween(TypeOpenCurly, TypeCloseCurly)
}

// IterateTokensBetween returns an iterator that iterates over all
// the tokens between the open and close type starting at the offset start
// it assumes that the offset start contains the opening token.
// Nested openers and closers will be contained in the output
// tokens until the matching closer is found.
// In addition it returns the start and end cursor position for the iterated portion
func (l *Lexer) IterateTokensBetween(open TokenType, close TokenType) (iter.Seq2[Token, error], *int, *int, error) {
	var err error
	token := l.CurrentToken()
	start := l.GetCursor() + len(token.Literal)
	end := l.GetCursor()
	level := 1

	if !token.TypeIs(open) {
		return nil, &start, &end, fmt.Errorf("Current token is not of opener type %s", open)
	}

	return func(yield func(Token, error) bool) {
		for !token.TypeIs(TypeEof) {
			end = l.GetCursor()
			token, err = l.NextToken()
			if err != nil {
				yield(token, err)
			}

			if token.TypeIs(TypeEof) {
				if !yield(token, nil) {
					return
				}

				return
			}

			if token.TypeIs(close) {
				level -= 1
				if level == 0 {
					end = l.GetCursor() - len(token.Literal)
					return
				}
			}

			if token.TypeIs(open) {
				level += 1
			}

			if !yield(token, nil) {
				return
			}
		}

		return
	}, &start, &end, nil
}

func (l *Lexer) IterateTokensDelimited(tokenType TokenType, delimiter TokenType) iter.Seq2[Token, error] {
	var err error
	token := l.CurrentToken()
	return func(yield func(Token, error) bool) {
		for !token.TypeIs(TypeEof) {
			if !token.TypeIs(tokenType) {
				yield(token, fmt.Errorf("expected %s but found %s", tokenType, token.Type))
			}

			if !yield(token, nil) {
				return
			}

			if !l.Lookahead(1).TypeIs(delimiter) {
				return
			}

			token, err = l.NextToken() // Just consume the delimiter
			if err != nil {
				if !yield(token, err) {
					return
				}
			}

			token, err = l.NextToken()
			if err != nil {
				if !yield(token, err) {
					return
				}
			}
		}

		return
	}
}

func (l *Lexer) IterateAnyTokenDelimited(delimiter TokenType) iter.Seq2[Token, error] {
	var err error
	token := l.CurrentToken()
	return func(yield func(Token, error) bool) {
		for !token.TypeIs(TypeEof) {
			if !yield(token, nil) {
				return
			}

			if !l.Lookahead(1).TypeIs(delimiter) {
				return
			}

			token, err = l.NextToken() // Just consume the delimiter
			if err != nil {
				if !yield(token, err) {
					return
				}
			}

			token, err = l.NextToken()
			if err != nil {
				if !yield(token, err) {
					return
				}
			}
		}

		return
	}
}

// ---------------------------------------------------------------
// Collectors
// ---------------------------------------------------------------

// CollectTokensBetweenParentheses collects all the tokens between
// the opening and closing parentheses starting at the offset start
// It assumes that the offset start contains the opening parenthesis.
// Nested parenthesis will be contained in the output tokens until
// the matching closing parenthesis to the start is found.
// In addition it returns the start and end cursor position for the collected portion
func (l *Lexer) CollectTokensBetweenParentheses() (Tokens, int, int, error) {
	return l.CollectTokensBetween(TypeOpenParen, TypeCloseParen)
}

// CollectTokensBetweenCurlyBraced collects all the tokens between
// the opening and closing curly braces starting at the offset start
// It assumes that the offset start contains the opening curly brace.
// Nested curly braces will be contained in the output tokens until
// the matching closing curly brace to the start is found.
// In addition it returns the start and end cursor position for the collected portion
func (l *Lexer) CollectTokensBetweenCurlyBraces() (Tokens, int, int, error) {
	return l.CollectTokensBetween(TypeOpenCurly, TypeCloseCurly)
}

// CollectTokensBetween collects all the tokens between the open
// and close type starting at the offset start
// it assumes that the offset start contains the opening token.
// Nested openers and closers will be contained in the output
// tokens until the matching closer is found.
// In addition it returns the start and end cursor position for the collected portion
func (l *Lexer) CollectTokensBetween(open TokenType, close TokenType) (Tokens, int, int, error) {
	var err error

	tokens := Tokens{}
	token := l.CurrentToken()

	if !token.TypeIs(open) {
		return tokens, -1, -1, fmt.Errorf("Current token is not of opener type %s", open)
	}

	start := l.GetCursor()
	end := start
	level := 1

	for !token.TypeIs(TypeEof) {
		end = l.GetCursor()
		token, err = l.NextToken()
		if err != nil {
			return tokens, start, end, err
		}

		if token.TypeIs(TypeEof) {
			return tokens, start, end, fmt.Errorf("Unexpected EndOfFile")
		}

		if token.TypeIs(close) {
			level -= 1
			if level == 0 {
				break
			}
		}

		if token.TypeIs(open) {
			level += 1
		}

		tokens = append(tokens, token)
	}

	return tokens, start, end, nil
}

func (l *Lexer) CollectTokensDelimited(tokenType TokenType, delimiter TokenType) (Tokens, error) {
	var err error

	tokens := Tokens{}

	token := l.CurrentToken()
	for !token.TypeIs(TypeEof) {
		if !token.TypeIs(tokenType) {
			return tokens, fmt.Errorf("expected %s but found %s", tokenType, token.Type)
		}

		tokens = append(tokens, token)

		if !l.Lookahead(1).TypeIs(delimiter) {
			break
		}

		token, err = l.NextToken() // Just consume the delimiter
		if err != nil {
			return tokens, err
		}

		token, err = l.NextToken()
		if err != nil {
			return tokens, err
		}
	}

	return tokens, nil
}

func (l *Lexer) CollectAnyTokenDelimited(delimiter TokenType) (Tokens, error) {
	var err error

	tokens := Tokens{}

	token := l.CurrentToken()
	for !token.TypeIs(TypeEof) {
		tokens = append(tokens, token)

		if !l.Lookahead(1).TypeIs(delimiter) {
			break
		}

		token, err = l.NextToken() // Just consume the delimiter
		if err != nil {
			return tokens, err
		}

		token, err = l.NextToken()
		if err != nil {
			return tokens, err
		}
	}

	return tokens, nil
}
