package interpreter_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
	"github.com/talon-one/talang/lexer"
	helpers "github.com/talon-one/talang/testhelpers"
)

func TestTemplate(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()

	require.NoError(t, interp.RegisterTemplate(shared.TaTemplate{
		CommonSignature: shared.CommonSignature{
			Name:    "Template1",
			Returns: block.DecimalKind,
		},
		Template: *lexer.MustLex("(* 2 (. Variable1))"),
	}))

	var result *block.Block

	require.NoError(t, interp.GenericSet("Variable1", 1))
	result = interp.MustLexAndEvaluate("(+ 1 (! Template1))")
	require.Equal(t, true, result.IsDecimal())
	require.Equal(t, "3", result.String)

	require.NoError(t, interp.GenericSet("Variable1", 2))
	result = interp.MustLexAndEvaluate("(+ 1 (! Template1))")
	require.Equal(t, true, result.IsDecimal())
	require.Equal(t, "5", result.String)
}

func TestFormatedTemplate(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	require.NoError(t, interp.RegisterTemplate(shared.TaTemplate{
		CommonSignature: shared.CommonSignature{
			Name: "MultiplyWith2",
			Arguments: []block.Kind{
				block.DecimalKind,
			},
			Returns: block.DecimalKind,
		},
		Template: *lexer.MustLex("(* 2 (# 0))"),
	}))

	var result *block.Block

	result = interp.MustLexAndEvaluate("(+ 1 (! MultiplyWith2 2))")
	require.Equal(t, true, result.IsDecimal())
	require.Equal(t, "5", result.String)

	result = interp.MustLexAndEvaluate("(+ 1 (! MultiplyWith2 4))")
	require.Equal(t, true, result.IsDecimal())
	require.Equal(t, "9", result.String)
}

func TestInvalidTemplateArgumentTypes(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	require.NoError(t, interp.RegisterTemplate(shared.TaTemplate{
		CommonSignature: shared.CommonSignature{
			Name: "MultiplyWith2",
			Arguments: []block.Kind{
				block.DecimalKind,
			},
			Returns: block.DecimalKind,
		},
		Template: *lexer.MustLex("(* 2 (# 0))"),
	}))

	require.Error(t, getError(interp.LexAndEvaluate("! MultiplyWith2 A")))
}

// // test if children got an error
// func TestInvalidTemplateArguments(t *testing.T) {
// 	interp := helpers.MustNewInterpreterWithLogger()
// 	require.NoError(t, interp.RegisterTemplate(shared.TaTemplate{
// 		CommonSignature: shared.CommonSignature{
// 			Name: "MultiplyWith2",
// 			Arguments: []block.Kind{
// 				block.DecimalKind,
// 			},
// 			Returns: block.DecimalKind,
// 		},
// 		Template: *lexer.MustLex("(* 2 (# 0))"),
// 	}))

// 	interp.RegisterFunction(shared.TaFunction{
// 		CommonSignature: shared.CommonSignature{
// 			Name: "FN",
// 		},
// 		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
// 			return nil, errors.New("SomeError")
// 		},
// 	})

// 	require.Error(t, getError(interp.LexAndEvaluate("! MultiplyWith2 (FN)")))
// }
