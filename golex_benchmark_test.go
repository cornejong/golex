package golex

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func BenchmarkGolex(b *testing.B) {
	lexer := NewLexer(WithKeywords("fun", "func", "def"))
	src := " func() { test = \"SomeStringValue\"; test = 1.2; test = 88 }\n"

	// lexer.TokenizeToSlice(strings.Repeat(src, 0000))

	for i := 0; i < b.N; i++ {
		lexer.TokenizeToSlice(strings.Repeat(src, 1000))
	}
}

func TestGolexSpeed(t *testing.T) {
	lines := 100000
	src := " func() { test = \"SomeStringValue\"; test = 1.2; test = 88 }"
	srcLong := strings.Repeat(src, lines)
	lexer := NewLexer(WithKeywords("fun", "func", "def"))

	start := time.Now()
	lexer.TokenizeToSlice(srcLong)
	fmt.Printf("Pared %s lines in %0.4fs\n", formatWithThousandsSeparator(lines), time.Since(start).Seconds())
}

func formatWithThousandsSeparator(num int) string {
	numStr := fmt.Sprint(num)
	length := len(numStr)
	if length <= 3 {
		return numStr
	}

	var result strings.Builder
	for i, digit := range numStr {
		if (length-i)%3 == 0 && i != 0 {
			result.WriteString(".")
		}
		result.WriteRune(digit)
	}

	return result.String()
}
