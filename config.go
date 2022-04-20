package main

import (
	"os"

	"github.com/pelletier/go-toml"
)

func (k *Kontaku) loadConfig() error {
	file, err := os.Open("kantoku.toml")
	if err != nil {
		return err
	}

	var config Config
	if err = toml.NewDecoder(file).Decode(&config); err != nil {
		return err
	}
	return nil
}

type Config struct {
	Kantoku KantokuConfig `toml:"kantoku"`
}

type KantokuConfig struct {
	PublicKey          string        `toml:"public_key"`
	PublishContentType string        `toml:"publish-content-type"`
	ExposeTestRoute    bool          `toml:"expose-test-route"`
	Server             ServerConfig  `toml:"server"`
	Amqp               AmqpConfig    `toml:"amqp"`
	Logging            LoggingConfig `toml:"logging"`
}

type ServerConfig struct {
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
	Prefork bool   `toml:"prefork"`
}

type AmqpConfig struct {
	URI   string `toml:"uri"`
	Group string `toml:"group"`
	Event string `toml:"event"`
}

type LoggingConfig struct {
	TimeFormat string `toml:"time_format"`
	Timezone   string `toml:"timezone"`
}
