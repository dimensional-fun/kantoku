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
	k := &Kantoku{
		Logger: logrus.New(),
	}

	var err error
	if err = k.loadConfig(); err != nil {
		k.Logger.Fatal("failed to load config: ", err)
	}

	/* configure logging */
	formatter := Formatter{
		TimestampFormat: k.Config.Kantoku.Logging.TimeFormat,
		PrintColors:     true,
	}

	k.Logger.SetFormatter(formatter)

	level, err := logrus.ParseLevel(k.Config.Kantoku.Logging.Level)
	if err != nil {
		k.Logger.Warnln("unable to parse configured log level:", err)
		level = logrus.InfoLevel
	}

	k.Logger.SetLevel(level)

	/* decode public key */
	if k.PublicKey, err = hex.DecodeString(k.Config.Kantoku.PublicKey); err != nil {
		k.Logger.Fatal("failed to decode public key: ", err)
	}

	/* setting up RPC using RabbitMQ */
	k.RpcClient = rpc.NewClient(k.Config.Kantoku.Amqp.URI).
		WithTimeout(3000 * time.Millisecond).
		WithConfirmMode(true).
		WithDebugLogger(k.Logger.Debugf).
		WithErrorLogger(k.Logger.Errorf)

	k.RpcClient.OnStarted(func(_, _ *amqp.Connection, inChan, _ *amqp.Channel) {
		k.Logger.Infoln("connected to rabbitmq")
	})

	/* starting Server */
	k.Logger.Infof("starting w/ version: %s...", version)
	defer k.Logger.Infoln("stopping...")

	handler := mux.NewRouter()

	/* setup router middleware */
	handler.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Powered-By", "catboys")
			handler.ServeHTTP(w, r)
		})
	})

	handler.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lw := NewLogResponseWriter(w)
			start := time.Now()
			handler.ServeHTTP(lw, r)
			stop := time.Now()
			k.Logger.Infof("%s %s %s %d %s", r.RemoteAddr, r.Method, r.URL, lw.StatusCode, stop.Sub(start).String())
		})
	})

	/* expose routes */
	handler.HandleFunc("/v1", k.GetIndex).Methods(http.MethodGet)
	handler.HandleFunc("/v1/info", k.GetInfo).Methods(http.MethodGet)
	handler.HandleFunc("/v1/interactions", k.PostInteractions).Methods(http.MethodPost)

	if k.Config.Kantoku.Server.ExposeTestRoute {
		k.Logger.Warnln("the interaction testing route has been exposed, interactions using any public-key can be published.")
		handler.HandleFunc("/v1/interactions-test", k.PostInteractionsTest).Methods(http.MethodPost)
	}

	addr := fmt.Sprintf("%s:%d", k.Config.Kantoku.Server.Host, k.Config.Kantoku.Server.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	k.Logger.Infoln("listening on", addr)
	if err = server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		k.Logger.Fatal("error while running server: ", err)
	}
}

type Kantoku struct {
	RpcClient *rpc.Client
	Config    Config
	Logger    *logrus.Logger
	PublicKey ed25519.PublicKey
}
