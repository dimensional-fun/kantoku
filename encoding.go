package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ugorji/go/codec"
)

var codecHandle = codec.MsgpackHandle{}

func Encode(contentType string, data interface{}) ([]byte, error) {
	if contentType == "application/msgpack" {
		bytes := new([]byte)
		return *bytes, codec.NewEncoderBytes(bytes, &codecHandle).Encode(data)
	} else if contentType == "application/json" {
		return json.Marshal(data)
	}

	return nil, errors.New("Unsupported Content-Type")
}

func Decode(contentType string, data []byte, result interface{}) error {
	switch contentType {
	case "application/msgpack":
		return codec.NewDecoderBytes(data, &codecHandle).Decode(result)
	case "application/json":
		return json.Unmarshal(data, result)
	default:
		panic(fmt.Sprintf("Unknown content-type: %s", contentType))
	}
}
