package main

import (
	"fmt"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/config"
)

func main() {
	fmt.Println("Iniciando processo de ETL...")
	cfg := config.LoadEnv()
	fmt.Printf("Conectando ao Neo4j em: %s\n", cfg.URI)
}
