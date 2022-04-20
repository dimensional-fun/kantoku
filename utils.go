package main

import (
	"fmt"
	"os"
	"time"

	rpc "github.com/0x4b53/amqp-rpc"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var InteractionsEvent string

func (k *Kontaku) initializeBroker() {
	k.RpcClient = rpc.NewClient(k.Config.Kantoku.Amqp.URI).
		WithTimeout(3000 * time.Millisecond).
		WithConfirmMode(true).
		WithDebugLogger(k.Logger.Printf).
		WithErrorLogger(k.Logger.Errorf)

	k.RpcClient.OnStarted(func(_, _ *amqp.Connection, inChan, _ *amqp.Channel) {
		log.Infoln("Connected to AMQP")
	})
}

func (k *Kontaku) initializeServer() {
	log.Infoln("Starting fiber...")

	app := fiber.New(fiber.Config{
		ErrorHandler:          createErrorMessage,
		AppName:               "Kantoku",
		DisableStartupMessage: false,
		Prefork:               k.Config.Kantoku.Server.Prefork,
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-Powered-By", "catboys")
		return c.Next()
	})

	app.Use(logger.New(logger.Config{
		Format:     "${black}[${time}]${reset} ${pid} ${cyan}HTTP:${reset} ${magenta}<${method} ${path}>${reset} ${status} ${yellow}${latency}${reset}\n",
		TimeFormat: k.Config.Kantoku.Logging.TimeFormat,
		TimeZone:   k.Config.Kantoku.Logging.Timezone,
		Output:     os.Stdout,
	}))

	v1 := app.Group("v1")

	v1.Get("/", GetIndex)
	v1.Post("/interactions", k.PostInteractions)

	if k.Config.Kantoku.ExposeTestRoute {
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

	addr := fmt.Sprintf("%s:%d", k.Config.Kantoku.Server.Host, k.Config.Kantoku.Server.Port)

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
