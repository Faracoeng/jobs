package loader

import (
	"context"
	"log"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (l *Neo4jLoader) LoadVaccineApprovals(ctx context.Context, approvals []entity.VaccineApproval) {
	session := l.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	for _, a := range approvals {
		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			query := `
			MERGE (v:VaccineApproval {vaccine: $vaccine, date: date($date)})`
			params := map[string]interface{}{
				"date":    a.Date.Format("2006-01-02"),
			}
			return tx.Run(ctx, query, params)
		})
		if err != nil {
			log.Printf("Erro ao inserir aprovação para %s (%s): %v", a.Date.Format("2006-01-02"), err)
		}
	}
}
