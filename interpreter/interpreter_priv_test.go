package interpreter

import (
	"log"
	"os"
	"testing"

	"github.com/ericlagergren/decimal"
	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
	"github.com/talon-one/talang/lexer"
)

func TestMatchesSignature(t *testing.T) {
	interp := MustNewInterpreter()
	interp.Logger = log.New(os.Stdout, "", log.LstdFlags)

	// 1 parameter
	t.Run("ShouldMatch(1)", func(t *testing.T) {
		matches, notMatching, evaluatedChildren, err := interp.matchesSignature(
			&shared.CommonSignature{
				Name: "fn",
				Arguments: []block.Kind{
					block.DecimalKind,
				},
				IsVariadic: false,
			},
			"fn",
			[]*block.Block{
				block.NewDecimal(decimal.New(0, 0)),
			},
		)
		require.Equal(t, true, matches)
		require.Equal(t, notMatchingDetail(0), notMatching)
		require.EqualValues(t, []*block.Block{
			block.NewDecimal(decimal.New(0, 0)),
		}, evaluatedChildren)
		require.NoError(t, err)
	})

	// 2 parameters
	t.Run("ShouldMatch(2)", func(t *testing.T) {
		matches, notMatching, evaluatedChildren, err := interp.matchesSignature(
			&shared.CommonSignature{
				Name: "fn",
				Arguments: []block.Kind{
					block.DecimalKind,
					block.StringKind,
				},
				IsVariadic: false,
			},
			"fn",
			[]*block.Block{
				block.NewDecimal(decimal.New(0, 0)),
				block.NewString("Hello"),
			},
		)
		require.Equal(t, true, matches)
		require.Equal(t, notMatchingDetail(0), notMatching)
		require.EqualValues(t, []*block.Block{
			block.NewDecimal(decimal.New(0, 0)),
			block.NewString("Hello"),
		}, evaluatedChildren)
		require.NoError(t, err)
	})

	// evaluate parameter
	t.Run("ShouldMatch(3)", func(t *testing.T) {
		matches, notMatching, evaluatedChildren, err := interp.matchesSignature(
			&shared.CommonSignature{
				Name: "fn",
				Arguments: []block.Kind{
					block.DecimalKind,
					block.StringKind,
				},
				IsVariadic: false,
			},
			"fn",
			[]*block.Block{
				lexer.MustLex("+ 1 2"),
				block.NewString("Hello"),
			},
		)
		require.Equal(t, true, matches)
		require.Equal(t, notMatchingDetail(0), notMatching)
		require.EqualValues(t, []*block.Block{
			block.NewDecimal(decimal.New(3, 0)),
			block.NewString("Hello"),
		}, evaluatedChildren)
		require.NoError(t, err)
	})

	// first parameter does not match type
	t.Run("ShouldNotMatch(1)", func(t *testing.T) {
		matches, notMatching, evaluatedChildren, err := interp.matchesSignature(
			&shared.CommonSignature{
				Name: "fn",
				Arguments: []block.Kind{
					block.StringKind,
				},
				IsVariadic: false,
			},
			"fn",
			[]*block.Block{
				block.NewDecimal(decimal.New(0, 0)),
			},
		)
		require.Equal(t, false, matches)
		require.Equal(t, invalidSignature, notMatching)
		require.Nil(t, evaluatedChildren)
		require.NoError(t, err)
	})

	// second parameter does not match type
	t.Run("ShouldNotMatch(2)", func(t *testing.T) {
		matches, notMatching, evaluatedChildren, err := interp.matchesSignature(
			&shared.CommonSignature{
				Name: "fn",
				Arguments: []block.Kind{
					block.StringKind,
					block.DecimalKind,
				},
				IsVariadic: false,
			},
			"fn",
			[]*block.Block{
				block.NewString("Hello"),
				block.NewString("World"),
			},
		)
		require.Equal(t, false, matches)
		require.Equal(t, invalidSignature, notMatching)
		require.Nil(t, evaluatedChildren)
		require.NoError(t, err)
	})

	// to few arguments
	t.Run("ShouldNotMatch(3)", func(t *testing.T) {
		matches, notMatching, evaluatedChildren, err := interp.matchesSignature(
			&shared.CommonSignature{
				Name: "fn",
				Arguments: []block.Kind{
					block.StringKind,
					block.StringKind,
					block.StringKind,
				},
				IsVariadic: false,
			},
			"fn",
			[]*block.Block{
				block.NewString("Hello"),
				block.NewString("World"),
			},
		)
		require.Equal(t, false, matches)
		require.Equal(t, invalidSignature, notMatching)
		require.Nil(t, evaluatedChildren)
		require.NoError(t, err)
	})

	// variadic signatures
	// same parameters
	t.Run("ShouldMatchVariadic(1)", func(t *testing.T) {
		matches, notMatching, evaluatedChildren, err := interp.matchesSignature(
			&shared.CommonSignature{
				Name: "fn",
				Arguments: []block.Kind{
					block.DecimalKind,
				},
				IsVariadic: true,
			},
			"fn",
			[]*block.Block{
				block.NewDecimal(decimal.New(0, 0)),
			},
		)
		require.Equal(t, true, matches)
		require.Equal(t, notMatchingDetail(0), notMatching)
		require.EqualValues(t, []*block.Block{
			block.NewDecimal(decimal.New(0, 0)),
		}, evaluatedChildren)
		require.NoError(t, err)
	})

	// variadic signatures
	// multiple parameters
	t.Run("ShouldMatchVariadic(2)", func(t *testing.T) {
		matches, notMatching, evaluatedChildren, err := interp.matchesSignature(
			&shared.CommonSignature{
				Name: "fn",
				Arguments: []block.Kind{
					block.DecimalKind,
				},
				IsVariadic: true,
			},
			"fn",
			[]*block.Block{
				block.NewDecimal(decimal.New(0, 0)),
				block.NewDecimal(decimal.New(0, 0)),
			},
		)
		require.Equal(t, true, matches)
		require.Equal(t, notMatchingDetail(0), notMatching)
		require.EqualValues(t, []*block.Block{
			block.NewDecimal(decimal.New(0, 0)),
			block.NewDecimal(decimal.New(0, 0)),
		}, evaluatedChildren)
		require.NoError(t, err)
	})

	// variadic signatures
	// multiple parameters
	t.Run("ShouldMatchVariadic(3)", func(t *testing.T) {
		matches, notMatching, evaluatedChildren, err := interp.matchesSignature(
			&shared.CommonSignature{
				Name: "fn",
				Arguments: []block.Kind{
					block.StringKind,
					block.DecimalKind,
				},
				IsVariadic: true,
			},
			"fn",
			[]*block.Block{
				block.NewString("Hello"),
				block.NewDecimal(decimal.New(0, 0)),
				block.NewDecimal(decimal.New(0, 0)),
			},
		)
		require.Equal(t, true, matches)
		require.Equal(t, notMatchingDetail(0), notMatching)
		require.EqualValues(t, []*block.Block{
			block.NewString("Hello"),
			block.NewDecimal(decimal.New(0, 0)),
			block.NewDecimal(decimal.New(0, 0)),
		}, evaluatedChildren)
		require.NoError(t, err)
	})

	// evaluate
	t.Run("ShouldMatchVariadic(4)", func(t *testing.T) {
		matches, notMatching, evaluatedChildren, err := interp.matchesSignature(
			&shared.CommonSignature{
				Name: "fn",
				Arguments: []block.Kind{
					block.StringKind,
					block.DecimalKind,
				},
				IsVariadic: true,
			},
			"fn",
			[]*block.Block{
				block.NewString("Hello"),
				lexer.MustLex("+ 1 2"),
			},
		)
		require.Equal(t, true, matches)
		require.Equal(t, notMatchingDetail(0), notMatching)
		require.EqualValues(t, []*block.Block{
			block.NewString("Hello"),
			block.NewDecimal(decimal.New(3, 0)),
		}, evaluatedChildren)
		require.NoError(t, err)
	})

	// Test KindTypes
	t.Run("MatchType(1)", func(t *testing.T) {
		matches, notMatching, evaluatedChildren, err := interp.matchesSignature(
			&shared.CommonSignature{
				Name: "fn",
				Arguments: []block.Kind{
					block.AnyKind,
				},
				IsVariadic: false,
			},
			"fn",
			[]*block.Block{
				block.NewString("Hello"),
			},
		)
		require.Equal(t, true, matches)
		require.Equal(t, notMatchingDetail(0), notMatching)
		require.EqualValues(t, []*block.Block{
			block.NewString("Hello"),
		}, evaluatedChildren)
		require.NoError(t, err)
	})

	t.Run("MatchType(2)", func(t *testing.T) {
		matches, notMatching, evaluatedChildren, err := interp.matchesSignature(
			&shared.CommonSignature{
				Name: "fn",
				Arguments: []block.Kind{
					block.AtomKind,
				},
				IsVariadic: false,
			},
			"fn",
			[]*block.Block{
				block.NewString("Hello"),
			},
		)
		require.Equal(t, true, matches)
		require.Equal(t, notMatchingDetail(0), notMatching)
		require.EqualValues(t, []*block.Block{
			block.NewString("Hello"),
		}, evaluatedChildren)
		require.NoError(t, err)
	})

	t.Run("MatchType(3)", func(t *testing.T) {
		matches, notMatching, evaluatedChildren, err := interp.matchesSignature(
			&shared.CommonSignature{
				Name: "fn",
				Arguments: []block.Kind{
					block.CollectionKind,
				},
				IsVariadic: false,
			},
			"fn",
			[]*block.Block{
				block.NewString("Hello"),
			},
		)
		require.Equal(t, false, matches)
		require.Equal(t, invalidSignature, notMatching)
		require.Nil(t, evaluatedChildren)
		require.NoError(t, err)
	})

	t.Run("MatchType(3)", func(t *testing.T) {
		matches, notMatching, evaluatedChildren, err := interp.matchesSignature(
			&shared.CommonSignature{
				Name: "fn",
				Arguments: []block.Kind{
					block.BlockKind,
				},
				IsVariadic: false,
			},
			"fn",
			[]*block.Block{
				lexer.MustLex("+ 1 2"),
			},
		)
		require.Equal(t, true, matches)
		require.Equal(t, notMatchingDetail(0), notMatching)
		require.EqualValues(t, []*block.Block{
			lexer.MustLex("+ 1 2"),
		}, evaluatedChildren)
		require.NoError(t, err)
	})
}
