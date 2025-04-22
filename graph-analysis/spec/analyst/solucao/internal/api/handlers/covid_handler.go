package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	//"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/model"
)

type CovidHandler struct {
	repo CovidRepository
}

func NewCovidHandler(repo CovidRepository) *CovidHandler {
	return &CovidHandler{repo: repo}
}

type CovidStatsResponse struct {
	TotalCases  int `json:"total_cases"`
	TotalDeaths int `json:"total_deaths"`
}
// GetCovidStats retorna o total de casos e mortes em uma data e país específicos
// @Summary Total de casos e mortes
// @Description Retorna o total acumulado de casos e mortes por COVID-19 em um país e data específicos
// @Tags covid
// @Accept json
// @Produce json
// @Param iso3 query string true "Código ISO3 do país"
// @Param date query string true "Data no formato YYYY-MM-DD"
// @Success 200 {object} CovidStatsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /covid-stats [get]
func (h *CovidHandler) GetCovidStats(c *gin.Context) {
	iso3 := c.Query("iso3")
	date := c.Query("date")

	if iso3 == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros iso3 e date são obrigatórios"})
		return
	}

	stats, err := h.repo.FetchStats(iso3, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar dados"})
		return
	}

	c.JSON(http.StatusOK, CovidStatsResponse{
		TotalCases:  stats.TotalCases,
		TotalDeaths: stats.TotalDeaths,
	})
}
