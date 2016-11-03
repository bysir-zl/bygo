package config

import (
	"encoding/json"
	"github.com/bysir-zl/bygo/db"
	"io/ioutil"
	"os"
)

type Config struct {
	Evn        string
	Debug      bool
	CacheDrive string
	DbConfigs  map[string]db.DbConfig
	RedisHost  string
}

type configFile struct {
	Evn      string
	App      struct {
			 Debug bool
		 }
	Cache    struct {
			 Driver string `json:"driver"`
		 }
	Database map[string]db.DbConfig
	Redis    struct {
			 Host string
		 }
}

func LoadConfigFromFile(filePath string) (config Config, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	configFile := configFile{}
	json.Unmarshal(bs, &configFile)

	config = Config{}
	config.Evn = configFile.Evn

	config.Debug = configFile.App.Debug
	config.CacheDrive = configFile.Cache.Driver
	config.RedisHost = configFile.Redis.Host

	config.DbConfigs = configFile.Database

	return
}
