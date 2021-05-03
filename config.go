package main

import (
	"errors"
	"io"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	BindAddr            string        `yaml:"bind_addr"`
	PriceUpdateInterval time.Duration `yaml:"price_update_interval"`
	MongoDB             MongoDBConfig `yaml:"mongodb"`
}

type MongoDBConfig struct {
	URI            string `yaml:"uri"`
	DB             string `yaml:"db"`
	CoinCollection string `yaml:"coin_collection"`
}

var DefaultServerConfig = ServerConfig{
	BindAddr:            "0.0.0.0:8000",
	PriceUpdateInterval: 5 * time.Second,
	MongoDB: MongoDBConfig{
		URI:            "mongodb://localhost",
		DB:             "random_price_oracle",
		CoinCollection: "coins",
	},
}

func ReadServerConfig(name string) (ServerConfig, error) {
	f, err := os.Open(name)
	if err != nil {
		return ServerConfig{}, err
	}
	defer f.Close()
	cfg := DefaultServerConfig
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil && !errors.Is(err, io.EOF) {
		return ServerConfig{}, err
	}
	return cfg, nil
}
