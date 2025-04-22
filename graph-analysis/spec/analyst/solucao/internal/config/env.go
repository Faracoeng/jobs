package config

import (
	"log"
	"os"
)

type Config struct {
	URI      string
	Username string
	Password string

	CountriesPath       string
	APIHost             string
	APIPort             string
	ETL_INTERVAL_SECONDS string
}

func LoadEnv() *Config {
	log.Println("[CONFIG] Carregando variáveis de ambiente...")

	cfg := &Config{
		URI:      os.Getenv("NEO4J_URI"),
		Username: os.Getenv("NEO4J_USERNAME"),
		Password: os.Getenv("NEO4J_PASSWORD"),

		CountriesPath:       os.Getenv("OUTPUT_DIR"),
		APIHost:             os.Getenv("API_HOST"),
		APIPort:             os.Getenv("API_PORT"),
		ETL_INTERVAL_SECONDS: os.Getenv("ETL_INTERVAL_SECONDS"),
	}

	log.Printf("[CONFIG] URI: %s", cfg.URI)
	log.Printf("[CONFIG] Username: %s", cfg.Username)
	log.Printf("[CONFIG] Path dos CSVs: %s", cfg.CountriesPath)
	log.Printf("[CONFIG] Host da API: %s", cfg.APIHost)
	log.Printf("[CONFIG] Porta da API: %s", cfg.APIPort)
	log.Printf("[CONFIG] Intervalo ETL: %s segundos", cfg.ETL_INTERVAL_SECONDS)

	return cfg
}
