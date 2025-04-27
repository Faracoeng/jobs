package reader

import (
	"log"
	"strconv"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/entity"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/util"
)

func ReadCountries(path string) []entity.Country {
	rows := util.ReadCSVFile(path)
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

func ReadVaccines(path string) []entity.Vaccine {
	rows := util.ReadCSVFile(path)
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

func ReadCovidCases(path string) []entity.CovidCase {
	rows := util.ReadCSVFile(path)
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

func ReadVaccinations(path string) []entity.VaccinationStat {
	rows := util.ReadCSVFile(path)
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

func ReadVaccineApprovals(path string) []entity.VaccineApproval {
	rows := util.ReadCSVFile(path)
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

func ReadCountryVaccines(path string) []entity.CountryVaccine {
	rows := util.ReadCSVFile(path)
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
