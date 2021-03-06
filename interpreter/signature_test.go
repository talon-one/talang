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
				token.Boolean,
			},
			IsVariadic: false,
		}

		require.Equal(t, true, sig.MatchesArguments([]token.Kind{
			token.String,
			token.Decimal,
			token.Boolean,
		}))
	})

	t.Run("Unequal Arguments", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []token.Kind{
				token.String,
				token.String,
				token.Boolean,
			},
			IsVariadic: false,
		}

		require.Equal(t, false, sig.MatchesArguments([]token.Kind{
			token.String,
			token.Decimal,
			token.Boolean,
		}))
	})

	t.Run("Unequal len of Arguments", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []token.Kind{
				token.String,
				token.Boolean,
			},
			IsVariadic: false,
		}

		require.Equal(t, false, sig.MatchesArguments([]token.Kind{
			token.String,
			token.Decimal,
			token.Boolean,
		}))
	})

	t.Run("Variadic Arguments", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []token.Kind{
				token.String,
				token.Boolean,
			},
			IsVariadic: true,
		}

		require.Equal(t, true, sig.MatchesArguments([]token.Kind{
			token.String,
			token.Boolean,
			token.Boolean,
			token.Boolean,
		}))
	})

	t.Run("Variadic Invalid Arguments", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []token.Kind{
				token.String,
				token.Boolean,
			},
			IsVariadic: true,
		}

		require.Equal(t, false, sig.MatchesArguments([]token.Kind{
			token.String,
			token.Boolean,
			token.Boolean,
			token.String,
		}))
	})

	t.Run("Variadic Invalid Arguments (2)", func(t *testing.T) {
		sig := CommonSignature{
			Arguments: []token.Kind{
				token.String,
				token.Boolean,
			},
			IsVariadic: true,
		}

		require.Equal(t, false, sig.MatchesArguments([]token.Kind{
			token.Decimal,
			token.Boolean,
			token.Boolean,
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

func TestSignatureParse(t *testing.T) {
	require.EqualValues(t, &CommonSignature{
		Name:      "Plus",
		lowerName: "plus",
		Arguments: []token.Kind{
			token.Decimal,
			token.String,
		},
		IsVariadic: true,
		Returns:    token.Boolean,
	}, NewCommonSignature("Plus (Decimal,String...)Boolean"))

	require.EqualValues(t, &CommonSignature{
		Name:      "Plus",
		lowerName: "plus",
		Arguments: []token.Kind{
			token.Decimal,
			token.String,
		},
		IsVariadic: true,
		Returns:    token.Boolean,
	}, NewCommonSignature("Plus(Decimal,String...)Boolean"))

	require.EqualValues(t, &CommonSignature{
		Name:      "Plus",
		lowerName: "plus",
		Arguments: []token.Kind{
			token.Decimal,
			token.String,
		},
		IsVariadic: true,
	}, NewCommonSignature("Plus(Decimal,String...)"))

	require.Nil(t, NewCommonSignature("(Decimal,String...)"))
	require.Nil(t, NewCommonSignature(""))

	sig := CommonSignature{
		Name:      "Plus",
		lowerName: "plus",
		Arguments: []token.Kind{
			token.Decimal,
			token.String,
		},
		IsVariadic: true,
	}

	require.EqualValues(t, &sig, NewCommonSignature(sig.String()))
}

func TestSignatureString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"Plus(Decimal, String...)Boolean",
			"Plus(Decimal, String...)Boolean",
		},
		{
			"Plus(Decimal, String...)",
			"Plus(Decimal, String...)",
		},
		{
			"Plus(Decimal)",
			"Plus(Decimal)",
		},
		{
			"Plus()Decimal",
			"Plus()Decimal",
		},
		{
			"Plus()",
			"Plus()",
		},
	}

	for i, test := range tests {
		require.Equal(t, test.expected, NewCommonSignature(test.input).String(), "Test %d failed", i)
	}
}
