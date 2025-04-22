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
	APIHost  string
	APIPort  string
	ETL_INTERVAL_SECONDS string
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
		APIHost: os.Getenv("API_HOST"),
		APIPort: os.Getenv("API_PORT"),
		ETL_INTERVAL_SECONDS: os.Getenv("ETL_INTERVAL_SECONDS"),
	}
}
