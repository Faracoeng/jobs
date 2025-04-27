package usecase

import (
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/entity"
)

type GetVaccinatedInputDTO struct {
	ISO3 string
	Date string
}

type GetVaccinatedOutputDTO struct {
	TotalVaccinated int
}

type GetVaccinatedUseCase struct {
	VaccinationRepo entity.VaccinationRepository
}

func NewGetVaccinatedUseCase(repo entity.VaccinationRepository) *GetVaccinatedUseCase {
	return &GetVaccinatedUseCase{VaccinationRepo: repo}
}

func (uc *GetVaccinatedUseCase) Execute(input GetVaccinatedInputDTO) (*GetVaccinatedOutputDTO, error) {
	total, err := uc.VaccinationRepo.FetchVaccinated(input.ISO3, input.Date)
	if err != nil {
		return nil, err
	}

	return &GetVaccinatedOutputDTO{
		TotalVaccinated: total,
	}, nil
}
