package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/config"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/loader"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/reader"
	repo "github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/repository/neo4j"
)

func main() {
	fmt.Println("Iniciando processo de ETL...")

	cfg := config.LoadEnv()
	ctx := context.Background()

	// Conexão compartilhada com Neo4j (utilizada também na API)
	driver, err := repo.GetDriver(ctx, cfg.URI, cfg.Username, cfg.Password)
	if err != nil {
		log.Fatalf("Erro de conexão com Neo4j: %v", err)
	}
	defer driver.Close(ctx)

	neoLoader := loader.NewNeo4jLoader(driver)

	// Leitura dos dados
	countries := reader.ReadCountries(cfg.CountriesPath + "/countries.csv")
	vaccines := reader.ReadVaccines(cfg.CountriesPath + "/vaccines.csv")
	cases := reader.ReadCovidCases(cfg.CountriesPath + "/covid_cases.csv")
	vaccs := reader.ReadVaccinations(cfg.CountriesPath + "/vaccinations.csv")
	approvals := reader.ReadVaccineApprovals(cfg.CountriesPath + "/vaccine_approvals.csv")
	relations := reader.ReadCountryVaccines(cfg.CountriesPath + "/country_vaccines.csv")

	fmt.Printf("Países: %d | Vacinas: %d | Casos: %d | Vacinações: %d | Aprovações: %d | Relações: %d\n",
		len(countries), len(vaccines), len(cases), len(vaccs), len(approvals), len(relations))

	// Constraints
	neoLoader.CreateConstraints(ctx)

	// Inserção de nós
	neoLoader.LoadCountries(ctx, countries)
	neoLoader.LoadVaccines(ctx, vaccines)
	neoLoader.LoadCovidCases(ctx, cases)
	neoLoader.LoadVaccinationStats(ctx, vaccs)
	neoLoader.LoadVaccineApprovals(ctx, approvals)

	// Inserção de relacionamentos
	neoLoader.LinkVaccineApprovals(ctx, approvals)
	neoLoader.LinkCountryVaccines(ctx, relations)

	fmt.Println("Carga no Neo4j finalizada com sucesso.")
}
