package entity

import (
	"time"
)
// Qual foi o total acumulado de casos e mortes de Covid-19 em um país específico em uma data determinada?
type CovidRepository interface {
	FetchStats(iso3 string, date string) (*CovidCase, error)
}
// Quantas pessoas foram vacinadas com pelo menos uma dose em um determinado país em uma data específica?
type VaccinationRepository interface {
	FetchVaccinated(iso3 string, date string) (int, error)
}

type VaccineRepository interface {
	GetVaccinesByCountry(iso3 string) ([]string, error)
	GetApprovalDates(name string) ([]time.Time, error)
	GetCountriesByVaccine(vaccine string) ([]string, error)
}
