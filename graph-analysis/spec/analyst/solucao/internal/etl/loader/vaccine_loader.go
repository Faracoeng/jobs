package loader

import (
	"context"
	"log"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (l *Neo4jLoader) LoadVaccines(ctx context.Context, vaccines []entity.Vaccine, batchSize int) {
	session := l.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	for i := 0; i < len(vaccines); i += batchSize {
		end := i + batchSize
		if end > len(vaccines) {
			end = len(vaccines)
		}
		batch := vaccines[i:end]

		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			for _, v := range batch {
				query := `MERGE (v:Vaccine {name: $name})`
				_, err := tx.Run(ctx, query, map[string]interface{}{
					"name": v.Name,
				})
				if err != nil {
					return nil, err
				}
			}
			return nil, nil
		})

		if err != nil {
			log.Println("Erro ao inserir batch de vacinas:", err)
		}
	}
}
