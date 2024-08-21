package app

import (
	"encoding/json"
	"os"
)

type (
	Config struct {
		Postgres ConfigPostgres `json:"postgres"`
		NATS     ConfigNATS     `json:"nats"`
		HTTP     ConfigHTTP     `json:"http"`
	}

	ConfigPostgres struct {
		URI string `json:"uri"`
	}

	ConfigNATS struct {
		URL          string `json:"url"`
		ClusterID    string `json:"cluster_id"`
		ClientID     string `json:"client_id"`
		Subscription string `json:"subscription"`
	}

	ConfigHTTP struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}
)

func NewConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	d := json.NewDecoder(f)
	err = d.Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
