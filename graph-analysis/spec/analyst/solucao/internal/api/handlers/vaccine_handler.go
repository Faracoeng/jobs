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
// GetVaccinesByCountry retorna vacinas aplicadas em um país específico
// @Summary Vacinas por país
// @Description Retorna a lista de vacinas aplicadas em um país (por código ISO3)
// @Tags vaccines
// @Accept json
// @Produce json
// @Param iso3 query string true "Código do país (ISO3)"
// @Success 200 {object} VaccinesResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /vaccines [get]

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
// GetApprovalDates retorna as datas de aprovação de uma vacina
// @Summary Datas de aprovação
// @Description Lista de datas em que uma vacina foi aprovada para uso
// @Tags vaccines
// @Accept json
// @Produce json
// @Param name query string true "Nome da vacina"
// @Success 200 {object} VaccineApprovalDatesResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /approval-dates [get]


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
// GetCountriesByVaccine retorna países que aprovaram uma vacina
// @Summary Países por vacina
// @Description Retorna os países onde uma vacina foi aprovada para uso
// @Tags vaccines
// @Accept json
// @Produce json
// @Param name query string true "Nome da vacina"
// @Success 200 {object} CountryListResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /countries-by-vaccine [get]


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