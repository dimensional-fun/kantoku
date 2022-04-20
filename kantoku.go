package main

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	rpc "github.com/0x4b53/amqp-rpc"
	"github.com/gorilla/mux"
	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const version = "dev"

func main() {
	logger := logrus.New()
	k := &Kantoku{
		Logger: logger,
	}

	var err error
	if err = k.loadConfig(); err != nil {
		k.Logger.Fatal("Failed to load config: ", err)
	}

	k.Logger.SetReportCaller(true)
	k.Logger.SetFormatter(Formatter{TimestampFormat: k.Config.Kantoku.Logging.TimeFormat})

	if k.PublicKey, err = hex.DecodeString(k.Config.Kantoku.PublicKey); err != nil {
		k.Logger.Fatal("Failed to decode public key: ", err)
	}

	/* setting up RPC using RabbitMQ */
	k.RpcClient = rpc.NewClient(k.Config.Kantoku.Amqp.URI).
		WithTimeout(3000 * time.Millisecond).
		WithConfirmMode(true).
		WithDebugLogger(k.Logger.Printf).
		WithErrorLogger(k.Logger.Errorf)

	k.RpcClient.OnStarted(func(_, _ *amqp.Connection, inChan, _ *amqp.Channel) {
		k.Logger.Infoln("Connected to AMQP")
	})

	/* starting Server */
	k.Logger.Infof("Starting Kantoku %s...", version)
	defer k.Logger.Infoln("Stopping Kantoku...")

	handler := mux.NewRouter()
	handler.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Powered-By", "catboys")
			handler.ServeHTTP(w, r)
		})
	})
	handler.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			k.Logger.Infof("%s %s", r.Method, r.URL.Path)
			handler.ServeHTTP(w, r)
		})
	})

	handler.HandleFunc("/v1", k.GetIndex).Methods(http.MethodGet)
	handler.HandleFunc("/v1/info", k.GetInfo).Methods(http.MethodGet)
	handler.HandleFunc("/v1/interactions", k.PostInteractions).Methods(http.MethodPost)

	if k.Config.Kantoku.Server.ExposeTestRoute {
		k.Logger.Warnln("The interaction testing route has been exposed, interactions using any public-key can be published.")
		handler.HandleFunc("/v1/interactions-test", k.PostInteractionsTest)
	}

	addr := fmt.Sprintf("%s:%d", k.Config.Kantoku.Server.Host, k.Config.Kantoku.Server.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	k.Logger.Infoln("Listening on", addr)
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
