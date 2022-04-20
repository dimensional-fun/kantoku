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

func PostInteractions(c *fiber.Ctx) error {
	if c.Get("Content-Type") != "application/json" {
		return c.Status(fiber.StatusBadRequest).JSON(createJson("Invalid Content-Type", false))
	}

	if !VerifyPayload(c, PublicKey) {
		return c.Status(fiber.StatusUnauthorized).JSON(createJson("Invalid Payload", false))
	}

	return handleInteraction(c)
}

func PostInteractionsTest(c *fiber.Ctx) error {
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

	return handleInteraction(c)
}

func handleInteraction(c *fiber.Ctx) error {
	interaction := new(discord.Interaction)
	if err := c.BodyParser(interaction); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if interaction.Type != 1 {
		resp, err := publishInteraction(interaction)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if resp != nil {
			if resp.Headers != nil {
				for key, value := range *resp.Headers {
					c.Set(key, value)
				}
			}

			return c.Status(fiber.StatusOK).Send(resp.Body)
		}
	} else {
		log.Debugln("Received Ping")
	}

	return c.Status(fiber.StatusOK).JSON(discord.InteractionResponse{
		Type: 1,
	})
}

func publishInteraction(i *discord.Interaction) (*KantokuReply, error) {
	contentType := config.GetDefault("kantoku.publish-content-type", "application/msgpack").(string)

	/* encode the interaction so that it can be sent to the message queue */
	body, err := Encode(contentType, i)
	if err != nil {
		return nil, err
	}

	/* publish the interaction and wait for a reply */
	req := rpc.NewRequest().
		WithExchange(config.Get("kantoku.amqp.group").(string)).
		WithRoutingKey(InteractionsEvent)

	req.Publishing.ContentType = contentType
	req.Publishing.Body = body

	res, err := RpcClient.Send(req)
	if err != nil {
		return nil, err
	}

	if err != nil {
		// TODO: handle errors correctly
		switch err {
		case rpc.ErrRequestTimeout:
			return nil, nil

		case rpc.ErrRequestRejected:
			log.Warnln("Interaction rejected?")
			return nil, nil

		case rpc.ErrUnexpectedConnClosed:
			log.Fatalln(err)
		}

		return nil, err
	}

	response := new(KantokuReply)
	if err := Decode(res.Body, &response); err != nil {
		return nil, err
	}

	return response, nil
}
