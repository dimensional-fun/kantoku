package main

import (
	"encoding/hex"
	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"mixtape.gg/betsu/kantoku/discord"
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
	/* encode the interaction so that it can be sent to the message queue */
	body, err := Encode(i)
	if err != nil {
		return nil, err
	}

	/* publish the interaction and wait for a reply */
	resp, err := Amqp.Call(InteractionsEvent, amqp091.Publishing{
		Body:        body,
		ContentType: "application/msgpack",
	})

	if err != nil {
		switch err {
		case ErrDisconnected:
			err := Amqp.Connect()
			if err != nil {
				log.Fatalln(err)
			}

			return publishInteraction(i)
		case ErrNoRes:
			return nil, nil
		}

		return nil, err
	}

	response := new(KantokuReply)
	if err := Decode(resp.Body, &response); err != nil {
		return nil, err
	}

	return response, nil
}
