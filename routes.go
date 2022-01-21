package main

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"mixtape.gg/betsu/kantoku/discord"
)

func GetIndex(c *fiber.Ctx) error {
	return c.JSON(createJson(
		fiber.Map{
			"message": "hello world",
		},
		true,
	))
}

func PostInteractions(c *fiber.Ctx) error {
	if c.Get("Content-Type") != "application/json" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if !verifyDiscordPayload(c) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

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
			return c.Status(fiber.StatusOK).JSON(resp)
		}

		return c.SendStatus(fiber.StatusNoContent)
	}

	log.Debugln("Received Ping")
	return c.Status(fiber.StatusOK).JSON(discord.InteractionResponse{
		Type: 1,
	})
}

func publishInteraction(i *discord.Interaction) (*discord.InteractionResponse, error) {
	body, err := Encode(i)
	if err != nil {
		return nil, err
	}

	resp, err := Amqp.Call("INTERACTION_CREATE", amqp091.Publishing{
		Body: body,
	})

	if err != nil {
		switch err {
		case ErrDisconnected:
			initializeBroker()
			return publishInteraction(i)
		case ErrNoRes:
			return nil, nil
		}

		return nil, err
	}

	response := new(discord.InteractionResponse)
	if err := Decode(resp, &response); err != nil {
		return nil, err
	}

	j, _ := json.Marshal(response)

	log.Infoln(string(j))

	return response, nil
}
