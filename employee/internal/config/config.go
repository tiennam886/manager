package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var config Config

type Config struct {
	LogLevel int
	LogFile  string

	MongoDbUrl string
	Database   string
	Collection string
}

func Get() Config {
	return config
}

func LoadEnvFromFile(configPrefix, envPath string) (err error) {
	if err := godotenv.Load(envPath); err != nil {
		return err
	}
	return envconfig.Process(configPrefix, &config)
}
