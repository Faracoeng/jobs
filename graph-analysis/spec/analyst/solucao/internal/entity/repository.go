package entity

import (
	"time"
)

type CovidRepository interface {
	FetchStats(iso3 string, date string) (*CovidCase, error)
}

type VaccinationRepository interface {
	FetchVaccinated(iso3 string, date string) (int, error)
}

type VaccineRepository interface {
	GetVaccinesByCountry(iso3 string) ([]string, error)
	GetApprovalDates(name string) ([]time.Time, error)
	GetCountriesByVaccine(vaccine string) ([]string, error)
}
