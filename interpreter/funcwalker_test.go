package interpreter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/lexer"
	"github.com/talon-one/talang/token"
)

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

func TestFuncToRunWalker(t *testing.T) {
	t.Run("Simple Overload", func(t *testing.T) {
		interp, err := NewInterpreter()
		require.NoError(t, err)
		require.NoError(t, interp.RemoveAllFunctions())

		interp.MustRegisterFunction(
			TaFunction{
				CommonSignature: MustNewCommonSignature("fn(String)String"),
				Func: func(interp *Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
					return nil, nil
				},
			},
			TaFunction{
				CommonSignature: MustNewCommonSignature("fn(Any)String"),
				Func: func(interp *Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
					return nil, nil
				},
			},
		)

		tkn := lexer.MustLex("fn Hello")

		walker := newFuncToRunWalker(interp, tkn, 0)

		require.Equal(t, "fn(String)String", walker.Next().CommonSignature.String())
		require.Equal(t, "fn(Any)String", walker.Next().CommonSignature.String())
		require.Nil(t, walker.Next())
	})

	t.Run("Variadic", func(t *testing.T) {
		interp, err := NewInterpreter()
		require.NoError(t, err)
		require.NoError(t, interp.RemoveAllFunctions())

		interp.MustRegisterFunction(
			TaFunction{
				CommonSignature: MustNewCommonSignature("fn(String...)String"),
				Func: func(interp *Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
					return nil, nil
				},
			},
			TaFunction{
				CommonSignature: MustNewCommonSignature("fn(Any...)String"),
				Func: func(interp *Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
					return nil, nil
				},
			},
		)

		tkn := lexer.MustLex("fn Hello")

		walker := newFuncToRunWalker(interp, tkn, 0)

		require.Equal(t, "fn(String...)String", walker.Next().CommonSignature.String())
		require.Equal(t, "fn(Any...)String", walker.Next().CommonSignature.String())
		require.Nil(t, walker.Next())
	})

	t.Run("Block", func(t *testing.T) {
		interp, err := NewInterpreter()
		require.NoError(t, err)
		require.NoError(t, interp.RemoveAllFunctions())

		interp.MustRegisterFunction(
			TaFunction{
				CommonSignature: MustNewCommonSignature("+(String...)Decimal"),
				Func: func(interp *Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
					return nil, nil
				},
			},
			TaFunction{
				CommonSignature: MustNewCommonSignature("+(String, String)String"),
				Func: func(interp *Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
					return nil, nil
				},
			},
			TaFunction{
				CommonSignature: MustNewCommonSignature("fn(String...)String"),
				Func: func(interp *Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
					return nil, nil
				},
			},
		)

		tkn := lexer.MustLex("fn (+ Hello World)")

		walker := newFuncToRunWalker(interp, tkn, 0)

		require.Equal(t, "fn(String...)String", walker.Next().CommonSignature.String())
		require.Nil(t, walker.Next())
	})
}
