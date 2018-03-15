package interpreter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/token"
)

func TestMatchesArguments(t *testing.T) {
	t.Run("Equal Arguments", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []token.Kind{
				token.String,
				token.Decimal,
				token.Bool,
			},
			IsVariadic: false,
		}

		require.Equal(t, true, sig.MatchesArguments([]token.Kind{
			token.String,
			token.Decimal,
			token.Bool,
		}))
	})

	t.Run("Unequal Arguments", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []token.Kind{
				token.String,
				token.String,
				token.Bool,
			},
			IsVariadic: false,
		}

		require.Equal(t, false, sig.MatchesArguments([]token.Kind{
			token.String,
			token.Decimal,
			token.Bool,
		}))
	})

	t.Run("Unequal len of Arguments", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []token.Kind{
				token.String,
				token.Bool,
			},
			IsVariadic: false,
		}

		require.Equal(t, false, sig.MatchesArguments([]token.Kind{
			token.String,
			token.Decimal,
			token.Bool,
		}))
	})

	t.Run("Variadic Arguments", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []token.Kind{
				token.String,
				token.Bool,
			},
			IsVariadic: true,
		}

		require.Equal(t, true, sig.MatchesArguments([]token.Kind{
			token.String,
			token.Bool,
			token.Bool,
			token.Bool,
		}))
	})

	t.Run("Variadic Invalid Arguments", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []token.Kind{
				token.String,
				token.Bool,
			},
			IsVariadic: true,
		}

		require.Equal(t, false, sig.MatchesArguments([]token.Kind{
			token.String,
			token.Bool,
			token.Bool,
			token.String,
		}))
	})

	t.Run("Variadic Invalid Arguments (2)", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []token.Kind{
				token.String,
				token.Bool,
			},
			IsVariadic: true,
		}

		require.Equal(t, false, sig.MatchesArguments([]token.Kind{
			token.Decimal,
			token.Bool,
			token.Bool,
		}))
	})
	t.Run("AnyKind", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []token.Kind{
				token.Any,
			},
		}

		require.Equal(t, true, sig.MatchesArguments([]token.Kind{
			token.Decimal,
		}))
		require.Equal(t, true, sig.MatchesArguments([]token.Kind{
			token.String,
		}))
		require.Equal(t, true, sig.MatchesArguments([]token.Kind{
			token.Token,
		}))
	})

	t.Run("AtomKind", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []token.Kind{
				token.Atom,
			},
		}

		require.Equal(t, true, sig.MatchesArguments([]token.Kind{
			token.Decimal,
		}))
		require.Equal(t, true, sig.MatchesArguments([]token.Kind{
			token.String,
		}))
		require.Equal(t, false, sig.MatchesArguments([]token.Kind{
			token.Token,
		}))
	})
}
