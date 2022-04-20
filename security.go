package main

import (
	"bytes"
	"encoding/hex"
	"github.com/gofiber/fiber/v2"
	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
)

func VerifyPayload(c *fiber.Ctx, key ed25519.PublicKey) bool {
	/* check for and decode the given signature. */
	signature := c.Get("X-Signature-Ed25519")
	if signature == "" {
		return false
	}

	sig, err := hex.DecodeString(signature)
	if err != nil {
		return false
	}

	/* check if the signature is the correct size (64 bytes) */
	if len(sig) != ed25519.SignatureSize || sig[63]&224 != 0 {
		return false
	}

	/* check if there is a given timestamp  */
	timestamp := c.Get("X-Signature-Timestamp")
	if timestamp == "" {
		return false
	}

	var msg bytes.Buffer
	msg.WriteString(timestamp)
	msg.Write(c.Body())

	return ed25519.Verify(key, msg.Bytes(), sig)
}
