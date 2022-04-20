package main

import (
	"encoding/hex"
	"fmt"
	rpc "github.com/0x4b53/amqp-rpc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
	"github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
	"time"
)

var config *toml.Tree
var PublicKey ed25519.PublicKey
var RpcClient *rpc.Client

var InteractionsEvent string

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

	PublicKey = hexDecodedKey

	/* get cool stuff */
	InteractionsEvent = config.GetDefault("kantoku.amqp.event", "INTERACTION_CREATE").(string)
}

func initializeBroker() {
	RpcClient = rpc.NewClient(config.Get("kantoku.amqp.uri").(string)).
		WithTimeout(3000 * time.Millisecond).
		WithConfirmMode(true).
		WithDebugLogger(log.Printf).
		WithErrorLogger(log.Errorf)

	RpcClient.OnStarted(func(_, _ *amqp.Connection, inChan, _ *amqp.Channel) {
		log.Infoln("Connected to AMQP")
	})
}

func initializeServer() {
	log.Infoln("Starting fiber...")

	app := fiber.New(fiber.Config{
		ErrorHandler:          createErrorMessage,
		AppName:               "Kantoku",
		DisableStartupMessage: false,
		Prefork:               config.GetDefault("kantoku.server.prefork", false).(bool),
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Powered-By", "catboys")
		return c.Next()
	})

	app.Use(logger.New(logger.Config{
		Format:     "${black}[${time}]${reset} ${pid} ${cyan}HTTP:${reset} ${magenta}<${method} ${path}>${reset} ${status} ${yellow}${latency}${reset}\n",
		TimeFormat: config.GetDefault("kantoku.logging.time_format", "01-02-06 15:04:0").(string),
		TimeZone:   config.GetDefault("kantoku.logging.timezone", "America/Los_Angeles").(string),
		Output:     os.Stdout,
	}))

	v1 := app.Group("v1")

	v1.Get("/", GetIndex)
	v1.Post("/interactions", PostInteractions)

	if config.GetDefault("kantoku.expose-test-route", false).(bool) {
		log.Warnln("The /v1/interactions-test route has been exposed, this allows any public key to be used.")
		v1.Post("/interactions-test", PostInteractionsTest)
	}

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
		config.GetDefault("kantoku.server.host", "127.0.0.1").(string),
		config.GetDefault("kantoku.server.port", "80").(int64),
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
