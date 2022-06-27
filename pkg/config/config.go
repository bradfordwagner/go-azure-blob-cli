package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"

	"github.com/sirupsen/logrus"
)

type Config struct {
	LogLevel           string `env:"LOG_LEVEL,default=info"`
	StorageAccountName string `env:"AZURE_STORAGE_ACCOUNT_NAME,default=blubblubblub"`
	StorageAccountKey  string `env:"AZURE_STORAGE_ACCOUNT_KEY"`
	ContainerName      string `env:"AZURE_STORAGE_CONTAINER_NAME"`
}

// New - gives a new config
func New() Config {
	var c Config
	if err := envconfig.Process(context.Background(), &c); err != nil {
		logrus.WithError(err).Panic()
	}
	return c
}
