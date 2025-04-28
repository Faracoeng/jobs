package loader

import (
	"context"
	"log"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (l *Neo4jLoader) LoadCovidCases(ctx context.Context, cases []entity.CovidCase, batchSize int) {
	session := l.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	for i := 0; i < len(cases); i += batchSize {
		end := i + batchSize
		if end > len(cases) {
			end = len(cases)
		}
		batch := cases[i:end]

		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			for _, c := range batch {
				query := `
				MERGE (co:Country {iso3: $iso3})
				MERGE (cc:CovidCase {date: date($date), iso3: $iso3})
				SET cc.totalCases = $cases, cc.totalDeaths = $deaths
				MERGE (co)-[:HAS_CASE]->(cc)
				`
				_, err := tx.Run(ctx, query, map[string]interface{}{
					"iso3":   c.ISO3,
					"date":   c.Date.Format("2006-01-02"),
					"cases":  c.TotalCases,
					"deaths": c.TotalDeaths,
				})
				if err != nil {
					return nil, err
				}
			}
			return nil, nil
		})

		if err != nil {
			log.Printf("Erro ao inserir batch de casos de Covid: %v", err)
		}
	}
}
