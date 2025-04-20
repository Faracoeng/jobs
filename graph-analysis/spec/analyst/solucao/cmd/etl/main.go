package main

import (
	"fmt"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/config"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/reader"
)

func main() {
	fmt.Println("Iniciando processo de ETL...")
	cfg := config.LoadEnv()
	fmt.Printf("Conectando ao Neo4j em: %s\n", cfg.URI)

	// Aqui você pode adicionar a lógica de ETL, como ler arquivos CSV, transformar dados e carregar no Neo4j
	// Exemplo de leitura de um arquivo CSV
	fmt.Println(cfg.CountriesPath)
	countries := reader.ReadeCountries(cfg.CountriesPath + "/countries.csv")
	fmt.Printf("Total de países lidos: %d\n", len(countries))

}
