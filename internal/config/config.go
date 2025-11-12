package config

import (
	"flag"
	"os"
)

type Config struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
	DatabaseDSN     string
}

func InitConfig() *Config {
	var config Config

	flag.StringVar(&config.ServerAddress, "a", ":8080", "address and port to run server")
	flag.StringVar(&config.BaseURL, "b", "http://localhost:8080", "base address for the resulting shortened URL")
	flag.StringVar(&config.FileStoragePath, "f", "", "file where the data is saved in JSON format")
	flag.StringVar(&config.DatabaseDSN, "d", "", "DB connection address")
	flag.Parse()

	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		config.ServerAddress = envAddr
	}

	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		config.BaseURL = envBaseURL
	}

	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" {
		config.DatabaseDSN = envDatabaseDSN
	}

	flagPath, ok := os.LookupEnv("FILE_STORAGE_PATH")
	if ok {
		config.FileStoragePath = flagPath
	}

	return &config
}
