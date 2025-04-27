package handler

import (
	"net/http"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/usecase/api"
	"github.com/gin-gonic/gin"
)

type CovidHandler struct {
	usecase *usecase.GetCovidStatsUseCase
}

func NewCovidHandler(uc *usecase.GetCovidStatsUseCase) *CovidHandler {
	return &CovidHandler{usecase: uc}
}
// Qual foi o total acumulado de casos e mortes de Covid-19 em um país específico em uma data determinada?
func (h *CovidHandler) GetCovidStats(c *gin.Context) {
	iso3 := c.Query("iso3")
	date := c.Query("date")

	if iso3 == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros iso3 e date são obrigatórios"})
		return
	}

	input := usecase.GetCovidStatsInputDTO{
		ISO3: iso3,
		Date: date,
	}

	result, err := h.usecase.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
