package main

import "github.com/ugorji/go/codec"

var codecHandle = codec.MsgpackHandle{}

func Encode(data interface{}) ([]byte, error) {
	b := new([]byte)
	return *b, codec.NewEncoderBytes(b, &codecHandle).Encode(data)
}

func Decode(data []byte, result interface{}) error {
	return codec.NewDecoderBytes(data, &codecHandle).Decode(result)
}
