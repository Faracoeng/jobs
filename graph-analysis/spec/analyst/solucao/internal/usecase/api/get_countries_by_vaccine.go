package usecase

import (
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/entity"
)

type GetCountriesByVaccineInputDTO struct {
	VaccineName string
}

type GetCountriesByVaccineOutputDTO struct {
	Countries []string
}

type GetCountriesByVaccineUseCase struct {
	VaccineRepo entity.VaccineRepository
}

func NewGetCountriesByVaccineUseCase(repo entity.VaccineRepository) *GetCountriesByVaccineUseCase {
	return &GetCountriesByVaccineUseCase{VaccineRepo: repo}
}

func (uc *GetCountriesByVaccineUseCase) Execute(input GetCountriesByVaccineInputDTO) (*GetCountriesByVaccineOutputDTO, error) {
	countries, err := uc.VaccineRepo.GetCountriesByVaccine(input.VaccineName)
	if err != nil {
		return nil, err
	}

	return &GetCountriesByVaccineOutputDTO{
		Countries: countries,
	}, nil
}
