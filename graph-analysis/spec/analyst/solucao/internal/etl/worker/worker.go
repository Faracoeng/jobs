package worker

import (
	"context"
	"log"
	"strconv"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/entity"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/etl/reader"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/etl/loader"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/etl/util"
)

// ProcessCountries lê países e carrega no Neo4j
func ProcessCountries(ctx context.Context, lines <-chan reader.Line, neo4jLoader *loader.Neo4jLoader, batchSize int) {
	var countries []entity.Country

	for line := range lines {
		if line.Err != nil {
			log.Printf("Erro lendo linha de país: %v", line.Err)
			continue
		}
		if len(line.Data) < 2 {
			log.Printf("Linha de país inválida: %+v", line.Data)
			continue
		}
		countries = append(countries, entity.Country{
			ISO3: line.Data[0],
			Name: line.Data[1],
		})
	}

	if len(countries) > 0 {
		neo4jLoader.LoadCountries(ctx, countries, batchSize)
	}
}

// ProcessVaccines lê vacinas e carrega no Neo4j
func ProcessVaccines(ctx context.Context, lines <-chan reader.Line, neo4jLoader *loader.Neo4jLoader, batchSize int) {
	var vaccines []entity.Vaccine

	for line := range lines {
		if line.Err != nil {
			log.Printf("Erro lendo linha de vacina: %v", line.Err)
			continue
		}
		if len(line.Data) < 1 {
			log.Printf("Linha de vacina inválida: %+v", line.Data)
			continue
		}
		vaccines = append(vaccines, entity.Vaccine{
			Name: line.Data[0],
		})
	}

	if len(vaccines) > 0 {
		neo4jLoader.LoadVaccines(ctx, vaccines, batchSize)
	}
}

// ProcessCovidCases lê casos de Covid-19 e carrega no Neo4j
func ProcessCovidCases(ctx context.Context, lines <-chan reader.Line, neo4jLoader *loader.Neo4jLoader, batchSize int) {
	var cases []entity.CovidCase

	for line := range lines {
		if line.Err != nil {
			log.Printf("Erro lendo linha de caso de Covid: %v", line.Err)
			continue
		}
		if len(line.Data) < 4 {
			log.Printf("Linha de caso de Covid inválida: %+v", line.Data)
			continue
		}
		date, err := util.ParseDate(line.Data[1])
		if err != nil {
			log.Printf("Erro parseando data: %v", err)
			continue
		}
		totalCases, _ := strconv.Atoi(line.Data[2])
		totalDeaths, _ := strconv.Atoi(line.Data[3])

		cases = append(cases, entity.CovidCase{
			ISO3:        line.Data[0],
			Date:        date,
			TotalCases:  totalCases,
			TotalDeaths: totalDeaths,
		})
	}

	if len(cases) > 0 {
		neo4jLoader.LoadCovidCases(ctx, cases, batchSize)
	}
}

// ProcessVaccinationStats lê estatísticas de vacinação e carrega no Neo4j
func ProcessVaccinationStats(ctx context.Context, lines <-chan reader.Line, neo4jLoader *loader.Neo4jLoader, batchSize int) {
	var stats []entity.VaccinationStat

	for line := range lines {
		if line.Err != nil {
			log.Printf("Erro lendo linha de vacinação: %v", line.Err)
			continue
		}
		if len(line.Data) < 3 {
			log.Printf("Linha de vacinação inválida: %+v", line.Data)
			continue
		}
		date, err := util.ParseDate(line.Data[1])
		if err != nil {
			log.Printf("Erro parseando data: %v", err)
			continue
		}
		totalVaccinated, _ := strconv.Atoi(line.Data[2])

		stats = append(stats, entity.VaccinationStat{
			ISO3:            line.Data[0],
			Date:            date,
			TotalVaccinated: totalVaccinated,
		})
	}

	if len(stats) > 0 {
		neo4jLoader.LoadVaccinationStats(ctx, stats, batchSize)
	}
}

// ProcessVaccineApprovals lê aprovações de vacinas e carrega no Neo4j
func ProcessVaccineApprovals(ctx context.Context, lines <-chan reader.Line, neo4jLoader *loader.Neo4jLoader, batchSize int) []entity.VaccineApproval {
	var approvals []entity.VaccineApproval

	for line := range lines {
		if line.Err != nil {
			log.Printf("Erro lendo linha de aprovação de vacina: %v", line.Err)
			continue
		}
		if len(line.Data) < 2 {
			log.Printf("Linha de aprovação inválida: %+v", line.Data)
			continue
		}
		date, err := util.ParseDate(line.Data[1])
		if err != nil {
			log.Printf("Erro parseando data: %v", err)
			continue
		}

		approvals = append(approvals, entity.VaccineApproval{
			VaccineName: line.Data[0],
			Date:        date,
		})
	}

	if len(approvals) > 0 {
		neo4jLoader.LoadVaccineApprovals(ctx, approvals, batchSize)
	}

	return approvals
}


// ProcessCountryVaccines lê relações país-vacina e carrega no Neo4j
func ProcessCountryVaccines(ctx context.Context, lines <-chan reader.Line, neo4jLoader *loader.Neo4jLoader, batchSize int) {
	var relations []entity.CountryVaccine

	for line := range lines {
		if line.Err != nil {
			log.Printf("Erro lendo linha de relação país-vacina: %v", line.Err)
			continue
		}
		if len(line.Data) < 2 {
			log.Printf("Linha de relação país-vacina inválida: %+v", line.Data)
			continue
		}

		relations = append(relations, entity.CountryVaccine{
			ISO3:        line.Data[0],
			VaccineName: line.Data[1],
		})
	}

	if len(relations) > 0 {
		neo4jLoader.LinkCountryVaccines(ctx, relations, batchSize)
	}
}