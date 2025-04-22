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
// GetVaccinated retorna o total de vacinados por país e data
// @Summary Total de vacinados
// @Description Retorna o total de vacinas aplicadas por país e data
// @Tags vaccination
// @Accept json
// @Produce json
// @Param iso3 query string true "Código ISO3 do país"
// @Param date query string true "Data no formato YYYY-MM-DD"
// @Success 200 {object} map[string]int
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /vaccinated [get]
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
