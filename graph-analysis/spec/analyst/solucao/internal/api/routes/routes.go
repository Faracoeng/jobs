package routes

import (
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/api/handlers"
	repoNeo4j "github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/infra/repository/neo4j"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/usecase/api"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

func SetupRouter(driver neo4j.DriverWithContext) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Criar Repository → UseCase → Handler
	covidRepo := repoNeo4j.NewCovidRepository(driver)
	covidUseCase := usecase.NewGetCovidStatsUseCase(covidRepo)
	covidHandler := handler.NewCovidHandler(covidUseCase)
	r.GET("/covid-stats", covidHandler.GetCovidStats)


	vaccRepo := repoNeo4j.NewVaccinationRepository(driver)
	vaccUseCase := usecase.NewGetVaccinatedUseCase(vaccRepo)
	vaccHandler := handler.NewVaccinationHandler(vaccUseCase)
	r.GET("/vaccination", vaccHandler.GetVaccinated)

	vaccineRepo := repoNeo4j.NewVaccineRepository(driver)
	vaccineHandler := handler.NewVaccineHandler(vaccineRepo)
	r.GET("/vaccines", vaccineHandler.GetVaccinesByCountry)
	r.GET("/approval-dates", vaccineHandler.GetApprovalDates)
	r.GET("/countries-by-vaccine", vaccineHandler.GetCountriesByVaccine)


	// adicionar rota para swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
