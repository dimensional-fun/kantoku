package main

import (
	"encoding/hex"

	rpc "github.com/0x4b53/amqp-rpc"
	"github.com/gofiber/fiber/v2"
	"github.com/mixtape-bot/kantoku/discord"
	log "github.com/sirupsen/logrus"
)

type KantokuReply struct {
	Headers *map[string]string `json:"headers"`
	Body    []byte             `json:"body"`
}

func GetIndex(c *fiber.Ctx) error {
	return c.JSON(createJson("Hello, World!", true))
}

func (k *Kontaku) PostInteractions(c *fiber.Ctx) error {
	if c.Get("Content-Type") != "application/json" {
		return c.Status(fiber.StatusBadRequest).JSON(createJson("Invalid Content-Type", false))
	}

	if !VerifyPayload(c, k.PublicKey) {
		return c.Status(fiber.StatusUnauthorized).JSON(createJson("Invalid Payload", false))
	}

	return k.handleInteraction(c)
}

func (k *Kontaku) PostInteractionsTest(c *fiber.Ctx) error {
	if c.Get("Content-Type") != "application/json" {
		return c.Status(fiber.StatusBadRequest).JSON(createJson("Invalid Content-Type", false))
	}

	key := c.Get("X-Kantoku-PublicKey")
	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(createJson("No X-Kantoku-PublicKey given.", false))
	}

	publicKey, err := hex.DecodeString(key)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(createJson(err, false))
	}

	if !VerifyPayload(c, publicKey) {
		return c.Status(fiber.StatusUnauthorized).JSON(createJson("Invalid Payload", false))
	}

	return k.handleInteraction(c)
}

func (k *Kontaku) handleInteraction(c *fiber.Ctx) error {
	var interaction discord.Interaction
	if err := c.BodyParser(&interaction); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if interaction.Type != 1 {
		resp, err := k.publishInteraction(interaction)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if resp.Headers != nil {
			for key, value := range *resp.Headers {
				c.Set(key, value)
			}
		}
		return c.Status(fiber.StatusOK).Send(resp.Body)

	}
	log.Debugln("Received Ping")

	return c.Status(fiber.StatusOK).JSON(discord.InteractionResponse{Type: 1})
}

func (k *Kontaku) publishInteraction(i discord.Interaction) (KantokuReply, error) {
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
		// TODO: handle errors correctly
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
