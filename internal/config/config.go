package config

import (
	"fmt"

	"github.com/artem-shestakov/autofaq-webhook/internal/apperror"
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Address string `mapstructure:"address"`
		Port    string `mapstructure:"port"`
	}
	Client struct {
		URL string `mapstructure:"url"`
	}
}

func LoadConfig(path string, errc chan *apperror.Error, warnc chan string) *Config {
	var config Config
	setDefault()
	viper.AutomaticEnv()
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		if notFoundErr, ok := err.(viper.ConfigFileNotFoundError); ok {
			warnc <- fmt.Sprintf("%s. Using environment variables or default settings", notFoundErr)
		} else {
			errc <- apperror.NewError("Read config file error", err.Error(), "0000", err)
		}
		getEnv()
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		errc <- apperror.NewError("Unmarshal config file error", err.Error(), "0000", err)
	}
	return &config
}

func setDefault() {
	viper.SetDefault("server.address", "0.0.0.0")
	viper.SetDefault("server.port", "8000")
	viper.SetDefault("client.url", "http://127.0.0.1:3000")
}

func getEnv() {
	viper.BindEnv("server.address", "SERVER_ADDRESS")
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("client.url", "CLIENT_URL")
}
