package interpreter

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/term"
)

func TestInterpreter(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"+ 1 1", "2"},
		{"* (+ 1 2) 3", "9"},
		{"/ (+ 1 2) 3", "1"},
		{"(/ (- 6 1) 2)", "2.5"},
		{"= 1 1", "true"},
	}

	interp := MustNewInterpreter()

	for _, test := range tests {
		require.Equal(t, test.expected, interp.MustLexAndEvaluate(test.input).Text, "Error in test `%s'", test.input)
	}
}

func TestInterpreterInvalidTerm(t *testing.T) {
	interp := MustNewInterpreter()
	require.Error(t, interp.Evaluate(nil))
	require.Error(t, interp.Evaluate(&term.Term{}))
}

func BenchmarkInterpreter(b *testing.B) {
	tests := []struct {
		input    string
		expected string
	}{
		{"(+ 1 1)", "2"},
		{"(* (+ 1 2) 3)", "9"},
		{"(/ (+ 1 2) 3)", "1"},
		{"(/ (- 6 1) 2)", "2.5"},
		{"(= 1 1)", "true"},
	}

	interp := MustNewInterpreter()

	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			require.Equal(b, test.expected, interp.MustLexAndEvaluate(test.input).Text)
		}
	}
}

func TestFoo(t *testing.T) {
	interp := MustNewInterpreter()
	interp.Logger = log.New(os.Stdout, "", log.LstdFlags)
	interp.MustLexAndEvaluate("(+ 1 1)")
}
