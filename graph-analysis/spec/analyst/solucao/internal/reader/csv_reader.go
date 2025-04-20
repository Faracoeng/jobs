package reader

import (
	"log"
	"strconv"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/model"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/util"
)

func ReadCountries(path string) []model.Country {
	rows := util.ReadCSVFile(path)
	var result []model.Country

	for i, row := range rows {
		if len(row) < 2 {
			log.Printf("Linha %d inválida: %+v", i+2, row)
			continue
		}
		result = append(result, model.Country{
			ISO3: row[0],
			Name: row[1],
		})
	}
	return result
}

func ReadVaccines(path string) []model.Vaccine {
	rows := util.ReadCSVFile(path)
	var result []model.Vaccine

	for i, row := range rows {
		if len(row) < 1 {
			log.Printf("Linha %d inválida: %+v", i+2, row)
			continue
		}
		result = append(result, model.Vaccine{Name: row[0]})
	}
	return result
}

func ReadCovidCases(path string) []model.CovidCase {
	rows := util.ReadCSVFile(path)
	var result []model.CovidCase

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
		result = append(result, model.CovidCase{
			ISO3:        row[0],
			Date:        date,
			TotalCases:  cases,
			TotalDeaths: deaths,
		})
	}
	return result
}

func ReadVaccinations(path string) []model.VaccinationStat {
	rows := util.ReadCSVFile(path)
	var result []model.VaccinationStat

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
		result = append(result, model.VaccinationStat{
			ISO3:            row[0],
			Date:            date,
			TotalVaccinated: total,
		})
	}
	return result
}

func ReadVaccineApprovals(path string) []model.VaccineApproval {
	rows := util.ReadCSVFile(path)
	var result []model.VaccineApproval

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
		result = append(result, model.VaccineApproval{
			VaccineName: row[0],
			Date:        date,
		})
	}
	return result
}

func ReadCountryVaccines(path string) []model.CountryVaccine {
	rows := util.ReadCSVFile(path)
	var result []model.CountryVaccine

	for i, row := range rows {
		if len(row) < 2 {
			log.Printf("Linha %d inválida: %+v", i+2, row)
			continue
		}
		result = append(result, model.CountryVaccine{
			ISO3:        row[0],
			VaccineName: row[1],
		})
	}
	return result
}
