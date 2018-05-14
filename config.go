package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/CactusDev/Xerophi/rethink"
)

// Config keeps track of the config set in config.json
type Config struct {
	Rethink rethinkCfg `json:"rethink"`
	Sentry  sentryCfg  `json:"sentry"`
	Server  serverCfg  `json:"server"`
	Redis   redisCfg   `json:"redis"`
}

type rethinkCfg struct {
	Connection rethink.ConnectionOpts `json:"connection"`
	DB         string                 `json:"db"`
}

type sentryCfg struct {
	DSN     string `json:"dsn"`
	Enabled bool   `json:"enabled"`
}

type redisCfg struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type serverCfg struct {
	Port int `json:"port"`
}

// LoadConfig tries to load the config from the default path "./config.json"
// By default the config for Sepal will be in the same directory as
// the executable
func LoadConfig() Config {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := path.Dir(ex)
	return LoadConfigFromPath(fmt.Sprintf("%s/config.json", exPath))
}

// LoadConfigFromPath loads config from a specific path
func LoadConfigFromPath(path string) Config {
	config := Config{}
	configFile, err := os.Open(path)
	defer configFile.Close()
	if err != nil {
		panic(err)
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		panic(err)
	}

	return config
}
