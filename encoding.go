package main

import (
	"github.com/ugorji/go/codec"
)

var codecHandle = codec.MsgpackHandle{}

func Encode(data interface{}) ([]byte, error) {
	bytes := new([]byte)
	return *bytes, codec.NewEncoderBytes(bytes, &codecHandle).Encode(data)
}

func Decode(data []byte, result interface{}) error {
	return codec.NewDecoderBytes(data, &codecHandle).Decode(result)
}
