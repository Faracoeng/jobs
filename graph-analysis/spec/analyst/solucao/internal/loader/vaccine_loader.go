package loader

import (
	"context"
	"log"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/model"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (l *Neo4jLoader) LoadVaccines(ctx context.Context, data []model.Vaccine) {
	session := l.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	for _, v := range data {
		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			query := `MERGE (v:Vaccine {name: $name})`
			_, err := tx.Run(ctx, query, map[string]interface{}{"name": v.Name})
			return nil, err
		})
		if err != nil {
			log.Println("Erro ao inserir vacina:", err)
		}
	}
}
