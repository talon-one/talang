package interpreter

import (
	"log"
	"os"
	"testing"

	"github.com/ericlagergren/decimal"
	"github.com/talon-one/talang/lexer"

	"github.com/stretchr/testify/require"

	"github.com/talon-one/talang/block"
)

func TestMatchesSignatureNonVariadic(t *testing.T) {
	interp := MustNewInterpreter()
	interp.Logger = log.New(os.Stdout, "", log.LstdFlags)

	type Expected struct {
		Matches           bool
		NotMatching       notMatchingDetail
		EvaluatedChildren []*block.TaToken
		Error             bool
	}

	type Result struct {
		Matches           bool
		NotMatching       notMatchingDetail
		EvaluatedChildren []*block.TaToken
		Error             error
	}

	makeResult := func(Matches bool, NotMatching notMatchingDetail, EvaluatedChildren []*block.TaToken, Error error) Result {
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
				[]*block.TaToken{
					block.NewDecimal(decimal.New(0, 0)),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []block.Kind{
						block.Decimal,
					},
					IsVariadic: false,
				},
				"fn",
				[]*block.TaToken{
					block.NewDecimal(decimal.New(0, 0)),
				},
			)),
		},

		// 2 parameters
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*block.TaToken{
					block.NewString("Hello"),
					block.NewDecimal(decimal.New(0, 0)),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []block.Kind{
						block.String,
						block.Decimal,
					},
					IsVariadic: false,
				},
				"fn",
				[]*block.TaToken{
					block.NewString("Hello"),
					block.NewDecimal(decimal.New(0, 0)),
				},
			)),
		},

		// AnyKind(Decimal)
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*block.TaToken{
					block.NewDecimal(decimal.New(0, 0)),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []block.Kind{
						block.Any,
					},
					IsVariadic: false,
				},
				"fn",
				[]*block.TaToken{
					block.NewDecimal(decimal.New(0, 0)),
				},
			)),
		},

		// AnyKind(String)
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*block.TaToken{
					block.NewString("Hello"),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []block.Kind{
						block.Any,
					},
					IsVariadic: false,
				},
				"fn",
				[]*block.TaToken{
					block.NewString("Hello"),
				},
			)),
		},

		// AtomKind(Decimal)
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*block.TaToken{
					block.NewDecimal(decimal.New(0, 0)),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []block.Kind{
						block.Atom,
					},
					IsVariadic: false,
				},
				"fn",
				[]*block.TaToken{
					block.NewDecimal(decimal.New(0, 0)),
				},
			)),
		},

		// AtomKind(String)
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*block.TaToken{
					block.NewString("Hello"),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []block.Kind{
						block.Atom,
					},
					IsVariadic: false,
				},
				"fn",
				[]*block.TaToken{
					block.NewString("Hello"),
				},
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
					Arguments: []block.Kind{
						block.String,
					},
					IsVariadic: false,
				},
				"fn2",
				[]*block.TaToken{
					block.NewString("Hello"),
				},
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
					Arguments:  []block.Kind{},
					IsVariadic: false,
				},
				"fn",
				[]*block.TaToken{
					block.NewString("Hello"),
				},
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
					Arguments: []block.Kind{
						block.Decimal,
					},
					IsVariadic: false,
				},
				"fn",
				[]*block.TaToken{
					block.NewString("Hello"),
				},
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
					Arguments: []block.Kind{
						block.Decimal,
					},
					IsVariadic: false,
				},
				"fn",
				[]*block.TaToken{
					lexer.MustLex("(panic)"),
				},
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
		EvaluatedChildren []*block.TaToken
		Error             bool
	}

	type Result struct {
		Matches           bool
		NotMatching       notMatchingDetail
		EvaluatedChildren []*block.TaToken
		Error             error
	}

	makeResult := func(Matches bool, NotMatching notMatchingDetail, EvaluatedChildren []*block.TaToken, Error error) Result {
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
				[]*block.TaToken{},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []block.Kind{
						block.Decimal,
					},
					IsVariadic: true,
				},
				"fn",
				[]*block.TaToken{},
			)),
		},
		// 0 parameters required, 1 parameter given
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*block.TaToken{
					block.NewDecimal(decimal.New(0, 0)),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []block.Kind{
						block.Decimal,
					},
					IsVariadic: true,
				},
				"fn",
				[]*block.TaToken{
					block.NewDecimal(decimal.New(0, 0)),
				},
			)),
		},

		// 1 parameters required, 2 parameters given
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*block.TaToken{
					block.NewString("Hello"),
					block.NewDecimal(decimal.New(0, 0)),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []block.Kind{
						block.String,
						block.Decimal,
					},
					IsVariadic: true,
				},
				"fn",
				[]*block.TaToken{
					block.NewString("Hello"),
					block.NewDecimal(decimal.New(0, 0)),
				},
			)),
		},

		// 1 parameters required, 3 parameters given
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*block.TaToken{
					block.NewString("Hello"),
					block.NewDecimal(decimal.New(0, 0)),
					block.NewDecimal(decimal.New(1, 0)),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []block.Kind{
						block.String,
						block.Decimal,
					},
					IsVariadic: true,
				},
				"fn",
				[]*block.TaToken{
					block.NewString("Hello"),
					block.NewDecimal(decimal.New(0, 0)),
					block.NewDecimal(decimal.New(1, 0)),
				},
			)),
		},

		// AnyKind(Decimal)
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*block.TaToken{
					block.NewDecimal(decimal.New(0, 0)),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []block.Kind{
						block.Any,
					},
					IsVariadic: true,
				},
				"fn",
				[]*block.TaToken{
					block.NewDecimal(decimal.New(0, 0)),
				},
			)),
		},

		// AnyKind(String)
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*block.TaToken{
					block.NewString("Hello"),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []block.Kind{
						block.Any,
					},
					IsVariadic: true,
				},
				"fn",
				[]*block.TaToken{
					block.NewString("Hello"),
				},
			)),
		},

		// AtomKind(Decimal)
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*block.TaToken{
					block.NewDecimal(decimal.New(0, 0)),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []block.Kind{
						block.Atom,
					},
					IsVariadic: true,
				},
				"fn",
				[]*block.TaToken{
					block.NewDecimal(decimal.New(0, 0)),
				},
			)),
		},

		// AtomKind(String)
		{
			Expected: Expected{
				true,
				notMatchingDetail(0),
				[]*block.TaToken{
					block.NewString("Hello"),
				},
				false,
			},
			Result: makeResult(interp.matchesSignature(
				&CommonSignature{
					Name:      "fn",
					lowerName: "fn",
					Arguments: []block.Kind{
						block.Atom,
					},
					IsVariadic: true,
				},
				"fn",
				[]*block.TaToken{
					block.NewString("Hello"),
				},
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
					Arguments: []block.Kind{
						block.String,
					},
					IsVariadic: true,
				},
				"fn2",
				[]*block.TaToken{
					block.NewString("Hello"),
				},
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
					Arguments: []block.Kind{
						block.Decimal,
					},
					IsVariadic: true,
				},
				"fn",
				[]*block.TaToken{
					block.NewString("Hello"),
				},
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
					Arguments: []block.Kind{
						block.Decimal,
					},
					IsVariadic: true,
				},
				"fn",
				[]*block.TaToken{
					lexer.MustLex("(panic)"),
				},
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
					Arguments: []block.Kind{
						block.Decimal,
						block.Decimal,
						block.Decimal,
					},
					IsVariadic: true,
				},
				"fn",
				[]*block.TaToken{
					block.NewDecimalFromString("1"),
				},
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
		Func: func(interp *Interpreter, args ...*block.TaToken) (*block.TaToken, error) {
			return nil, nil
		},
	}))

	interp = interp.NewScope()
	require.NoError(t, interp.RemoveAllFunctions())

	require.NoError(t, interp.RegisterFunction(TaFunction{
		CommonSignature: CommonSignature{
			Name: "Scope1FN",
		},
		Func: func(interp *Interpreter, args ...*block.TaToken) (*block.TaToken, error) {
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
			Func: func(interp *Interpreter, args ...*block.TaToken) (*block.TaToken, error) {
				return nil, nil
			},
		},
		TaFunction{
			CommonSignature: CommonSignature{
				Name: "Scope2FN2",
			},
			Func: func(interp *Interpreter, args ...*block.TaToken) (*block.TaToken, error) {
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
