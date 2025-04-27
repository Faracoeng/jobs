package usecase

import (
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/entity"
)

type GetCovidStatsInputDTO struct {
	ISO3 string
	Date string
}

type GetCovidStatsOutputDTO struct {
	TotalCases  int
	TotalDeaths int
}

type GetCovidStatsUseCase struct {
	CovidRepo entity.CovidRepository
}

func NewGetCovidStatsUseCase(repo entity.CovidRepository) *GetCovidStatsUseCase {
	return &GetCovidStatsUseCase{CovidRepo: repo}
}
// Qual foi o total acumulado de casos e mortes de Covid-19 em um país específico em uma data determinada?
func (uc *GetCovidStatsUseCase) Execute(input GetCovidStatsInputDTO) (*GetCovidStatsOutputDTO, error) {
	stats, err := uc.CovidRepo.FetchStats(input.ISO3, input.Date)
	if err != nil {
		return nil, err
	}

	return &GetCovidStatsOutputDTO{
		TotalCases:  stats.TotalCases,
		TotalDeaths: stats.TotalDeaths,
	}, nil
}
