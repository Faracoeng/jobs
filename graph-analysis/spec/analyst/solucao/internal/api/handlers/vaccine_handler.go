package handler

import (
	"net/http"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/usecase/api"
	"github.com/gin-gonic/gin"
)

type VaccineHandler struct {
	getVaccinesByCountryUC    *usecase.GetVaccinesByCountryUseCase
	getApprovalDatesUC        *usecase.GetApprovalDatesUseCase
	getCountriesByVaccineUC   *usecase.GetCountriesByVaccineUseCase
}

func NewVaccineHandler(
	getVaccinesByCountryUC *usecase.GetVaccinesByCountryUseCase,
	getApprovalDatesUC *usecase.GetApprovalDatesUseCase,
	getCountriesByVaccineUC *usecase.GetCountriesByVaccineUseCase,
) *VaccineHandler {
	return &VaccineHandler{
		getVaccinesByCountryUC:    getVaccinesByCountryUC,
		getApprovalDatesUC:        getApprovalDatesUC,
		getCountriesByVaccineUC:   getCountriesByVaccineUC,
	}
}
// GetVaccinesByCountry retorna as vacinas utilizadas por um país
// @Summary Vacinas usadas em um país
// @Description Lista as vacinas utilizadas por um país específico
// @Tags vaccine
// @Accept json
// @Produce json
// @Param iso3 query string true "Código ISO3 do país"
// @Success 200 {object} map[string][]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /vaccines [get]

func (h *VaccineHandler) GetVaccinesByCountry(c *gin.Context) {
	iso3 := c.Query("iso3")
	if iso3 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro iso3 é obrigatório"})
		return
	}

	input := usecase.GetVaccinesByCountryInputDTO{ISO3: iso3}
	result, err := h.getVaccinesByCountryUC.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
// GetApprovalDates retorna as datas de aprovação das vacinas
// @Summary Datas de aprovação de vacinas
// @Description Lista todas as vacinas e suas datas de aprovação
// @Tags vaccine
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /approval-dates [get]

func (h *VaccineHandler) GetApprovalDates(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro name é obrigatório"})
		return
	}

	input := usecase.GetApprovalDatesInputDTO{VaccineName: name}
	result, err := h.getApprovalDatesUC.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
// GetCountriesByVaccine retorna os países que utilizaram uma vacina específica
// @Summary Países por vacina
// @Description Lista os países que utilizaram uma vacina específica
// @Tags vaccine
// @Accept json
// @Produce json
// @Param vaccine query string true "Nome da vacina"
// @Success 200 {object} map[string][]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /countries-by-vaccine [get]

func (h *VaccineHandler) GetCountriesByVaccine(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro name é obrigatório"})
		return
	}

	input := usecase.GetCountriesByVaccineInputDTO{VaccineName: name}
	result, err := h.getCountriesByVaccineUC.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
