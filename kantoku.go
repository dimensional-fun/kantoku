package main

import (
	"encoding/hex"
	"os"

	rpc "github.com/0x4b53/amqp-rpc"
	"github.com/oasisprotocol/curve25519-voi/primitives/ed25519"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	k := &Kontaku{
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

	k.initializeBroker()
	k.initializeServer()
}

type Kontaku struct {
	RpcClient *rpc.Client
	Config    Config
	Logger    *logrus.Logger
	PublicKey ed25519.PublicKey
}
