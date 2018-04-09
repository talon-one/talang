package interpreter

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/talon-one/talang/lexer"
	"github.com/talon-one/talang/token"
)

func TestFuncToRunWalkerNonVariadic(t *testing.T) {
	interp := MustNewInterpreter()
	interp.Logger = log.New(os.Stdout, "", log.LstdFlags)
	require.NoError(t, interp.RemoveAllFunctions())

	tests := []struct {
		Input    *token.TaToken
		Register []CommonSignature
		Expected []CommonSignature
	}{

		// MATCHING SIGNATURES
		// 1 parameter
		{
			Input: lexer.MustLex("fn 0"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(Decimal)"),
			},
			Expected: []CommonSignature{
				MustNewCommonSignature("fn(Decimal)"),
			},
		},

		// 2 parameters
		{
			Input: lexer.MustLex("fn Hello 0"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(String,Decimal)"),
			},
			Expected: []CommonSignature{
				MustNewCommonSignature("fn(String,Decimal)"),
			},
		},

		// AnyKind(Decimal)
		{
			Input: lexer.MustLex("fn 0"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(Any)"),
			},
			Expected: []CommonSignature{
				MustNewCommonSignature("fn(Any)"),
			},
		},

		// AnyKind(String)
		{
			Input: lexer.MustLex("fn Hello"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(Any)"),
			},
			Expected: []CommonSignature{
				MustNewCommonSignature("fn(Any)"),
			},
		},

		// AtomKind(Decimal)
		{
			Input: lexer.MustLex("fn 0"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(Atom)"),
			},
			Expected: []CommonSignature{
				MustNewCommonSignature("fn(Atom)"),
			},
		},

		// AtomKind(String)
		{
			Input: lexer.MustLex("fn Hello"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(Atom)"),
			},
			Expected: []CommonSignature{
				MustNewCommonSignature("fn(Atom)"),
			},
		},

		// Multiple possibilities
		{
			Input: lexer.MustLex("fn Hello"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(Atom)"),
				MustNewCommonSignature("fn(Decimal)"),
				MustNewCommonSignature("fn(String)"),
			},
			Expected: []CommonSignature{
				MustNewCommonSignature("fn(Atom)"),
				MustNewCommonSignature("fn(String)"),
			},
		},

		// Nested
		{
			Input: lexer.MustLex("fn (sub Hello)"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(String)"),
				MustNewCommonSignature("sub(String)String"),
			},
			Expected: []CommonSignature{
				MustNewCommonSignature("fn(String)"),
			},
		},

		// Nested multiple
		{
			Input: lexer.MustLex("fn (sub Hello)"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(String)"),
				MustNewCommonSignature("sub(String)Decimal"),
				MustNewCommonSignature("sub(Atom)String"),
			},
			Expected: []CommonSignature{
				MustNewCommonSignature("fn(String)"),
			},
		},

		// Nested multiple
		{
			Input: lexer.MustLex("fn (sub Hello)"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(String)"),
				MustNewCommonSignature("fn(Decimal)"),
				MustNewCommonSignature("fn(Token)"),
				MustNewCommonSignature("sub(Any)Any"),
			},
			Expected: []CommonSignature{
				MustNewCommonSignature("fn(String)"),
				MustNewCommonSignature("fn(Decimal)"),
				MustNewCommonSignature("fn(Token)"),
			},
		},

		// NOT MATCHING
		// Invalid name
		{
			Input: lexer.MustLex("fn1 Hello"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn2(String)"),
			},
			Expected: nil,
		},

		// Invalid count of arguments
		{
			Input: lexer.MustLex("fn Hello"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn()"),
			},
			Expected: nil,
		},

		// Invalid type
		{
			Input: lexer.MustLex("fn Hello"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(Decimal)"),
			},
			Expected: nil,
		},

		// Nested invalid type
		{
			Input: lexer.MustLex("fn (sub Hello)"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(String)"),
				MustNewCommonSignature("sub(String)Decimal"),
			},
			Expected: nil,
		},

		// Nested nested invalid type
		{
			Input: lexer.MustLex("fn (sub Hello)"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(String)"),
				MustNewCommonSignature("sub(Decimal)String"),
			},
			Expected: nil,
		},
	}

	for i, test := range tests {
		scope := interp.NewScope()
		for j := 0; j < len(test.Register); j++ {
			scope.MustRegisterFunction(TaFunction{CommonSignature: test.Register[j]})
		}
		walker := newFuncToRunWalker(scope, test.Input, 0)

		for j := 0; j < len(test.Expected); j++ {
			require.Equal(t, test.Expected[j].String(), walker.Next().CommonSignature.String(), "Test %d failed", i)
		}
		require.Nil(t, walker.Next())
	}
}

