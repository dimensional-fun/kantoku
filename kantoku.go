package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
	"github.com/sirupsen/logrus"
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

	/* prepare no_responders reply. */
	if k.Config.Kantoku.Nats.NoResponders != nil {
		b, err := json.Marshal(map[string]any{
			"type": 4,
			"data": k.Config.Kantoku.Nats.NoResponders,
		})

		if err != nil {
			k.Logger.Warnln("unable to encode 'no_responders' reply: ", err)
		} else {
			k.NoResponders = b
		}
	}

	/* prepare nats client */
	nc, err := nats.Connect(strings.Join(k.Config.Kantoku.Nats.Servers, ", "))
	if err != nil {
		k.Logger.Fatal("connecting to NATS server failed: ", err)
	}

	k.NatsConn = nc
	k.Logger.Infoln("connected to NATS server!")

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
	NatsConn     *nats.Conn
	NoResponders []byte
	Config       Config
	Logger       *logrus.Logger
	PublicKey    ed25519.PublicKey
}
