package main

import (
	"bytes"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
)

func (k *Kontaku) VerifyRequest(r *http.Request, publicKey ed25519.PublicKey) bool {
	var msg bytes.Buffer

	signature := r.Header.Get("X-Signature-Ed25519")
	if signature == "" {
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	if len(sig) != ed25519.SignatureSize || sig[63]&224 != 0 {
		return false
	}

	timestamp := r.Header.Get("X-Signature-Timestamp")
	if timestamp == "" {
		return false
	}

	msg.WriteString(timestamp)

	defer func() {
		err = r.Body.Close()
		if err != nil {
			k.Logger.Error("error while closing request body: ", err)
		}
	}()
	var body bytes.Buffer

	defer func() {
		r.Body = io.NopCloser(&body)
	}()

	_, err = io.Copy(&msg, io.TeeReader(r.Body, &body))
	if err != nil {
		return false
	}

	return ed25519.Verify(publicKey, msg.Bytes(), sig)
}
