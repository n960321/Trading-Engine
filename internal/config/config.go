package config

import (
	"Trading-Engine/internal/server"
	"Trading-Engine/internal/storage/mysql"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	Http  server.Config `mapstructure:"http"`
	Mysql mysql.Config  `mapstructure:"mysql"`
}

func GetConfig(configFile string) *Config {
	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()

	if err != nil {
		log.Panic().Msgf("load config failed, err: %v", err)
	}

	var c Config
	err = viper.Unmarshal(&c)

	if err != nil {
		log.Panic().Msgf("unmarshal config failed, err: %v", err)
	}

	return &c
}
