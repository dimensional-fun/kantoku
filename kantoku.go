package main

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	rpc "github.com/0x4b53/amqp-rpc"
	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func main() {
	logger := logrus.New()
	k := &Kantoku{
		Logger: logger,
	}

	var err error
	if err = k.loadConfig(); err != nil {
		k.Logger.Fatal("Failed to load config: ", err)
	}
	if os.Getenv("FIBER_PREFORK_CHILD") == "1" {
		k.Logger.SetOutput(NopWriter{})
	} else {
		k.Logger.SetReportCaller(true)
		k.Logger.SetFormatter(Formatter{TimestampFormat: k.Config.Kantoku.Logging.TimeFormat})
	}

	if k.PublicKey, err = hex.DecodeString(k.Config.Kantoku.PublicKey); err != nil {
		k.Logger.Fatal("Failed to decode public key: ", err)
	}

	k.RpcClient = rpc.NewClient(k.Config.Kantoku.Amqp.URI).
		WithTimeout(3000 * time.Millisecond).
		WithConfirmMode(true).
		WithDebugLogger(k.Logger.Printf).
		WithErrorLogger(k.Logger.Errorf)

	k.RpcClient.OnStarted(func(_, _ *amqp.Connection, inChan, _ *amqp.Channel) {
		k.Logger.Infoln("Connected to AMQP")
	})

	// Starting Server
	k.Logger.Infoln("Starting server...")
	defer k.Logger.Infoln("Stopping server...")

	mux := http.NewServeMux()
	mux.HandleFunc("/", k.GetIndex)
	mux.HandleFunc("/interactions", k.PostInteractions)

	if k.Config.Kantoku.ExposeTestRoute {
		k.Logger.Warnln("The /v1/interactions-test route has been exposed, this allows any public key to be used.")
		mux.HandleFunc("/interactions-test", k.PostInteractionsTest)
	}

	addr := fmt.Sprintf("%s:%d", k.Config.Kantoku.Server.Host, k.Config.Kantoku.Server.Port)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	k.Logger.Infof("Listening on: %s", addr)

	if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		k.Logger.Fatal("Error while running server: ", err)
	}
}

type Kantoku struct {
	RpcClient *rpc.Client
	Config    Config
	Logger    *logrus.Logger
	PublicKey ed25519.PublicKey
}
