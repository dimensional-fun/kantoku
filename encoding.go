package main

import (
	"encoding/json"
	"github.com/ugorji/go/codec"
)

var codecHandle = codec.MsgpackHandle{}

func Encode(contentType string, data interface{}) ([]byte, error) {
	switch contentType {
	case "application/msgpack":
		bytes := new([]byte)
		return *bytes, codec.NewEncoderBytes(bytes, &codecHandle).Encode(data)
	case "application/json":
		return json.Marshal(data)
	default:
		panic("Unsupported Content-Type: " + contentType)
	}
}

func Decode(data []byte, result interface{}) error {
	return codec.NewDecoderBytes(data, &codecHandle).Decode(result)
}
