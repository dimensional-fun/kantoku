package main

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	rpc "github.com/0x4b53/amqp-rpc"
	log "github.com/sirupsen/logrus"
)

type KantokuReply struct {
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

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

	if !k.VerifyRequest(r, k.PublicKey) {
		w.WriteHeader(http.StatusUnauthorized)
		k.createJsonResponse(w, "Invalid Payload", false)
		return
	}

	k.handleInteraction(w, r)
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

	if !k.VerifyRequest(r, publicKey) {
		w.WriteHeader(http.StatusUnauthorized)
		k.createJsonResponse(w, "Invalid Payload", false)
		return
	}

	k.handleInteraction(w, r)
}

func (k *Kantoku) handleInteraction(w http.ResponseWriter, r *http.Request) {
	var interaction map[string]any
	if err := json.NewDecoder(r.Body).Decode(&interaction); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		k.createJsonResponse(w, err.Error(), false)
		return
	}

	if interaction["type"] == 1 {
		log.Debugln("Received Ping")
		k.createJson(w, InteractionResponse{Type: 1})
		return
	}

	resp, err := k.publishInteraction(interaction)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		k.createJsonResponse(w, err.Error(), false)
		return
	}

	for key, value := range resp.Headers {
		w.Header().Set(key, value)
	}

	if _, err = w.Write(resp.Body); err != nil {
		k.Logger.Error("Error writing response body: ", err.Error())
	}
}

func (k *Kantoku) publishInteraction(i map[string]any) (KantokuReply, error) {
	contentType := k.Config.Kantoku.PublishContentType

	/* encode the interaction so that it can be sent to the message queue */
	body, err := Encode(contentType, i)
	if err != nil {
		return KantokuReply{}, err
	}

	/* publish the interaction and wait for a reply */
	req := rpc.NewRequest().
		WithExchange(k.Config.Kantoku.Amqp.Group).
		WithRoutingKey(k.Config.Kantoku.Amqp.Event)

	req.Publishing.ContentType = contentType
	req.Publishing.Body = body

	res, err := k.RpcClient.Send(req)
	if err != nil {
		return KantokuReply{}, err
	}

	if err != nil {
		// TODO: handle rpc errors correctly
		switch err {
		case rpc.ErrRequestTimeout:
			return KantokuReply{}, nil

		case rpc.ErrRequestRejected:
			log.Warnln("Interaction rejected?")
			return KantokuReply{}, nil

		case rpc.ErrUnexpectedConnClosed:
			log.Fatalln(err)
		}

		return KantokuReply{}, err
	}

	var response KantokuReply
	return response, Decode(res.Body, &response)
}

type InteractionResponse struct {
	Type int `json:"type"`
}
