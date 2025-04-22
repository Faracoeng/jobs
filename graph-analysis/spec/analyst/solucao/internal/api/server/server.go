package server

import (
	"context"
	"fmt"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/api/routes"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/config"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/repository/neo4j"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func getCFG() *config.Config {
	cfg := config.LoadEnv()

	return cfg

}

func NewServer() *Server {
	cfg := getCFG()
	ctx := context.Background()

	driver, err := neo4j.GetDriver(ctx, cfg.URI, cfg.Username, cfg.Password)
	if err != nil {
		panic("Erro ao conectar ao Neo4j: " + err.Error())
	}

	r := routes.SetupRouter(driver)
	return &Server{router: r}
}

func (s *Server) Run() error {
	cfg := getCFG()
	return s.router.Run(fmt.Sprintf("%s:%s", cfg.APIHost, cfg.APIPort))

}
