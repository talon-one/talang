package testhelpers

/*
type LastError struct {
}

func (LastError) Error() string { return "" }

type StackError struct {
	StackTrace error
}

func (StackError) Error() string { return "" }

type AnotherError struct {
	StackTrace error
}

func (AnotherError) Error() string { return "" }

func TestTestErrorStackTrace(t *testing.T) {
	// 1 level test
	TestErrorStackTrace(
		&LastError{},
		ErrorStackTrace{
			Error: &LastError{},
		},
	)

	TestErrorStackTrace(
		&StackError{},
		ErrorStackTrace{
			Error: &StackError{},
		},
	)

	// 2 level test
	TestErrorStackTrace(
		&StackError{
			StackTrace: &LastError{},
		},
		ErrorStackTrace{
			Error: &StackError{},
			ErrorStackTrace: &ErrorStackTrace{
				Error: &LastError{},
			},
		},
	)

	TestErrorStackTrace(
		&StackError{
			StackTrace: &StackError{},
		},
		ErrorStackTrace{
			Error: &StackError{},
			ErrorStackTrace: &ErrorStackTrace{
				Error: &StackError{},
			},
		},
	)

	require.Panics(t, func() {
		TestErrorStackTrace(
			&StackError{
				StackTrace: &LastError{},
			},
			ErrorStackTrace{
				Error: &AnotherError{},
				ErrorStackTrace: &ErrorStackTrace{
					Error: &LastError{},
				},
			},
		)
	})

	require.Panics(t, func() {
		TestErrorStackTrace(
			&StackError{
				StackTrace: &LastError{},
			},
			ErrorStackTrace{
				Error: &StackError{},
				ErrorStackTrace: &ErrorStackTrace{
					Error: &AnotherError{},
				},
			},
		)
	})
}
*/
