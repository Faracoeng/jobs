package loader

import (
	"context"
	"log"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (l *Neo4jLoader) LoadVaccineApprovals(ctx context.Context, approvals []entity.VaccineApproval, batchSize int) {
	session := l.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	for i := 0; i < len(approvals); i += batchSize {
		end := i + batchSize
		if end > len(approvals) {
			end = len(approvals)
		}
		batch := approvals[i:end]

		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			for _, a := range batch {
				query := `
				MERGE (va:VaccineApproval {date: date($date)})
				`
				_, err := tx.Run(ctx, query, map[string]interface{}{
					"date": a.Date.Format("2006-01-02"),
				})
				if err != nil {
					return nil, err
				}
			}
			return nil, nil
		})

		if err != nil {
			log.Printf("Erro ao inserir batch de aprovações: %v", err)
		}
	}
}
