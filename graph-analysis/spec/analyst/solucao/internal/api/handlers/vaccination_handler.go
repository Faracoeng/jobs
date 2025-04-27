package handler

import (
	"net/http"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/usecase/api"
	"github.com/gin-gonic/gin"
)

type VaccinationHandler struct {
	usecase *usecase.GetVaccinatedUseCase
}

func NewVaccinationHandler(uc *usecase.GetVaccinatedUseCase) *VaccinationHandler {
	return &VaccinationHandler{usecase: uc}
}
// Quantas pessoas foram vacinadas com pelo menos uma dose em um determinado país em uma data específica?
func (h *VaccinationHandler) GetVaccinated(c *gin.Context) {
	iso3 := c.Query("iso3")
	date := c.Query("date")

	if iso3 == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros iso3 e date são obrigatórios"})
		return
	}

	input := usecase.GetVaccinatedInputDTO{
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
