package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/config"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/etl/loader"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/etl/reader"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/etl/source"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/etl/worker"
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
	readyFlagPath := filepath.Join(cfg.CountriesPath, "ready.flag")
	if _, err := os.Stat(readyFlagPath); os.IsNotExist(err) {
		log.Println("Arquivo 'ready.flag' não encontrado. Pulando execução do ETL.")
		return
	}

	log.Println("Arquivo 'ready.flag' encontrado. Iniciando processamento...")

	log.Println("Criando constraints no Neo4j...")
	neoLoader.CreateConstraints(ctx)

	log.Println("Processando países...")
	countrySource := &source.CSVSource{Path: filepath.Join(cfg.CountriesPath, "countries.csv")}
	countryLines := reader.StreamData(countrySource)
	worker.ProcessCountries(ctx, countryLines, neoLoader, cfg.ETLBatchSize)

	log.Println("Processando vacinas...")
	vaccineSource := &source.CSVSource{Path: filepath.Join(cfg.CountriesPath, "vaccines.csv")}
	vaccineLines := reader.StreamData(vaccineSource)
	worker.ProcessVaccines(ctx, vaccineLines, neoLoader, cfg.ETLBatchSize)

	log.Println("Processando casos de Covid...")
	caseSource := &source.CSVSource{Path: filepath.Join(cfg.CountriesPath, "covid_cases.csv")}
	caseLines := reader.StreamData(caseSource)
	worker.ProcessCovidCases(ctx, caseLines, neoLoader, cfg.ETLBatchSize)

	log.Println("Processando vacinação...")
	vaccinationSource := &source.CSVSource{Path: filepath.Join(cfg.CountriesPath, "vaccinations.csv")}
	vaccinationLines := reader.StreamData(vaccinationSource)
	worker.ProcessVaccinationStats(ctx, vaccinationLines, neoLoader, cfg.ETLBatchSize)

	log.Println("Processando aprovações de vacinas...")
	approvalSource := &source.CSVSource{Path: filepath.Join(cfg.CountriesPath, "vaccine_approvals.csv")}
	approvalLines := reader.StreamData(approvalSource)
	approvals := worker.ProcessVaccineApprovals(ctx, approvalLines, neoLoader, cfg.ETLBatchSize)

	log.Println("Processando relações país-vacina...")
	relationSource := &source.CSVSource{Path: filepath.Join(cfg.CountriesPath, "country_vaccines.csv")}
	relationLines := reader.StreamData(relationSource)
	worker.ProcessCountryVaccines(ctx, relationLines, neoLoader, cfg.ETLBatchSize)

	log.Println("Criando relações entre Vacinas e Aprovações...")
	neoLoader.LinkVaccinesToApprovals(ctx, approvals, cfg.ETLBatchSize)

	log.Println("Carga concluída com sucesso.")

	// Remove o ready.flag depois da carga
	err := os.Remove(readyFlagPath)
	if err != nil {
		log.Printf("Erro ao remover o arquivo 'ready.flag': %v", err)
	} else {
		log.Println("Arquivo 'ready.flag' removido com sucesso após o processamento.")
	}
}
