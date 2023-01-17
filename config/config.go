package config

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	RpcAddr     string
	ServerPort  string
	AccountSeed string
	AccountID   string
}

var DefaultConfigPath = "./config/config.toml"

var config Config

func InitConfig(path string) error {
	if path == "" {
		path = DefaultConfigPath
	}
	if _, err := os.Stat(path); err != nil {
		return errors.Wrap(err, "config file not exist")
	}
	viper.SetConfigFile(path)
	viper.SetConfigType("toml")
	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "can not load config file")
	}
	if err := viper.Unmarshal(&config); err != nil {
		return errors.Wrap(err, "unmarshal config file error")
	}
	return nil
}

func GetConfig() Config {
	return config
}
