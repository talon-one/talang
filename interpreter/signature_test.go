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
				block.StringKind,
				block.DecimalKind,
				block.BoolKind,
			},
			IsVariadic: false,
		}

		require.Equal(t, true, sig.MatchesArguments([]block.Kind{
			block.StringKind,
			block.DecimalKind,
			block.BoolKind,
		}))
	})

	t.Run("Unequal Arguments", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []block.Kind{
				block.StringKind,
				block.StringKind,
				block.BoolKind,
			},
			IsVariadic: false,
		}

		require.Equal(t, false, sig.MatchesArguments([]block.Kind{
			block.StringKind,
			block.DecimalKind,
			block.BoolKind,
		}))
	})

	t.Run("Unequal len of Arguments", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []block.Kind{
				block.StringKind,
				block.BoolKind,
			},
			IsVariadic: false,
		}

		require.Equal(t, false, sig.MatchesArguments([]block.Kind{
			block.StringKind,
			block.DecimalKind,
			block.BoolKind,
		}))
	})

	t.Run("Variadic Arguments", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []block.Kind{
				block.StringKind,
				block.BoolKind,
			},
			IsVariadic: true,
		}

		require.Equal(t, true, sig.MatchesArguments([]block.Kind{
			block.StringKind,
			block.BoolKind,
			block.BoolKind,
			block.BoolKind,
		}))
	})

	t.Run("Variadic Invalid Arguments", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []block.Kind{
				block.StringKind,
				block.BoolKind,
			},
			IsVariadic: true,
		}

		require.Equal(t, false, sig.MatchesArguments([]block.Kind{
			block.StringKind,
			block.BoolKind,
			block.BoolKind,
			block.StringKind,
		}))
	})

	t.Run("Variadic Invalid Arguments (2)", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []block.Kind{
				block.StringKind,
				block.BoolKind,
			},
			IsVariadic: true,
		}

		require.Equal(t, false, sig.MatchesArguments([]block.Kind{
			block.DecimalKind,
			block.BoolKind,
			block.BoolKind,
		}))
	})
	t.Run("AnyKind", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []block.Kind{
				block.AnyKind,
			},
		}

		require.Equal(t, true, sig.MatchesArguments([]block.Kind{
			block.DecimalKind,
		}))
		require.Equal(t, true, sig.MatchesArguments([]block.Kind{
			block.StringKind,
		}))
		require.Equal(t, true, sig.MatchesArguments([]block.Kind{
			block.BlockKind,
		}))
	})

	t.Run("AtomKind", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []block.Kind{
				block.AtomKind,
			},
		}

		require.Equal(t, true, sig.MatchesArguments([]block.Kind{
			block.DecimalKind,
		}))
		require.Equal(t, true, sig.MatchesArguments([]block.Kind{
			block.StringKind,
		}))
		require.Equal(t, false, sig.MatchesArguments([]block.Kind{
			block.BlockKind,
		}))
	})
}
