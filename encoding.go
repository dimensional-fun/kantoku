package main

import (
	"encoding/json"

	"github.com/ugorji/go/codec"
)

var codecHandle = &codec.MsgpackHandle{}

func Encode(contentType string, data interface{}) ([]byte, error) {
	/* due to certain libraries using JSON-specific code (cough... kord... cough) we need to support multiple encodings */
	switch contentType {
	case "application/msgpack":
		var bytes []byte
		return bytes, codec.NewEncoderBytes(&bytes, codecHandle).Encode(data)
	case "application/json":
		return json.Marshal(data)
	default:
		panic("Unsupported Content-Type: " + contentType)
	}
}

func Decode(data []byte, result interface{}) error {
	/* only msgpack encoded responses are allowed */
	return codec.NewDecoderBytes(data, codecHandle).Decode(result)
}
