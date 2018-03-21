package interpreter

import (
	"log"
	"os"
	"testing"

	"github.com/ericlagergren/decimal"
	"github.com/talon-one/talang/lexer"

	"github.com/stretchr/testify/require"

	"github.com/talon-one/talang/token"
)

func TestMatchesSignatureNonVariadic(t *testing.T) {
	interp := MustNewInterpreter()
	interp.Logger = log.New(os.Stdout, "", log.LstdFlags)

	type Expected struct {
		Matches           bool
		NotMatching       notMatchingDetail
		EvaluatedChildren []*token.TaToken
		Error             bool
	}

	type Result struct {
		Matches           bool
		NotMatching       notMatchingDetail
		EvaluatedChildren []*token.TaToken
		Error             error
	}

	makeResult := func(Matches bool, NotMatching notMatchingDetail, EvaluatedChildren []*token.TaToken, Error error) Result {
		return Result{
			Matches:           Matches,
			NotMatching:       NotMatching,
			EvaluatedChildren: EvaluatedChildren,
			Error:             Error,
		}
	}

	tests := []struct {
		Expected Expected
		Result   Result
	}{

		// MATCHING SIGNATURES
		// 1 parameter
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*token.TaToken{
					token.NewDecimal(decimal.New(0, 0)),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.Decimal,
					},
					IsVariadic: false,
				},
				"fn",
				[]*token.TaToken{
					token.NewDecimal(decimal.New(0, 0)),
				},
				0,
			)),
		},

		// 2 parameters
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*token.TaToken{
					token.NewString("Hello"),
					token.NewDecimal(decimal.New(0, 0)),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.String,
						token.Decimal,
					},
					IsVariadic: false,
				},
				"fn",
				[]*token.TaToken{
					token.NewString("Hello"),
					token.NewDecimal(decimal.New(0, 0)),
				},
				0,
			)),
		},

		// AnyKind(Decimal)
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*token.TaToken{
					token.NewDecimal(decimal.New(0, 0)),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.Any,
					},
					IsVariadic: false,
				},
				"fn",
				[]*token.TaToken{
					token.NewDecimal(decimal.New(0, 0)),
				},
				0,
			)),
		},

		// AnyKind(String)
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*token.TaToken{
					token.NewString("Hello"),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.Any,
					},
					IsVariadic: false,
				},
				"fn",
				[]*token.TaToken{
					token.NewString("Hello"),
				},
				0,
			)),
		},

		// AtomKind(Decimal)
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*token.TaToken{
					token.NewDecimal(decimal.New(0, 0)),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.Atom,
					},
					IsVariadic: false,
				},
				"fn",
				[]*token.TaToken{
					token.NewDecimal(decimal.New(0, 0)),
				},
				0,
			)),
		},

		// AtomKind(String)
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*token.TaToken{
					token.NewString("Hello"),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.Atom,
					},
					IsVariadic: false,
				},
				"fn",
				[]*token.TaToken{
					token.NewString("Hello"),
				},
				0,
			)),
		},

		// NOT MATCHING
		// invalid name
		{
			Expected: Expected{
				false,
				invalidName,
				nil,
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn1",
					lowerName: "fn1",
					Arguments: []token.Kind{
						token.String,
					},
					IsVariadic: false,
				},
				"fn2",
				[]*token.TaToken{
					token.NewString("Hello"),
				},
				0,
			)),
		},
		// invalid count of arguments
		{
			Expected: Expected{
				false,
				invalidSignature,
				nil,
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:       "fn",
					lowerName:  "fn",
					Arguments:  []token.Kind{},
					IsVariadic: false,
				},
				"fn",
				[]*token.TaToken{
					token.NewString("Hello"),
				},
				0,
			)),
		},
		// invalid type
		{
			Expected: Expected{
				false,
				invalidSignature,
				nil,
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.Decimal,
					},
					IsVariadic: false,
				},
				"fn",
				[]*token.TaToken{
					token.NewString("Hello"),
				},
				0,
			)),
		},
		// error in child
		{
			Expected: Expected{
				false,
				errorInChildrenEvaluation,
				nil,
				true,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.Decimal,
					},
					IsVariadic: false,
				},
				"fn",
				[]*token.TaToken{
					lexer.MustLex("(panic)"),
				},
				0,
			)),
		},
	}

	for i, test := range tests {
		require.Equal(t, test.Expected.Matches, test.Result.Matches, "Test #%d failed", i)
		require.Equal(t, test.Expected.NotMatching, test.Result.NotMatching, "Test #%d failed", i)
		require.EqualValues(t, test.Expected.EvaluatedChildren, test.Result.EvaluatedChildren, "Test #%d failed", i)
		if test.Expected.Error == true {
			require.Error(t, test.Result.Error, "Test %d failed", i)
		} else {
			require.NoError(t, test.Result.Error, "Test %d failed", i)
		}
	}
}

