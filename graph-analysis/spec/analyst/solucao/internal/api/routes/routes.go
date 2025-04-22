package routes

import (
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/api/handlers"
	repoNeo4j "github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/repository/neo4j"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func SetupRouter(driver neo4j.DriverWithContext) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Repositório e handler
	covidRepo := repoNeo4j.NewCovidRepository(driver)
	covidHandler := handler.NewCovidHandler(covidRepo)

	// Endpoint da API
	r.GET("/covid-stats", covidHandler.GetCovidStats)


	vaccRepo := repoNeo4j.NewVaccinationRepository(driver)
	vaccHandler := handler.NewVaccinationHandler(vaccRepo)
	r.GET("/vaccination", vaccHandler.GetVaccinated)

	vaccineRepo := repoNeo4j.NewVaccineRepository(driver)
	vaccineHandler := handler.NewVaccineHandler(vaccineRepo)
	r.GET("/vaccines", vaccineHandler.GetVaccinesByCountry)
	r.GET("/approval-dates", vaccineHandler.GetApprovalDates)
	r.GET("/countries-by-vaccine", vaccineHandler.GetCountriesByVaccine)




	return r
}
