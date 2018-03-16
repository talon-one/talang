package cmp_test

import (
	"testing"

	"github.com/talon-one/talang/lexer"

	helpers "github.com/talon-one/talang/testhelpers"
	"github.com/talon-one/talang/token"
)

func TestEqual(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`(= 1)`,
			nil,
			lexer.MustLex(`= 1`),
		},
		helpers.Test{
			`(= 1 1)`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`(= "Hello World" "Hello World")`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`(= true true)`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`(= 2006-01-02T15:04:05Z 2006-01-02T15:04:05Z)`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`(= 1 "1")`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`(= "Hello" "Bye")`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`(= "Hello" "Hello" "Bye")`,
			nil,
			token.NewBool(false),
		},
	)
}

func TestNotEqual(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`(!= 1)`,
			nil,
			lexer.MustLex(`!= 1`),
		},
		helpers.Test{
			`(!= 1 1)`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`(!= "Hello World" "Hello World")`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`(!= true true)`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`(!= 2006-01-02T15:04:05Z 2006-01-02T15:04:05Z)`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`(!= 1 "1")`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`(!= "Hello" "Bye")`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`(!= "Hello" "Hello" "Bye")`,
			nil,
			token.NewBool(false),
		},
	)
}

func TestGreaterThanDecimal(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`(> 0)`,
			nil,
			lexer.MustLex(`> 0`),
		},
		helpers.Test{
			`(> 0 1)`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`(> 1 1)`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`(> 2 1)`,
			nil,
			token.NewBool(true),
		},
	)
}

func TestGreaterThanTime(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`(> 2006-01-02T15:04:05Z)`,
			nil,
			lexer.MustLex(`> 2006-01-02T15:04:05Z`),
		},
		helpers.Test{
			`(> 2006-01-02T15:04:05Z 2007-01-02T15:04:05Z)`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`(> 2007-01-02T15:04:05Z 2007-01-02T15:04:05Z)`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`(> 2008-01-02T15:04:05Z 2007-01-02T15:04:05Z)`,
			nil,
			token.NewBool(true),
		},
	)
}

func TestLessThanDecimal(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`(< 0)`,
			nil,
			lexer.MustLex(`< 0`),
		},
		helpers.Test{
			`(< 0 1)`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`(< 1 1)`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`(< 2 1)`,
			nil,
			token.NewBool(false),
		},
	)
}

func TestLessThanTime(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`(< 2006-01-02T15:04:05Z)`,
			nil,
			lexer.MustLex(`< 2006-01-02T15:04:05Z`),
		},
		helpers.Test{
			`(< 2006-01-02T15:04:05Z 2007-01-02T15:04:05Z)`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`(< 2007-01-02T15:04:05Z 2007-01-02T15:04:05Z)`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`(< 2008-01-02T15:04:05Z 2007-01-02T15:04:05Z)`,
			nil,
			token.NewBool(false),
		},
	)
}

func TestGreaterThanOrEqualDecimal(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`(>= 0)`,
			nil,
			lexer.MustLex(`>= 0`),
		},
		helpers.Test{
			`(>= 0 1)`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`(>= 1 1)`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`(>= 2 1)`,
			nil,
			token.NewBool(true),
		},
	)
}

func TestGreaterThanOrEqualTime(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`(>= 2006-01-02T15:04:05Z)`,
			nil,
			lexer.MustLex(`>= 2006-01-02T15:04:05Z`),
		},
		helpers.Test{
			`(>= 2006-01-02T15:04:05Z 2007-01-02T15:04:05Z)`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`(>= 2007-01-02T15:04:05Z 2007-01-02T15:04:05Z)`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`(>= 2008-01-02T15:04:05Z 2007-01-02T15:04:05Z)`,
			nil,
			token.NewBool(true),
		},
	)
}

func TestLessThanOrEqualDecimal(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`(<= 0)`,
			nil,
			lexer.MustLex(`<= 0`),
		},
		helpers.Test{
			`(<= 0 1)`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`(<= 1 1)`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`(<= 2 1)`,
			nil,
			token.NewBool(false),
		},
	)
}

