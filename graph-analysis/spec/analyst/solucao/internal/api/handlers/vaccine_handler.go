package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type VaccineHandler struct {
	repo VaccineRepository
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
