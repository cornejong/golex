package golex

type BuildInType string

func (bit BuildInType) String() string {
	return string(bit)
}

const (
	TypeSof     BuildInType = "StartOfFile"
	TypeEof     BuildInType = "EndOfFile"
	TypeInvalid BuildInType = "Invalid"

	TypeString               BuildInType = "String"
	TypeDoubleQuoteString    BuildInType = "DoubleQuoteString"
	TypeSingleQuoteString    BuildInType = "SingleQuoteString"
	TypeBacktickString       BuildInType = "BacktickString"
	TypeTripleBacktickString BuildInType = "TripleBacktickString"
	TypeNumber               BuildInType = "Number"
	TypeInteger              BuildInType = "Integer"
	TypeFloat                BuildInType = "Float"
	TypeBool                 BuildInType = "Boolean"
	TypeNull                 BuildInType = "Null"
	TypeNil                  BuildInType = "Nil"

	TypeComment    BuildInType = "Comment"
	TypeKeyword    BuildInType = "Keyword"
	TypeIdentifier BuildInType = "Identifier"
	TypeSymbol     BuildInType = "Symbol"

	TypePlus               BuildInType = "Plus"               // +
	TypeMinus              BuildInType = "Minus"              // -
	TypeMultiply           BuildInType = "Multiply"           // *
	TypeDivide             BuildInType = "Divide"             // /
	TypeModulo             BuildInType = "Modulo"             // %
	TypeAssign             BuildInType = "Assign"             // =
	TypeEqual              BuildInType = "Equal"              // ==
	TypeNotEqual           BuildInType = "NotEqual"           // !=
	TypeLessThan           BuildInType = "LessThan"           // <
	TypeGreaterThan        BuildInType = "GreaterThan"        // >
	TypeLessThanOrEqual    BuildInType = "LessThanOrEqual"    // <=
	TypeGreaterThanOrEqual BuildInType = "GreaterThanOrEqual" // >=
	TypeAnd                BuildInType = "And"                // &&
	TypeOr                 BuildInType = "Or"                 // ||
	TypeNot                BuildInType = "Not"                // !

	TypeOpenParen   BuildInType = "OpenParenthesis"    // (
	TypeCloseParen  BuildInType = "CloseParenthesis"   // )
	TypeOpenCurly   BuildInType = "OpenCurlyBracket"   // {
	TypeCloseCurly  BuildInType = "CloseCurlyBracket"  // }
	TypeOpenSquare  BuildInType = "OpenSquareBracket"  // [
	TypeCloseSquare BuildInType = "CloseSquareBracket" // ]
	TypeComma       BuildInType = "Comma"              // ,
	TypeDot         BuildInType = "Dot"                // .
	TypeColon       BuildInType = "Colon"              // :
	TypeSemicolon   BuildInType = "Semicolon"          // ;

	TypeArrowRight   BuildInType = "ArrowRight"   // ->
	TypeArrowLeft    BuildInType = "ArrowLeft"    // <-
	TypeQuestionMark BuildInType = "QuestionMark" // ?
	TypeTilde        BuildInType = "Tilde"        // ~
	TypeAmpersand    BuildInType = "Ampersand"    // &
	TypePipe         BuildInType = "Pipe"         // |
	TypeCaret        BuildInType = "Caret"        // ^
	TypeDollar       BuildInType = "Dollar"       // $
	TypeHash         BuildInType = "Hash"         // #
	TypeAt           BuildInType = "At"           // @
	TypeEllipses     BuildInType = "Ellipses"     //...

	TypeSpace          BuildInType = "Space"
	TypeTab            BuildInType = "Tab"
	TypeNewline        BuildInType = "Newline"
	TypeCarriageReturn BuildInType = "CarriageReturn"
	TypeFormFeed       BuildInType = "FormFeed"
)

var buildInLiteralTokens []LiteralToken = []LiteralToken{
	LiteralToken{TypeEllipses, "..."},
	LiteralToken{TypeOpenCurly, "{"},
	LiteralToken{TypeCloseCurly, "}"},
	LiteralToken{TypeOpenParen, "("},
	LiteralToken{TypeCloseParen, ")"},
	LiteralToken{TypeOpenSquare, "["},
	LiteralToken{TypeCloseSquare, "]"},
	LiteralToken{TypeComma, ","},
	LiteralToken{TypeDot, "."},
	LiteralToken{TypeColon, ":"},
	LiteralToken{TypeSemicolon, ";"},
	LiteralToken{TypePlus, "+"},
	LiteralToken{TypeMinus, "-"},
	LiteralToken{TypeMultiply, "*"},
	LiteralToken{TypeDivide, "/"},
	LiteralToken{TypeModulo, "%"},
	LiteralToken{TypeAssign, "="},
	LiteralToken{TypeEqual, "=="},
	LiteralToken{TypeNotEqual, "!="},
	LiteralToken{TypeLessThan, "<"},
	LiteralToken{TypeGreaterThan, ">"},
	LiteralToken{TypeLessThanOrEqual, "<="},
	LiteralToken{TypeGreaterThanOrEqual, ">="},
	LiteralToken{TypeAnd, "&&"},
	LiteralToken{TypeOr, "||"},
	LiteralToken{TypeNot, "!"},
	LiteralToken{TypeArrowRight, "->"},
	LiteralToken{TypeArrowLeft, "<-"},
	LiteralToken{TypeQuestionMark, "?"},
	LiteralToken{TypeTilde, "~"},
	LiteralToken{TypeAmpersand, "&"},
	LiteralToken{TypePipe, "|"},
	LiteralToken{TypeCaret, "^"},
	LiteralToken{TypeDollar, "$"},
	LiteralToken{TypeHash, "#"},
	LiteralToken{TypeAt, "@"},

	LiteralToken{TypeSpace, " "},
	LiteralToken{TypeTab, "\t"},
	LiteralToken{TypeNewline, "\n"},
	LiteralToken{TypeCarriageReturn, "\r"},
	LiteralToken{TypeFormFeed, "\f"},
}
