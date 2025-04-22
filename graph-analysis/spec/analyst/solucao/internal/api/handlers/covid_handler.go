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
