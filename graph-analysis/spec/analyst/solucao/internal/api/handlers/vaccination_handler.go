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
// GetVaccinated retorna o total de pessoas vacinadas em um país na data informada
// @Summary Total de pessoas vacinadas por país e data
// @Description Retorna o número de pessoas vacinadas com pelo menos uma dose em um país em uma data específica
// @Tags vaccination
// @Accept json
// @Produce json
// @Param iso3 query string true "Código ISO3 do país"
// @Param date query string true "Data no formato YYYY-MM-DD"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /vaccination [get]

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