func TestFuncToRunWalkerVariadic(t *testing.T) {
	interp := MustNewInterpreter()
	interp.Logger = log.New(os.Stdout, "", log.LstdFlags)
	require.NoError(t, interp.RemoveAllFunctions())

	tests := []struct {
		Input    *token.TaToken
		Register []CommonSignature
		Expected []CommonSignature
	}{

		// MATCHING SIGNATURES
		// 0 parameters required, 0 parameter given
		{
			Input: token.NewToken("fn"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(Decimal...)"),
			},
			Expected: []CommonSignature{
				MustNewCommonSignature("fn(Decimal...)"),
			},
		},

		// 0 parameters required, 1 parameter given
		{
			Input: lexer.MustLex("fn 0"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(Decimal...)"),
			},
			Expected: []CommonSignature{
				MustNewCommonSignature("fn(Decimal...)"),
			},
		},

		// 1 parameter required, 2 parameter given
		{
			Input: lexer.MustLex("fn Hello 0"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(String, Decimal...)"),
			},
			Expected: []CommonSignature{
				MustNewCommonSignature("fn(String, Decimal...)"),
			},
		},

		// 1 parameter required, 3 parameters given
		{
			Input: lexer.MustLex("fn Hello 0 1"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(String, Decimal...)"),
			},
			Expected: []CommonSignature{
				MustNewCommonSignature("fn(String, Decimal...)"),
			},
		},

		// AnyKind(Decimal)
		{
			Input: lexer.MustLex("fn 0"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(Any...)"),
			},
			Expected: []CommonSignature{
				MustNewCommonSignature("fn(Any...)"),
			},
		},

		// AnyKind(String)
		{
			Input: lexer.MustLex("fn Hello"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(Any...)"),
			},
			Expected: []CommonSignature{
				MustNewCommonSignature("fn(Any...)"),
			},
		},

		// AtomKind(Decimal)
		{
			Input: lexer.MustLex("fn 0"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(Atom...)"),
			},
			Expected: []CommonSignature{
				MustNewCommonSignature("fn(Atom...)"),
			},
		},

		// AtomKind(String)
		{
			Input: lexer.MustLex("fn Hello"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(Atom...)"),
			},
			Expected: []CommonSignature{
				MustNewCommonSignature("fn(Atom...)"),
			},
		},

		// NOT MATCHING
		// Invalid name
		{
			Input: lexer.MustLex("fn1 Hello"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn2(String...)"),
			},
			Expected: nil,
		},
		// Invalid type
		{
			Input: lexer.MustLex("fn Hello"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(Decimal...)"),
			},
			Expected: nil,
		},

		// To few arguments
		{
			Input: lexer.MustLex("fn 0"),
			Register: []CommonSignature{
				MustNewCommonSignature("fn(Decimal, Decimal, Decimal...)"),
			},
			Expected: nil,
		},
	}

	for i, test := range tests {
		scope := interp.NewScope()
		for j := 0; j < len(test.Register); j++ {
			scope.MustRegisterFunction(TaFunction{CommonSignature: test.Register[j]})
		}
		walker := newFuncToRunWalker(scope, test.Input, 0)

		for j := 0; j < len(test.Expected); j++ {
			require.Equal(t, test.Expected[j].String(), walker.Next().CommonSignature.String(), "Test %d failed", i)
		}
		require.Nil(t, walker.Next())
	}
}

func TestAllFunctions(t *testing.T) {
	interp := Interpreter{
		Functions: []TaFunction{
			{
				CommonSignature: CommonSignature{
					Name: "Root1",
				},
			},
			{
				CommonSignature: CommonSignature{
					Name: "Root2",
				},
			},
		},
	}

	subInterp := Interpreter{
		Functions: []TaFunction{
			{
				CommonSignature: CommonSignature{
					Name: "Sub1",
				},
			},
			{
				CommonSignature: CommonSignature{
					Name: "Sub2",
				},
			},
		},
		Parent: &interp,
	}

	require.EqualValues(t, interp.Functions, interp.AllFunctions())
	require.EqualValues(t, append(subInterp.Functions, interp.Functions...), subInterp.AllFunctions())
}

func TestAllTemplates(t *testing.T) {
	interp := Interpreter{
		Templates: []TaTemplate{
			{
				CommonSignature: CommonSignature{
					Name: "Root1",
				},
			},
			{
				CommonSignature: CommonSignature{
					Name: "Root2",
				},
			},
		},
	}

	subInterp := Interpreter{
		Templates: []TaTemplate{
			{
				CommonSignature: CommonSignature{
					Name: "Sub1",
				},
			},
			{
				CommonSignature: CommonSignature{
					Name: "Sub2",
				},
			},
		},
		Parent: &interp,
	}

	require.EqualValues(t, interp.Templates, interp.AllTemplates())
	require.EqualValues(t, append(subInterp.Templates, interp.Templates...), subInterp.AllTemplates())
}
