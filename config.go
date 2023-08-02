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
	PublicKey string        `toml:"public_key"`
	Server    ServerConfig  `toml:"server"`
	Nats      NatsConfig    `toml:"nats"`
	Logging   LoggingConfig `toml:"logging"`
}

type ServerConfig struct {
	Host            string `toml:"host"`
	Port            int64  `toml:"port"`
	ExposeTestRoute bool   `toml:"expose_test_route"`
}

type NatsConfig struct {
	URI          string       `toml:"uri"`
	Subject      string       `toml:"group"`
	NoResponders *interface{} `toml:"no_responders"`
}

type LoggingConfig struct {
	TimeFormat string `toml:"time_format"`
	Level      string `toml:"level"`
}
