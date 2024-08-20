package golex

type TokenizerType string

const (
	TypeNoTokenizer      TokenizerType = ""
	TypeCommentTokenizer TokenizerType = "BuildInCommentTokenizer"
	TypeStringTokenizer  TokenizerType = "BuildInStringTokenizer"
	TypeNumberTokenizer  TokenizerType = "BuildInNumberTokenizer"
	TypeLiteralTokenizer TokenizerType = "BuildInLiteralTokenizer"
	TypeSymbolTokenizer  TokenizerType = "BuildInSymbolTokenizer"
	TypeBooleanTokenizer TokenizerType = "BuildInBooleanTokenizer"
)

type Tokenizer interface {
	CanTokenize(*Lexer) bool
	Tokenize(*Lexer) Token
}

type TokenizerInserter struct {
	tokenizerType TokenizerType
	tokenizer     Tokenizer
	Before        TokenizerType
	After         TokenizerType
}

func (ti TokenizerInserter) Insert(tokenizers map[TokenizerType]Tokenizer, tokenizationOrder []TokenizerType) (map[TokenizerType]Tokenizer, []TokenizerType) {
	tokenizers[ti.tokenizerType] = ti.tokenizer

	for i, t := range tokenizationOrder {
		if t == ti.Before {
			tokenizationOrder = append(tokenizationOrder[:i], append([]TokenizerType{ti.tokenizerType}, tokenizationOrder[i:]...)...)
			return tokenizers, tokenizationOrder
		}

		if t == ti.After {
			tokenizationOrder = append(tokenizationOrder[:i+1], append([]TokenizerType{ti.tokenizerType}, tokenizationOrder[i+1:]...)...)
			return tokenizers, tokenizationOrder
		}
	}

	// The specified Before or After was not in the tokenization order,
	if ti.Before != TypeNoTokenizer {
		// The insertion point was specified as "Before", so lets prepend it to the list
		tokenizationOrder = append([]TokenizerType{ti.tokenizerType}, tokenizationOrder...)
	} else {
		// The insertion point was specified as "After", so lets append it to the list
		tokenizationOrder = append(tokenizationOrder, ti.tokenizerType)
	}

	return tokenizers, tokenizationOrder
}

func InsertBefore(before TokenizerType, tokenizerType TokenizerType, tokenizer Tokenizer) TokenizerInserter {
	return TokenizerInserter{
		tokenizerType: tokenizerType,
		tokenizer:     tokenizer,
		Before:        before,
	}
}

func InsertAfter(after TokenizerType, tokenizerType TokenizerType, tokenizer Tokenizer) TokenizerInserter {
	return TokenizerInserter{
		tokenizerType: tokenizerType,
		tokenizer:     tokenizer,
		After:         after,
	}
}
