package main

import (
	"fmt"
	"context"
	"log"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/config"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/reader"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/loader"
)

func main() {
	fmt.Println("Iniciando processo de ETL...")

	cfg := config.LoadEnv()
	
	ctx := context.Background()

	driver, err := loader.NewNeo4jDriver(ctx, cfg)
		if err != nil {
			log.Fatalf("Erro de conexão com Neo4j: %v", err)
	}

	defer driver.Close(ctx)


	fmt.Printf("Conectado ao Neo4j em: %s\n", cfg.URI)

	// Países
	countries := reader.ReadCountries(cfg.CountriesPath + "/countries.csv")
	fmt.Printf(" Países carregados: %d\n", len(countries))

	// Vacinas
	vaccines := reader.ReadVaccines(cfg.CountriesPath + "/vaccines.csv")
	fmt.Printf(" Vacinas carregadas: %d\n", len(vaccines))

	// Casos de Covid
	cases := reader.ReadCovidCases(cfg.CountriesPath + "/covid_cases.csv")
	fmt.Printf(" Casos de Covid carregados: %d\n", len(cases))

	// Vacinações
	vaccs := reader.ReadVaccinations(cfg.CountriesPath + "/vaccinations.csv")
	fmt.Printf(" Estatísticas de vacinação carregadas: %d\n", len(vaccs))

	// Aprovações
	approvals := reader.ReadVaccineApprovals(cfg.CountriesPath + "/vaccine_approvals.csv")
	fmt.Printf(" Aprovações de vacinas carregadas: %d\n", len(approvals))

	// Relacionamentos País-Vacina
	relations := reader.ReadCountryVaccines(cfg.CountriesPath + "/country_vaccines.csv")
	fmt.Printf(" Relações país-vacina carregadas: %d\n", len(relations))


	loader := loader.NewNeo4jLoader(driver)
	fmt.Println("Carregando dados no Neo4j...")
	loader.LoadCountries(ctx, countries)

	fmt.Println("Carregando dados de vacinas no Neo4j...")
	loader.LoadVaccines(ctx, vaccines)
	fmt.Println("Carregando dados de aprovações de vacinas no Neo4j...")
	loader.LoadVaccineApprovals(ctx, approvals)

	fmt.Println("Carregando dados de casos de Covid no Neo4j...")
	loader.LoadCovidCases(ctx, cases)

}
