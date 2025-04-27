package loader

import (
	"context"
	"log"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (l *Neo4jLoader) LoadVaccinationStats(ctx context.Context, data []entity.VaccinationStat) {
	session := l.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	for _, v := range data {
		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			query := `
			MATCH (co:Country {iso3: $iso3})
			MERGE (vs:VaccinationStats {date: date($date)})
			SET vs.totalVaccinated = $vacc
			MERGE (co)-[:VACCINATED_ON]->(vs)`
			params := map[string]interface{}{
				"iso3": v.ISO3,
				"date": v.Date.Format("2006-01-02"),
				"vacc": v.TotalVaccinated,
			}
			return tx.Run(ctx, query, params)
		})
		if err != nil {
			log.Printf("Erro ao inserir vacinação para %s em %s: %v", v.ISO3, v.Date.Format("2006-01-02"), err)
		}
	}
}
