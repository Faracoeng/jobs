package main

import (
	"log"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/api/server"
)

// @title API de Análise COVID-19
// @version 1.0
// @description Documentação da API que responde perguntas analíticas sobre COVID-19 e vacinas
// @host localhost:8080
// @BasePath /


func main() {
	s := server.NewServer()
	if err := s.Run(); err != nil {
		log.Fatalf("Erro ao iniciar servidor HTTP: %v", err)
	}
}
