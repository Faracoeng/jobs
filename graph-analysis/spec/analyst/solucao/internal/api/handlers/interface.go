package handler

import "github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/model"

type CovidRepository interface {
	FetchStats(iso3 string, date string) (*model.CovidCase, error)
}

type VaccinationRepository interface {
	FetchVaccinated(iso3 string, date string) (int, error)
}

type VaccineRepository interface {
	GetVaccinesByCountry(iso3 string) ([]string, error)
}
