package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	URI      string
	Username string
	Password string
	CountriesPath string
}

func LoadEnv() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(".env não encontrado")
	}

	return &Config{
		URI:      os.Getenv("NEO4J_URI"),
		Username: os.Getenv("NEO4J_USER"),
		Password: os.Getenv("NEO4J_PASSWORD"),
		CountriesPath: os.Getenv("OUTPUT_DIR"),
	}
}
