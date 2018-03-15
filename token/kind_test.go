package token

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKindString(t *testing.T) {
	require.Equal(t, "Decimal", Decimal.String())
	require.Equal(t, "String", String.String())
	require.Equal(t, "Bool", Bool.String())
	require.Equal(t, "Time", Time.String())
	require.Equal(t, "Null", Null.String())
	require.Equal(t, "List", List.String())
	require.Equal(t, "Map", Map.String())
	require.Equal(t, "Token", Token.String())
	require.Equal(t, "Atom", Atom.String())
	require.Equal(t, "Collection", Collection.String())
	require.Equal(t, "Any", Any.String())

	require.Equal(t, "Decimal|String", (Decimal | String).String())

	unknown := Kind(300)
	require.Equal(t, "Bool|Time|List|Unknown(256)", unknown.String())
}

func TestKindFromString(t *testing.T) {
	require.Equal(t, Decimal, KindFromString("Decimal"))
	require.Equal(t, String, KindFromString("String"))
	require.Equal(t, Bool, KindFromString("Bool"))
	require.Equal(t, Time, KindFromString("Time"))
	require.Equal(t, Null, KindFromString("Null"))
	require.Equal(t, List, KindFromString("List"))
	require.Equal(t, Map, KindFromString("Map"))
	require.Equal(t, Token, KindFromString("Token"))
	require.Equal(t, Atom, KindFromString("Atom"))
	require.Equal(t, Collection, KindFromString("Collection"))
	require.Equal(t, Any, KindFromString("Any"))

	require.Equal(t, Bool|Time|List, KindFromString("Bool|Time|List|Unknown(256)"))
}
