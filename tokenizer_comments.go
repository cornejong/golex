package golex

import "fmt"

var (
	SlashSingleLineCommentSyntax   = CommentSyntax{Opener: "//"}
	SlashMultilineCommentSyntax    = CommentSyntax{Opener: "/*", Closer: "*/"}
	HashtagSingleLineCommentSyntax = CommentSyntax{Opener: "#"}

	cachedCommentSyntax *CommentSyntax
)

type CommentSyntax struct {
	Opener string
	Closer string
}

type CommentTokenizer struct{}

func (c CommentTokenizer) CanTokenize(l *Lexer) bool {
	if len(l.CommentSyntaxes) < 1 {
		return false
	}

	for _, syntax := range l.CommentSyntaxes {
		if l.NextCharsAre([]rune(syntax.Opener)) {
			cachedCommentSyntax = &syntax
			return true
		}
	}

	return false
}

func (c CommentTokenizer) Tokenize(l *Lexer) (Token, error) {
	if cachedCommentSyntax == nil {
		if !c.CanTokenize(l) {
			return Token{Type: TypeInvalid, Position: l.GetPosition()},
				NewError(fmt.Sprintf("Invalid token '%c' found", l.CharAtCursor()), l.GetPosition(), l.GetCursor(), l.state.Content)
		} else {
			return c.Tokenize(l)
		}
	}

	var reachedEndOfComment func(*Lexer) bool
	if cachedCommentSyntax.Closer == "" {
		reachedEndOfComment = func(l *Lexer) bool {
			return l.CharAtCursor() == '\n'
		}
	} else {
		reachedEndOfComment = func(l *Lexer) bool {
			return l.NextCharsAre([]rune(cachedCommentSyntax.Closer))
		}
	}

	token := Token{Type: TypeComment, Position: l.GetPosition()}
	for l.GetCursor() <= l.state.ContentLength && !reachedEndOfComment(l) {
		token.AppendChar(l.CharAtCursor())
		l.IncrementCursor(1)
	}

	l.IncrementCursor(len(cachedCommentSyntax.Closer))

	cachedCommentSyntax = nil

	if l.IgnoreComments {
		return l.NextToken()
	}

	return token, nil
}
