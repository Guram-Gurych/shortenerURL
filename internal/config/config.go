package config

import "flag"

type Config struct {
	ServerAddress string
	BaseURL       string
}

func InitConfig() *Config {
	var config Config

	flag.StringVar(&config.ServerAddress, "a", ":8080", "address and port to run server")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "base address for the resulting shortened URL")
	flag.Parse()

	return &config
}
