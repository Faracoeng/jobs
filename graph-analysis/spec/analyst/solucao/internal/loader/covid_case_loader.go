package loader

import (
	"context"
	"fmt"
	"log"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (l *Neo4jLoader) LoadCovidCases(ctx context.Context, data []entity.CovidCase) {
	session := l.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	for _, c := range data {
		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			query := `
			MATCH (co:Country {iso3: $iso3})
			MERGE (cc:CovidCase {date: date($date)})
			SET cc.totalCases = $cases, cc.totalDeaths = $deaths
			MERGE (co)-[:HAS_CASE]->(cc)`
			params := map[string]interface{}{
				"iso3":   c.ISO3,
				"date":   c.Date.Format("2006-01-02"),
				"cases":  c.TotalCases,
				"deaths": c.TotalDeaths,
			}
			return tx.Run(ctx, query, params)
		})
		if err != nil {
			log.Printf("Erro ao inserir caso de Covid (%s): %v", c.Date.Format("2006-01-02"), err)
		}
	}

	fmt.Printf("Inserção de %d casos de Covid concluída.\n", len(data))
}
