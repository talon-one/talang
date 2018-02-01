package interpreter

type dummyWriter struct{}

func (*dummyWriter) Write([]byte) (int, error) {
	return 0, nil
}
