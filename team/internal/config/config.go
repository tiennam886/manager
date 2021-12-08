package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var config Config

type Config struct {
	LogLevel        int
	MongoDbUrl      string
	Database        string
	Collection      string
	RedisUrl        string
	MysqlUrl        string
	TeamTable       string
	TeamMemberTable string
	PostgresqlUrl   string
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
