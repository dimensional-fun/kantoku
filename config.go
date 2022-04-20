package main

import (
	"os"

	"github.com/pelletier/go-toml"
)

func (k *Kantoku) loadConfig() error {
	file, err := os.Open("kantoku.toml")
	if err != nil {
		return err
	}

	return toml.NewDecoder(file).Decode(&k.Config)
}

type Config struct {
	Kantoku KantokuConfig `toml:"kantoku"`
}

type KantokuConfig struct {
	PublicKey          string        `toml:"public_key"`
	PublishContentType string        `toml:"publish_content_type"`
	Server             ServerConfig  `toml:"server"`
	Amqp               AmqpConfig    `toml:"amqp"`
	Logging            LoggingConfig `toml:"logging"`
}

type ServerConfig struct {
	Host            string `toml:"host"`
	Port            int64  `toml:"port"`
	ExposeTestRoute bool   `toml:"expose_test_route"`
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
