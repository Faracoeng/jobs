package reader

import (
	"log"
	"strconv"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/entity"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/etl/util"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/etl/source"
)

func ReadCountries(ds source.DataSource) []entity.Country {
	rows, err := ds.ReadAll()
	if err != nil {
		log.Printf("Erro ao ler countries: %v", err)
		return nil
	}

	var result []entity.Country
	for i, row := range rows {
		if len(row) < 2 {
			log.Printf("Linha %d inválida: %+v", i+2, row)
			continue
		}
		result = append(result, entity.Country{
			ISO3: row[0],
			Name: row[1],
		})
	}
	return result
}

func ReadVaccines(ds source.DataSource) []entity.Vaccine {
	rows, err := ds.ReadAll()
	if err != nil {
		log.Printf("Erro ao ler vaccines: %v", err)
		return nil
	}

	var result []entity.Vaccine
	for i, row := range rows {
		if len(row) < 1 {
			log.Printf("Linha %d inválida: %+v", i+2, row)
			continue
		}
		result = append(result, entity.Vaccine{Name: row[0]})
	}
	return result
}

func ReadCovidCases(ds source.DataSource) []entity.CovidCase {
	rows, err := ds.ReadAll()
	if err != nil {
		log.Printf("Erro ao ler covid cases: %v", err)
		return nil
	}

	var result []entity.CovidCase
	for i, row := range rows {
		if len(row) < 4 {
			log.Printf("Linha %d inválida: %+v", i+2, row)
			continue
		}
		date, err := util.ParseDate(row[1])
		if err != nil {
			log.Printf("Erro parseando data na linha %d: %v", i+2, err)
			continue
		}
		cases, _ := strconv.Atoi(row[2])
		deaths, _ := strconv.Atoi(row[3])
		result = append(result, entity.CovidCase{
			ISO3:        row[0],
			Date:        date,
			TotalCases:  cases,
			TotalDeaths: deaths,
		})
	}
	return result
}

func ReadVaccinations(ds source.DataSource) []entity.VaccinationStat {
	rows, err := ds.ReadAll()
	if err != nil {
		log.Printf("Erro ao ler vaccination stats: %v", err)
		return nil
	}

	var result []entity.VaccinationStat
	for i, row := range rows {
		if len(row) < 3 {
			log.Printf("Linha %d inválida: %+v", i+2, row)
			continue
		}
		date, err := util.ParseDate(row[1])
		if err != nil {
			log.Printf("Erro parseando data na linha %d: %v", i+2, err)
			continue
		}
		total, _ := strconv.Atoi(row[2])
		result = append(result, entity.VaccinationStat{
			ISO3:            row[0],
			Date:            date,
			TotalVaccinated: total,
		})
	}
	return result
}

func ReadVaccineApprovals(ds source.DataSource) []entity.VaccineApproval {
	rows, _ := ds.ReadAll()
	var result []entity.VaccineApproval

	for i, row := range rows {
		if len(row) < 2 {
			log.Printf("Linha %d inválida: %+v", i+2, row)
			continue
		}
		date, err := util.ParseDate(row[1])
		if err != nil {
			log.Printf("Erro parseando data na linha %d: %v", i+2, err)
			continue
		}
		result = append(result, entity.VaccineApproval{
			VaccineName: row[0],
			Date:        date,
		})
	}
	return result
}


func ReadCountryVaccines(ds source.DataSource) []entity.CountryVaccine {
	rows, err := ds.ReadAll()
	if err != nil {
		log.Printf("Erro ao ler country vaccines: %v", err)
		return nil
	}

	var result []entity.CountryVaccine
	for i, row := range rows {
		if len(row) < 2 {
			log.Printf("Linha %d inválida: %+v", i+2, row)
			continue
		}
		result = append(result, entity.CountryVaccine{
			ISO3:        row[0],
			VaccineName: row[1],
		})
	}
	return result
}
