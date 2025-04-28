package loader

import (
	"context"
	"log"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (l *Neo4jLoader) LoadCountries(ctx context.Context, countries []entity.Country, batchSize int) {
	session := l.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	for i := 0; i < len(countries); i += batchSize {
		end := i + batchSize
		if end > len(countries) {
			end = len(countries)
		}
		batch := countries[i:end]

		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			for _, c := range batch {
				query := `
				MERGE (country:Country {iso3: $iso3})
				SET country.name = $name`
				_, err := tx.Run(ctx, query, map[string]interface{}{
					"iso3": c.ISO3,
					"name": c.Name,
				})
				if err != nil {
					return nil, err
				}
			}
			return nil, nil
		})

		if err != nil {
			log.Printf("Erro ao inserir batch de países: %v", err)
		}
	}
}
