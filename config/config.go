package config

import (
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func Load() *Config {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if envFileExists(".env") {
		viper.SetConfigFile(".env")
		viper.ReadInConfig()
	}

	return &Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}
}

func envFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
