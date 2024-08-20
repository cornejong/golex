package golex

var (
	SlashSingleLineCommentSyntax   = CommentSyntax{Opener: "//"}
	SlashMultilineCommentSyntax    = CommentSyntax{Opener: "/*", Closer: "*/"}
	HashtagSingleLineCommentSyntax = CommentSyntax{Opener: "#"}
)

type CommentTokenizerCacheKey string

type CommentSyntax struct {
	Opener string
	Closer string
}

type CommentTokenizer struct {
	syntaxes      []CommentSyntax
	currentSyntax *CommentSyntax
}

func NewCommentTokenizer() CommentTokenizer {
	return CommentTokenizer{}
}

func (c CommentTokenizer) CanTokenize(l *Lexer) bool {
	if len(c.syntaxes) < 1 {
		return false
	}

	for _, syntax := range l.CommentSyntaxes {
		if l.NextCharsAre([]rune(syntax.Opener)) {
			l.state.Cache[CommentTokenizerCacheKey("currentSyntax")] = syntax
			return true
		}
	}

	return false
}

func (c CommentTokenizer) Tokenize(l *Lexer) Token {
	var syntax CommentSyntax
	if syntaxI, ok := l.state.Cache[CommentTokenizerCacheKey("currentSyntax")]; ok {
		syntax = syntaxI.(CommentSyntax)
		delete(l.state.Cache, CommentTokenizerCacheKey("currentSyntax"))
	} else {
		if !c.CanTokenize(l) {
			return Token{Type: TypeInvalid, Position: l.GetPosition()}
		} else {
			return c.Tokenize(l)
		}
	}

	var reachedEndOfComment func(*Lexer) bool
	if syntax.Closer == "" {
		reachedEndOfComment = func(l *Lexer) bool {
			return l.CharAtCursor() == '\n'
		}
	} else {
		reachedEndOfComment = func(l *Lexer) bool {
			return l.NextCharsAre([]rune(syntax.Closer))
		}
	}

	token := Token{Type: TypeComment, Position: l.GetPosition()}
	for l.GetCursor() <= l.state.ContentLength && !reachedEndOfComment(l) {
		token.AppendChar(l.CharAtCursor())
		l.IncrementCursor(1)
	}

	l.IncrementCursor(len(c.currentSyntax.Closer))

	if l.IgnoreComments {
		return l.NextToken()
	}

	return token
}
