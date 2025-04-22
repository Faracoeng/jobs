package main

import (
	"log"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/api/server"
)

func main() {
	s := server.NewServer()
	if err := s.Run(); err != nil {
		log.Fatalf("Erro ao iniciar servidor HTTP: %v", err)
	}
}
