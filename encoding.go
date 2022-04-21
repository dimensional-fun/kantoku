package main

import (
	"encoding/json"

	"github.com/ugorji/go/codec"
)

var codecHandle = &codec.MsgpackHandle{}

func Encode(contentType string, data any) (bytes []byte, err error) {
	/* due to certain libraries using JSON-specific code (cough... kord... cough) we need to support multiple encodings */
	switch contentType {
	case "application/msgpack":
		err = codec.NewEncoderBytes(&bytes, codecHandle).Encode(data)
	case "application/json":
		bytes, err = json.Marshal(data)
	default:
		panic("Unsupported Content-Type: " + contentType)
	}
	return
}

func Decode(data []byte, result any) error {
	/* only msgpack encoded responses are allowed */
	return codec.NewDecoderBytes(data, codecHandle).Decode(result)
}
