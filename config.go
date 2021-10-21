package manager

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type TomlConfig struct {
	ServerHost       string
	ServerPort       string
	CachePort        string
	MongoPort        string
	MongoDatabase    string
	MongoEmployeeCol string
	MongoTeamsCol    string
	MySqlPort        string
	MySqlUser        string
	MySqlHost        string
	MySqlDatabase    string
	MySqlEmployee    string
	MySqlTeams       string
	MySqlTeamMem     string
}

var conf TomlConfig

func loadConfig() (TomlConfig, error) {
	in, err := os.Open("/home/namtt/go/src/manager/config.toml")
	if err != nil {
		return conf, err
	}

	_, err = toml.DecodeReader(in, &conf)
	fmt.Println(conf.MySqlHost)
	return conf, err
}
