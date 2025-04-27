package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/config"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/loader"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/reader"
	repo "github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/infra/db/neo4j"
)

func main() {
	cfg := config.LoadEnv()
	ctx := context.Background()

	intervalStr := os.Getenv("ETL_INTERVAL_SECONDS")
	intervalSec, err := strconv.Atoi(intervalStr)
	if err != nil || intervalSec <= 0 {
		log.Println("Intervalo inválido ou ausente. Usando 60 segundos como padrão.")
		intervalSec = 60
	}

	driver, err := repo.GetDriver(ctx, cfg.URI, cfg.Username, cfg.Password)
	if err != nil {
		log.Fatalf("Erro ao conectar ao Neo4j: %v", err)
	}
	defer driver.Close(ctx)

	neoLoader := loader.NewNeo4jLoader(driver)

	for {
		log.Println("Iniciando novo ciclo de ETL.")
		RunETL(ctx, cfg, neoLoader)
		log.Printf("Aguardando %d segundos até o próximo ciclo.\n", intervalSec)
		time.Sleep(time.Duration(intervalSec) * time.Second)
	}
}

func RunETL(ctx context.Context, cfg *config.Config, neoLoader *loader.Neo4jLoader) {
	countries := reader.ReadCountries(cfg.CountriesPath + "/countries.csv")
	vaccines := reader.ReadVaccines(cfg.CountriesPath + "/vaccines.csv")
	cases := reader.ReadCovidCases(cfg.CountriesPath + "/covid_cases.csv")
	vaccs := reader.ReadVaccinations(cfg.CountriesPath + "/vaccinations.csv")
	approvals := reader.ReadVaccineApprovals(cfg.CountriesPath + "/vaccine_approvals.csv")
	relations := reader.ReadCountryVaccines(cfg.CountriesPath + "/country_vaccines.csv")

	log.Printf("Países: %d | Vacinas: %d | Casos: %d | Vacinações: %d | Aprovações: %d | Relações: %d",
		len(countries), len(vaccines), len(cases), len(vaccs), len(approvals), len(relations))

	neoLoader.CreateConstraints(ctx)

	neoLoader.LoadCountries(ctx, countries)
	neoLoader.LoadVaccines(ctx, vaccines)
	neoLoader.LoadCovidCases(ctx, cases)
	neoLoader.LoadVaccinationStats(ctx, vaccs)
	neoLoader.LoadVaccineApprovals(ctx, approvals)

	neoLoader.LinkVaccineApprovals(ctx, approvals)
	neoLoader.LinkCountryVaccines(ctx, relations)

	log.Println("Carga concluída com sucesso.")
}
