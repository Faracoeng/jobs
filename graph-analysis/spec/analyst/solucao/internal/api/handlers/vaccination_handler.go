package handler


import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type VaccinationHandler struct {
	repo VaccinationRepository
}

func NewVaccinationHandler(repo VaccinationRepository) *VaccinationHandler {
	return &VaccinationHandler{repo: repo}
}

func (h *VaccinationHandler) GetVaccinated(c *gin.Context) {
	iso3 := c.Query("iso3")
	date := c.Query("date")

	if iso3 == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parâmetros iso3 e date são obrigatórios"})
		return
	}

	total, err := h.repo.FetchVaccinated(iso3, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar dados de vacinação"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_vaccinated": total})
}
