package interpreter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
)

func TestMatchesArguments(t *testing.T) {
	t.Run("Equal Arguments", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []block.Kind{
				block.String,
				block.Decimal,
				block.Bool,
			},
			IsVariadic: false,
		}

		require.Equal(t, true, sig.MatchesArguments([]block.Kind{
			block.String,
			block.Decimal,
			block.Bool,
		}))
	})

	t.Run("Unequal Arguments", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []block.Kind{
				block.String,
				block.String,
				block.Bool,
			},
			IsVariadic: false,
		}

		require.Equal(t, false, sig.MatchesArguments([]block.Kind{
			block.String,
			block.Decimal,
			block.Bool,
		}))
	})

	t.Run("Unequal len of Arguments", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []block.Kind{
				block.String,
				block.Bool,
			},
			IsVariadic: false,
		}

		require.Equal(t, false, sig.MatchesArguments([]block.Kind{
			block.String,
			block.Decimal,
			block.Bool,
		}))
	})

	t.Run("Variadic Arguments", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []block.Kind{
				block.String,
				block.Bool,
			},
			IsVariadic: true,
		}

		require.Equal(t, true, sig.MatchesArguments([]block.Kind{
			block.String,
			block.Bool,
			block.Bool,
			block.Bool,
		}))
	})

	t.Run("Variadic Invalid Arguments", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []block.Kind{
				block.String,
				block.Bool,
			},
			IsVariadic: true,
		}

		require.Equal(t, false, sig.MatchesArguments([]block.Kind{
			block.String,
			block.Bool,
			block.Bool,
			block.String,
		}))
	})

	t.Run("Variadic Invalid Arguments (2)", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []block.Kind{
				block.String,
				block.Bool,
			},
			IsVariadic: true,
		}

		require.Equal(t, false, sig.MatchesArguments([]block.Kind{
			block.Decimal,
			block.Bool,
			block.Bool,
		}))
	})
	t.Run("AnyKind", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []block.Kind{
				block.Any,
			},
		}

		require.Equal(t, true, sig.MatchesArguments([]block.Kind{
			block.Decimal,
		}))
		require.Equal(t, true, sig.MatchesArguments([]block.Kind{
			block.String,
		}))
		require.Equal(t, true, sig.MatchesArguments([]block.Kind{
			block.Token,
		}))
	})

	t.Run("AtomKind", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []block.Kind{
				block.Atom,
			},
		}

		require.Equal(t, true, sig.MatchesArguments([]block.Kind{
			block.Decimal,
		}))
		require.Equal(t, true, sig.MatchesArguments([]block.Kind{
			block.String,
		}))
		require.Equal(t, false, sig.MatchesArguments([]block.Kind{
			block.Token,
		}))
	})
}