func TestLessThanOrEqualTime(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`(<= 2006-01-02T15:04:05Z)`,
			nil,
			lexer.MustLex(`<= 2006-01-02T15:04:05Z`),
		},
		helpers.Test{
			`(<= 2006-01-02T15:04:05Z 2007-01-02T15:04:05Z)`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`(<= 2007-01-02T15:04:05Z 2007-01-02T15:04:05Z)`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`(<= 2008-01-02T15:04:05Z 2007-01-02T15:04:05Z)`,
			nil,
			token.NewBool(false),
		},
	)
}

func TestBetweenDecimal(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`(between 1 0 3)`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`(between 1 2 0 3)`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`(between 0 0 2)`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`(between 2 0 2)`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`(between 1 4 0 3)`,
			nil,
			token.NewBool(false),
		},
	)
}

func TestBetweenTime(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`(between 2007-01-02T00:00:00Z 2006-01-02T00:00:00Z 2009-01-02T00:00:00Z)`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`(between 2007-01-02T00:00:00Z 2008-01-02T00:00:00Z 2006-01-02T00:00:00Z 2009-01-02T00:00:00Z)`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`(between 2006-01-02T00:00:00Z 2006-01-02T00:00:00Z 2008-01-02T00:00:00Z)`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`(between 2008-01-02T00:00:00Z 2006-01-02T00:00:00Z 2008-01-02T00:00:00Z)`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`(between 2007-01-02T00:00:00Z 2010-01-02T00:00:00Z 2006-01-02T00:00:00Z 2009-01-02T00:00:00Z)`,
			nil,
			token.NewBool(false),
		},
	)
}

func TestOr(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`or false false false false false`,
			nil,
			token.NewBool(false),
		}, helpers.Test{
			`or false false false true false`,
			nil,
			token.NewBool(true),
		}, helpers.Test{
			`or (. List True) false`,
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewMap(map[string]*token.TaToken{
					"True":  token.NewBool(true),
					"False": token.NewBool(false),
				}),
			}),
			token.NewBool(true),
		}, helpers.Test{
			`or false`,
			nil,
			token.NewBool(false),
		}, helpers.Test{
			`or true`,
			nil,
			token.NewBool(true),
		}, helpers.Test{
			`or (> 1 1)`,
			nil,
			token.NewBool(false),
		}, helpers.Test{
			`or (> 2 2) false true`,
			nil,
			token.NewBool(true),
		}, helpers.Test{
			`or (+ 2 2)`,
			nil,
			helpers.Error{},
		},
		// Collections evaluate to true
		helpers.Test{
			`or (list 1 2 3)`,
			nil,
			token.NewBool(true),
		}, helpers.Test{
			`or (map (list "World" "Universe") ((x) (+ "Hello " (. x))))`,
			nil,
			token.NewBool(true),
		},
	)
}

func TestAnd(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`and false false false false false`,
			nil,
			token.NewBool(false),
		}, helpers.Test{
			`and true true true`,
			nil,
			token.NewBool(true),
		}, helpers.Test{
			`and (. List True) true`,
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewMap(map[string]*token.TaToken{
					"True":  token.NewBool(true),
					"False": token.NewBool(false),
				}),
			}),
			token.NewBool(true),
		}, helpers.Test{
			`and false false`,
			nil,
			token.NewBool(false),
		}, helpers.Test{
			`and true`,
			nil,
			token.NewBool(true),
		}, helpers.Test{
			`and (> 1 1)`,
			nil,
			token.NewBool(false),
		}, helpers.Test{
			`and (> 2 1) true`,
			nil,
			token.NewBool(true),
		}, helpers.Test{
			`and (+ 2 2)`,
			nil,
			helpers.Error{},
		}, helpers.Test{
			`and (list 1 2 3) false`,
			nil,
			token.NewBool(false),
		}, helpers.Test{
			`and`,
			nil,
			token.NewBool(true),
		},
	)
}
