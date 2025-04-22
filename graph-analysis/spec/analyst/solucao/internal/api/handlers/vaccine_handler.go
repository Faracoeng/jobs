package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"


)

type VaccineHandler struct {
	repo VaccineRepository
}

type VaccineApprovalDatesResponse struct {
	Dates []string `json:"dates"`
}

type CountryListResponse struct {
	Countries []string `json:"countries"`
}




func NewVaccineHandler(repo VaccineRepository) *VaccineHandler {
	return &VaccineHandler{repo: repo}
}

type VaccinesResponse struct {
	Vaccines []string `json:"vaccines"`
}

func (h *VaccineHandler) GetVaccinesByCountry(c *gin.Context) {
	iso3 := c.Query("iso3")
	if iso3 == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro iso3 é obrigatório"})
		return
	}

	vaccines, err := h.repo.GetVaccinesByCountry(iso3)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar vacinas"})
		return
	}

	c.JSON(http.StatusOK, VaccinesResponse{Vaccines: vaccines})
}


func (h *VaccineHandler) GetApprovalDates(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro name é obrigatório"})
		return
	}

	dates, err := h.repo.GetApprovalDates(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar datas de aprovação"})
		return
	}

	// Formatar para string
	var formatted []string
	for _, d := range dates {
		formatted = append(formatted, d.Format("2006-01-02"))
	}

	c.JSON(http.StatusOK, VaccineApprovalDatesResponse{Dates: formatted})
}


func (h *VaccineHandler) GetCountriesByVaccine(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetro name é obrigatório"})
		return
	}

	countries, err := h.repo.GetCountriesByVaccine(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar países"})
		return
	}

	c.JSON(http.StatusOK, CountryListResponse{Countries: countries})
}