func TestMatchesSignatureVariadic(t *testing.T) {
	interp := MustNewInterpreter()
	interp.Logger = log.New(os.Stdout, "", log.LstdFlags)

	type Expected struct {
		Matches           bool
		NotMatching       notMatchingDetail
		EvaluatedChildren []*token.TaToken
		Error             bool
	}

	type Result struct {
		Matches           bool
		NotMatching       notMatchingDetail
		EvaluatedChildren []*token.TaToken
		Error             error
	}

	makeResult := func(Matches bool, NotMatching notMatchingDetail, EvaluatedChildren []*token.TaToken, Error error) Result {
		return Result{
			Matches:           Matches,
			NotMatching:       NotMatching,
			EvaluatedChildren: EvaluatedChildren,
			Error:             Error,
		}
	}

	tests := []struct {
		Expected Expected
		Result   Result
	}{

		// MATCHING SIGNATURES

		// 0 parameters required, 0 parameter given
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*token.TaToken{},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.Decimal,
					},
					IsVariadic: true,
				},
				"fn",
				[]*token.TaToken{},
				0,
			)),
		},
		// 0 parameters required, 1 parameter given
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*token.TaToken{
					token.NewDecimal(decimal.New(0, 0)),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.Decimal,
					},
					IsVariadic: true,
				},
				"fn",
				[]*token.TaToken{
					token.NewDecimal(decimal.New(0, 0)),
				},
				0,
			)),
		},

		// 1 parameters required, 2 parameters given
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*token.TaToken{
					token.NewString("Hello"),
					token.NewDecimal(decimal.New(0, 0)),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.String,
						token.Decimal,
					},
					IsVariadic: true,
				},
				"fn",
				[]*token.TaToken{
					token.NewString("Hello"),
					token.NewDecimal(decimal.New(0, 0)),
				},
				0,
			)),
		},

		// 1 parameters required, 3 parameters given
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*token.TaToken{
					token.NewString("Hello"),
					token.NewDecimal(decimal.New(0, 0)),
					token.NewDecimal(decimal.New(1, 0)),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.String,
						token.Decimal,
					},
					IsVariadic: true,
				},
				"fn",
				[]*token.TaToken{
					token.NewString("Hello"),
					token.NewDecimal(decimal.New(0, 0)),
					token.NewDecimal(decimal.New(1, 0)),
				},
				0,
			)),
		},

		// AnyKind(Decimal)
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*token.TaToken{
					token.NewDecimal(decimal.New(0, 0)),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.Any,
					},
					IsVariadic: true,
				},
				"fn",
				[]*token.TaToken{
					token.NewDecimal(decimal.New(0, 0)),
				},
				0,
			)),
		},

		// AnyKind(String)
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*token.TaToken{
					token.NewString("Hello"),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.Any,
					},
					IsVariadic: true,
				},
				"fn",
				[]*token.TaToken{
					token.NewString("Hello"),
				},
				0,
			)),
		},

		// AtomKind(Decimal)
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*token.TaToken{
					token.NewDecimal(decimal.New(0, 0)),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.Atom,
					},
					IsVariadic: true,
				},
				"fn",
				[]*token.TaToken{
					token.NewDecimal(decimal.New(0, 0)),
				},
				0,
			)),
		},

		// AtomKind(String)
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*token.TaToken{
					token.NewString("Hello"),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.Atom,
					},
					IsVariadic: true,
				},
				"fn",
				[]*token.TaToken{
					token.NewString("Hello"),
				},
				0,
			)),
		},

		// NOT MATCHING
		// invalid name
		{
			Expected: Expected{
				false,
				invalidName,
				nil,
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn1",
					lowerName: "fn1",
					Arguments: []token.Kind{
						token.String,
					},
					IsVariadic: true,
				},
				"fn2",
				[]*token.TaToken{
					token.NewString("Hello"),
				},
				0,
			)),
		},
		// invalid type
		{
			Expected: Expected{
				false,
				invalidSignature,
				nil,
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.Decimal,
					},
					IsVariadic: true,
				},
				"fn",
				[]*token.TaToken{
					token.NewString("Hello"),
				},
				0,
			)),
		},
		// error in child
		{
			Expected: Expected{
				false,
				errorInChildrenEvaluation,
				nil,
				true,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.Decimal,
					},
					IsVariadic: true,
				},
				"fn",
				[]*token.TaToken{
					lexer.MustLex("(panic)"),
				},
				0,
			)),
		},
		// to few arguments
		{
			Expected: Expected{
				false,
				invalidSignature,
				nil,
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []token.Kind{
						token.Decimal,
						token.Decimal,
						token.Decimal,
					},
					IsVariadic: true,
				},
				"fn",
				[]*token.TaToken{
					token.NewDecimalFromString("1"),
				},
				0,
			)),
		},
	}

	for i, test := range tests {
		require.Equal(t, test.Expected.NotMatching, test.Result.NotMatching, "Test #%d failed", i)
		require.Equal(t, test.Expected.Matches, test.Result.Matches, "Test #%d failed", i)
		require.EqualValues(t, test.Expected.EvaluatedChildren, test.Result.EvaluatedChildren, "Test #%d failed", i)
		if test.Expected.Error == true {
			require.Error(t, test.Result.Error, "Test #%d failed", i)
		} else {
			require.NoError(t, test.Result.Error, "Test #%d failed", i)
		}
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

func TestFuncWalker(t *testing.T) {
	interp, err := NewInterpreter()

	require.NoError(t, err)
	require.NoError(t, interp.RemoveAllFunctions())

	require.NoError(t, interp.RegisterFunction(TaFunction{
		CommonSignature: CommonSignature{
			Name: "ROOTFN",
		},
		Func: func(interp *Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return nil, nil
		},
	}))

	interp = interp.NewScope()
	require.NoError(t, interp.RemoveAllFunctions())

	require.NoError(t, interp.RegisterFunction(TaFunction{
		CommonSignature: CommonSignature{
			Name: "Scope1FN",
		},
		Func: func(interp *Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return nil, nil
		},
	}))

	interp = interp.NewScope()
	require.NoError(t, interp.RemoveAllFunctions())

	require.NoError(t, interp.RegisterFunction(
		TaFunction{
			CommonSignature: CommonSignature{
				Name: "Scope2FN1",
			},
			Func: func(interp *Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
				return nil, nil
			},
		},
		TaFunction{
			CommonSignature: CommonSignature{
				Name: "Scope2FN2",
			},
			Func: func(interp *Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
				return nil, nil
			},
		},
	))

	walker := funcWalker{
		interp: interp,
	}

	require.Equal(t, "Scope2FN1", walker.Next().Name)
	require.Equal(t, "Scope2FN2", walker.Next().Name)
	require.Equal(t, "Scope1FN", walker.Next().Name)
	require.Equal(t, "ROOTFN", walker.Next().Name)
	require.Nil(t, walker.Next())
}
