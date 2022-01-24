package main

type NopWriter struct {
}

func (_ *NopWriter) Write(b []byte) (n int, err error) {
	return len(b), nil
}
