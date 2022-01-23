package main

import (
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/pelletier/go-toml"
	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

var config *toml.Tree
var publicKey ed25519.PublicKey
var Amqp *AMQP

func loadConfig() {
	/* get config */
	t, err := toml.LoadFile("kantoku.toml")
	if err != nil {
		log.Fatalln(err)
	}

	config = t

	/* get public key */
	hexDecodedKey, err := hex.DecodeString(config.Get("kantoku.public_key").(string))
	if err != nil {
		log.Fatalf("error while decoding public key: %s", err)
	}

	publicKey = hexDecodedKey
}

func initializeBroker() {
	Amqp = &AMQP{
		Group:   config.Get("kantoku.amqp.exchange").(string),
		Timeout: time.Duration(time.Duration(15).Minutes()),
	}

	conn, err := amqp091.Dial(config.Get("kantoku.amqp.uri").(string))
	if err != nil {
		log.Fatalln(err)
	}

	err = Amqp.Init(conn)
	if err != nil {
		return
	}

	log.Infoln("Connected to AMQP")
}

func initializeServer() {
	log.Infoln("Starting fiber...")

	app := fiber.New(fiber.Config{
		ErrorHandler:          createErrorMessage,
		DisableStartupMessage: true,
	})

	app.Use(logger.New(logger.Config{
		Format:     "${black}[${time}] ${cyan}HTTP:${reset} ${magenta}<${method} ${path}>${reset} ${status} ${yellow}${latency}${reset}\n",
		TimeFormat: config.Get("kantoku.logging.time_format").(string),
		TimeZone:   config.Get("kantoku.logging.timezone").(string),
		Output:     os.Stdout,
	}))

	v1 := app.Group("v1")

	v1.Get("/", GetIndex)
	v1.Post("/interactions", PostInteractions)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(createJson(
			fiber.Map{
				"message": fmt.Sprintf("I was unable to find %s %s", c.Method(), c.Path()),
			},
			true,
		))
	})

	addr := fmt.Sprintf(
		"%s:%d",
		config.Get("kantoku.server.host").(string),
		config.Get("kantoku.server.port").(int64),
	)

	log.Infof("Listening on %s", addr)

	err := app.Listen(addr)
	if err != nil {
		log.Fatalln(err)
		return
	}
}

func createJson(data interface{}, success bool) fiber.Map {
	return fiber.Map{
		"data":    data,
		"success": success,
	}
}

func createErrorMessage(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "A server side error has occurred."

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	} else if err.Error() != "" {
		message = err.Error()
	}

	log.Println("Error: ", message)

	return ctx.Status(code).JSON(createJson(
		fiber.Map{
			"code":    code,
			"message": message,
		},
		false,
	))
}

func verifyDiscordPayload(c *fiber.Ctx) bool {
	signature := c.Get("X-Signature-Ed25519")
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

	timestamp := c.Get("X-Signature-Timestamp")
	if timestamp == "" {
		return false
	}

	var msg bytes.Buffer
	msg.WriteString(timestamp)
	msg.Write(c.Body())

	return ed25519.Verify(publicKey, msg.Bytes(), sig)
}
