package server

import (
	"context"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/api/routes"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/config"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/repository/neo4j"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	cfg := config.LoadEnv()
	ctx := context.Background()

	driver, err := neo4j.GetDriver(ctx, cfg.URI, cfg.Username, cfg.Password)
	if err != nil {
		panic("Erro ao conectar ao Neo4j: " + err.Error())
	}

	r := routes.SetupRouter(driver)
	return &Server{router: r}
}

func (s *Server) Run() error {
	return s.router.Run(":8080")
}
