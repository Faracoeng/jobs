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
