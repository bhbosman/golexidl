package main

type nullWriter struct {
}

func (self nullWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}
