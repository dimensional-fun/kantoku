package main

type NopWriter struct{}

func (NopWriter) Write(b []byte) (n int, err error) {
	return len(b), nil
}
