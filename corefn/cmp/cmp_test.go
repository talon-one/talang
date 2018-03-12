package cmp_test

import (
	"testing"

	"github.com/talon-one/talang/lexer"

	"github.com/talon-one/talang/block"
	helpers "github.com/talon-one/talang/testhelpers"
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
			block.NewBool(true),
		},
		helpers.Test{
			`(= "Hello World" "Hello World")`,
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			`(= true true)`,
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			`(= "2006-01-02T15:04:05Z" "2006-01-02T15:04:05Z")`,
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			`(= 1 "1")`,
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			`(= "Hello" "Bye")`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`(= "Hello" "Hello" "Bye")`,
			nil,
			block.NewBool(false),
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
			block.NewBool(false),
		},
		helpers.Test{
			`(!= "Hello World" "Hello World")`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`(!= true true)`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`(!= "2006-01-02T15:04:05Z" "2006-01-02T15:04:05Z")`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`(!= 1 "1")`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`(!= "Hello" "Bye")`,
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			`(!= "Hello" "Hello" "Bye")`,
			nil,
			block.NewBool(false),
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
			block.NewBool(false),
		},
		helpers.Test{
			`(> 1 1)`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`(> 2 1)`,
			nil,
			block.NewBool(true),
		},
	)
}

func TestGreaterThanTime(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`(> "2006-01-02T15:04:05Z")`,
			nil,
			lexer.MustLex(`> "2006-01-02T15:04:05Z"`),
		},
		helpers.Test{
			`(> "2006-01-02T15:04:05Z" "2007-01-02T15:04:05Z")`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`(> "2007-01-02T15:04:05Z" "2007-01-02T15:04:05Z")`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`(> "2008-01-02T15:04:05Z" "2007-01-02T15:04:05Z")`,
			nil,
			block.NewBool(true),
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
			block.NewBool(true),
		},
		helpers.Test{
			`(< 1 1)`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`(< 2 1)`,
			nil,
			block.NewBool(false),
		},
	)
}

func TestLessThanTime(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`(< "2006-01-02T15:04:05Z")`,
			nil,
			lexer.MustLex(`< "2006-01-02T15:04:05Z"`),
		},
		helpers.Test{
			`(< "2006-01-02T15:04:05Z" "2007-01-02T15:04:05Z")`,
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			`(< "2007-01-02T15:04:05Z" "2007-01-02T15:04:05Z")`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`(< "2008-01-02T15:04:05Z" "2007-01-02T15:04:05Z")`,
			nil,
			block.NewBool(false),
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
			block.NewBool(false),
		},
		helpers.Test{
			`(>= 1 1)`,
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			`(>= 2 1)`,
			nil,
			block.NewBool(true),
		},
	)
}

func TestGreaterThanOrEqualTime(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`(>= "2006-01-02T15:04:05Z")`,
			nil,
			lexer.MustLex(`>= "2006-01-02T15:04:05Z"`),
		},
		helpers.Test{
			`(>= "2006-01-02T15:04:05Z" "2007-01-02T15:04:05Z")`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`(>= "2007-01-02T15:04:05Z" "2007-01-02T15:04:05Z")`,
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			`(>= "2008-01-02T15:04:05Z" "2007-01-02T15:04:05Z")`,
			nil,
			block.NewBool(true),
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
			block.NewBool(true),
		},
		helpers.Test{
			`(<= 1 1)`,
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			`(<= 2 1)`,
			nil,
			block.NewBool(false),
		},
	)
}

func TestLessThanOrEqualTime(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`(<= "2006-01-02T15:04:05Z")`,
			nil,
			lexer.MustLex(`<= "2006-01-02T15:04:05Z"`),
		},
		helpers.Test{
			`(<= "2006-01-02T15:04:05Z" "2007-01-02T15:04:05Z")`,
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			`(<= "2007-01-02T15:04:05Z" "2007-01-02T15:04:05Z")`,
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			`(<= "2008-01-02T15:04:05Z" "2007-01-02T15:04:05Z")`,
			nil,
			block.NewBool(false),
		},
	)
}

func TestBetweenDecimal(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`(between 1 0 3)`,
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			`(between 1 2 0 3)`,
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			`(between 0 0 2)`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`(between 2 0 2)`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`(between 1 4 0 3)`,
			nil,
			block.NewBool(false),
		},
	)
}

func TestBetweenTime(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`(between "2007-01-02T00:00:00Z" "2006-01-02T00:00:00Z" "2009-01-02T00:00:00Z")`,
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			`(between "2007-01-02T00:00:00Z" "2008-01-02T00:00:00Z" "2006-01-02T00:00:00Z" "2009-01-02T00:00:00Z")`,
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			`(between "2006-01-02T00:00:00Z" "2006-01-02T00:00:00Z" "2008-01-02T00:00:00Z")`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`(between "2008-01-02T00:00:00Z" "2006-01-02T00:00:00Z" "2008-01-02T00:00:00Z")`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`(between "2007-01-02T00:00:00Z" "2010-01-02T00:00:00Z" "2006-01-02T00:00:00Z" "2009-01-02T00:00:00Z")`,
			nil,
			block.NewBool(false),
		},
	)
}
