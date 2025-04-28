package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	URI      string
	Username string
	Password string

	CountriesPath       string
	APIHost             string
	APIPort             string
	ETLIntervalSeconds  int
	ETLBatchSize        int
}

func LoadEnv() *Config {
	log.Println("[CONFIG] Carregando variáveis de ambiente...")

	intervalStr := os.Getenv("ETL_INTERVAL_SECONDS")
	interval, err := strconv.Atoi(intervalStr)
	if err != nil || interval <= 0 {
		log.Println("Intervalo inválido ou ausente. Usando 60 segundos como padrão.")
		interval = 60
	}

	batchStr := os.Getenv("ETL_BATCH_SIZE")
	batchSize, err := strconv.Atoi(batchStr)
	if err != nil || batchSize <= 0 {
		log.Println("Batch size inválido ou ausente. Usando 500 como padrão.")
		batchSize = 500
	}

	cfg := &Config{
		URI:      os.Getenv("NEO4J_URI"),
		Username: os.Getenv("NEO4J_USERNAME"),
		Password: os.Getenv("NEO4J_PASSWORD"),

		CountriesPath:      os.Getenv("OUTPUT_DIR"),
		APIHost:            os.Getenv("API_HOST"),
		APIPort:            os.Getenv("API_PORT"),
		ETLIntervalSeconds: interval,
		ETLBatchSize:       batchSize,
	}

	log.Printf("[CONFIG] URI: %s", cfg.URI)
	log.Printf("[CONFIG] Username: %s", cfg.Username)
	log.Printf("[CONFIG] Path dos CSVs: %s", cfg.CountriesPath)
	log.Printf("[CONFIG] Host da API: %s", cfg.APIHost)
	log.Printf("[CONFIG] Porta da API: %s", cfg.APIPort)
	log.Printf("[CONFIG] Intervalo ETL: %d segundos", cfg.ETLIntervalSeconds)
	log.Printf("[CONFIG] Batch Size: %d registros", cfg.ETLBatchSize)

	return cfg
}
