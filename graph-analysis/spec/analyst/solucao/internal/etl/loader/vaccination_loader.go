package loader

import (
	"context"
	"log"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (l *Neo4jLoader) LoadVaccinationStats(ctx context.Context, stats []entity.VaccinationStat, batchSize int) {
	session := l.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	for i := 0; i < len(stats); i += batchSize {
		end := i + batchSize
		if end > len(stats) {
			end = len(stats)
		}
		batch := stats[i:end]

		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			for _, s := range batch {
				query := `
				MERGE (co:Country {iso3: $iso3})
				MERGE (vs:VaccinationStats {date: date($date), iso3: $iso3})
				SET vs.totalVaccinated = $vacc
				MERGE (co)-[:VACCINATED_ON]->(vs)
				`
				_, err := tx.Run(ctx, query, map[string]interface{}{
					"iso3":  s.ISO3,
					"date":  s.Date.Format("2006-01-02"),
					"vacc":  s.TotalVaccinated,
				})
				if err != nil {
					return nil, err
				}
			}
			return nil, nil
		})

		if err != nil {
			log.Printf("Erro ao inserir batch de estatísticas de vacinação: %v", err)
		}
	}
}
