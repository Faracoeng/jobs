package usecase

import (
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/entity"
)

type GetVaccinesByCountryInputDTO struct {
	ISO3 string
}

type GetVaccinesByCountryOutputDTO struct {
	Vaccines []string
}

type GetVaccinesByCountryUseCase struct {
	VaccineRepo entity.VaccineRepository
}

func NewGetVaccinesByCountryUseCase(repo entity.VaccineRepository) *GetVaccinesByCountryUseCase {
	return &GetVaccinesByCountryUseCase{VaccineRepo: repo}
}

func (uc *GetVaccinesByCountryUseCase) Execute(input GetVaccinesByCountryInputDTO) (*GetVaccinesByCountryOutputDTO, error) {
	vaccines, err := uc.VaccineRepo.GetVaccinesByCountry(input.ISO3)
	if err != nil {
		return nil, err
	}

	return &GetVaccinesByCountryOutputDTO{
		Vaccines: vaccines,
	}, nil
}
