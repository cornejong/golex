# Golex
A lexing and parsing toolkit for go

## Features
- **Multiple Tokenizers**: Supports built-in tokenizers for comments, literals, numbers, booleans, strings, and symbols, with the ability to add custom tokenizers.
- **Flexible Lexer Options**: Configure the lexer with options like retaining whitespace or customizing keyword sets.

### WIP
This package is still a work-in-progress. The lexer is pretty much there but the parsing tools are not completely implemented.

## Installation
```sh
go get github.com/cornejong/golex
```

Include the library in your project:
```go
import "github.com/cornejong/golex"
```


## Usage
### Basic Example
Here's an example of how to use the lexer to tokenize a simple source string:
```go
package main

import (
    "fmt"
    "github.com/cornejong/golex"
)

func main() {
    source := `func() { test = "SomeStringValue"; test = 1.2; test = 88 }`
    lexer := golex.NewLexer()

    for token, err := range lexer.Iterate(source) {
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        token.Dump()
    }
}

// Output:
//   1:   2 -> Symbol                func                  (<nil>)
//   1:   6 -> OpenParenthesis       (                     (<nil>)
//   1:   7 -> CloseParenthesis      )                     (<nil>)
//   1:   9 -> OpenCurlyBracket      {                     (<nil>)
//   1:  11 -> Symbol                test                  (<nil>)
//   1:  16 -> Assign                =                     (<nil>)
//   1:  18 -> DoubleQuoteString     "SomeStringValue"     (SomeStringValue)
//   1:  35 -> Semicolon             ;                     (<nil>)
//   1:  37 -> Symbol                test                  (<nil>)
//   1:  42 -> Assign                =                     (<nil>)
//   1:  44 -> Float                 1.2                   (1.2)
//   1:  47 -> Semicolon             ;                     (<nil>)
//   1:  49 -> Symbol                test                  (<nil>)
//   1:  54 -> Assign                =                     (<nil>)
//   1:  56 -> Integer               88                    (88)
//   1:  59 -> CloseCurlyBracket     }                     (<nil>)
//   1:  60 -> EndOfFile                                   (<nil>)
```

## Lexer Options
```go
lexer := NewLexer(
    // Print each token as it is parsed 
    DebugPrintTokens(),

    // Don't add the token position to the token
    OmitTokenPosition(),

    // Ignore specific tokens. Tokens will be parsed but lexer.NextToken will be returned
    IgnoreTokens(TypeComment),

    // Retain whitespace tokens
    RetainWhitespace(),

    // Turn symbols into keyword tokens
    WithKeywords("func", "const", "def"),

    // Specify the symbol character maps
    // - arg1: the start character of a symbol
    // - arg2: the continuation of the symbol
    SymbolCharacterMap("a-zA-Z_", "a-zA-Z0-9_"),

    // Register a custom tokenizer
    WithTokenizer(InsertBefore(TypeStringTokenizer, TokenizerType("MyCustomTokenizer"), MyCustomTokenizer{})),

    // Extend the literal tokens
    WithLiteralTokens(LiteralToken{Type: Type("MyLiteralToken", Literal: "__!__")}),

    // unset a build-in literal token
    WithoutLiteralTokens(TypeEllipses, TypeSemicolon),

    // Add a comment syntax
    WithCommentSyntax(CommentSyntax{Opener: "#"}, CommentSyntax{Opener: "/*", Closer: "*/"}),

    // Unset a build-in comment syntax
    WithoutCommentSyntax(CommentSyntax{Opener: "//"}),

    // Add a string enclosure
    WithStringEnclosure(StringEnclosure{Enclosure: "```"}),
    
    // Unset a build-in enclosure
    WithoutStringEnclosure(StringEnclosure{Enclosure: "\""})
)
```

## Tokens

```go
type Token struct {
    // The token Type
    Type     TokenType
    // The literal representation of the token
    Literal  string
    // The parsed value (if available)
    // Currently just for strings, numbers and booleans
    Value    any
    // The token Position within the source
    Position Position
}
```

## Build-in Types
All basic token types are build-in and can be unset or extended using the lexer options.
For a full list of build-in types check [build_in_types.go](build_in_types.go)



## TODO:
- [ ] Better lexer errors with positional info
- [ ] 