package main

import (
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/api/routes"
)

func main() {
	r := routes.SetupRouter()
	r.Run(":8080")
}
