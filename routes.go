package main

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
)

func (k *Kantoku) GetIndex(w http.ResponseWriter, _ *http.Request) {
	k.createJsonResponse(w, "Hello, World!", true)
}

func (k *Kantoku) GetInfo(w http.ResponseWriter, _ *http.Request) {
	k.createJsonResponse(w, version, true)
}

func (k *Kantoku) PostInteractions(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		k.createJsonResponse(w, "Invalid Content-Type", false)
		return
	}

	k.handleInteraction(w, r, k.PublicKey)
}

func (k *Kantoku) PostInteractionsTest(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		k.createJsonResponse(w, "Invalid Content-Type", false)
		return
	}

	key := r.Header.Get("X-Kantoku-PublicKey")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		k.createJsonResponse(w, "No X-Kantoku-PublicKey given", false)
		return
	}

	publicKey, err := hex.DecodeString(key)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		k.createJsonResponse(w, err.Error(), false)
		return
	}

	k.handleInteraction(w, r, publicKey)
}

func (k *Kantoku) handleInteraction(w http.ResponseWriter, r *http.Request, pk ed25519.PublicKey) {
	if !k.VerifyRequest(r, pk) {
		w.WriteHeader(http.StatusUnauthorized)
		k.createJsonResponse(w, "Invalid Payload", false)
		return
	}

	var interaction struct {
		Type int `json:"type"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		k.createJsonResponse(w, err.Error(), false)
		return
	}

	if err := json.Unmarshal(body, &interaction); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		k.createJsonResponse(w, err.Error(), false)
		return
	}

	if interaction.Type == 1 {
		k.Logger.Debugln("Received Ping")
		k.createJson(w, map[string]any{"type": 1})
		return
	}

	resp, err := k.publishInteraction(body)
	if err != nil {
		if err == nats.ErrNoResponders && k.NoResponders != nil {
			k.reply(w, k.NoResponders, "application/json")
			return
		}

		k.Logger.Errorln("Error publishing interaction:", err)
		w.WriteHeader(http.StatusBadRequest)
		k.createJsonResponse(w, err.Error(), false)
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		return
	}

	k.reply(w, resp.Data, contentType)
}

func (k *Kantoku) reply(w http.ResponseWriter, body []byte, contentType string) {
	w.Header().Set("Content-Type", contentType)
	if _, err := w.Write(body); err != nil {
		k.Logger.Errorln("Error writing response body:", err.Error())
	}
}

func (k *Kantoku) publishInteraction(body []byte) (*nats.Msg, error) {
	msg := nats.NewMsg(k.Config.Kantoku.Nats.Subject)
	msg.Data = body
	msg.Header.Add("Content-Type", "application/json")

	return k.NatsConn.RequestMsg(msg, 3*time.Second)
}
