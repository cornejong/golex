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

	for i := 0; i < b.N; i++ {
		lexer.TokenizeToSlice(strings.Repeat(src, 1000))
	}
}

func BenchmarkGolexSpeedLongInput(b *testing.B) {
	lines := 100000
	src := " func() { test = \"SomeStringValue\"; test = 1.2; test = 88 }\n"
	srcLong := strings.Repeat(src, lines)
	lexer := NewLexer(WithKeywords("fun", "func", "def"))
	total := time.Duration(0)

	for i := 0; i < b.N; i++ {
		start := time.Now()
		lexer.TokenizeToSlice(srcLong)
		total += time.Since(start)
	}

	fmt.Printf("Pared %s lines %d times in averagely %0.4fs\n", formatWithThousandsSeparator(lines), b.N, total.Seconds()/float64(b.N))
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
