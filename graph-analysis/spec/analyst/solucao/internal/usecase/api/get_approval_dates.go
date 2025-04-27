package usecase

import (
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/entity"
)

type GetApprovalDatesInputDTO struct {
	VaccineName string
}

type GetApprovalDatesOutputDTO struct {
	Dates []string
}

type GetApprovalDatesUseCase struct {
	VaccineRepo entity.VaccineRepository
}

func NewGetApprovalDatesUseCase(repo entity.VaccineRepository) *GetApprovalDatesUseCase {
	return &GetApprovalDatesUseCase{VaccineRepo: repo}
}

func (uc *GetApprovalDatesUseCase) Execute(input GetApprovalDatesInputDTO) (*GetApprovalDatesOutputDTO, error) {
	dates, err := uc.VaccineRepo.GetApprovalDates(input.VaccineName)
	if err != nil {
		return nil, err
	}

	var formatted []string
	for _, d := range dates {
		formatted = append(formatted, d.Format("2006-01-02"))
	}

	return &GetApprovalDatesOutputDTO{
		Dates: formatted,
	}, nil
}
