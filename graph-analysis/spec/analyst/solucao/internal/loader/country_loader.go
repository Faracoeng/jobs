package loader

import (
	"context"
	"log"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/model"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (l *Neo4jLoader) LoadCountries(ctx context.Context, countries []model.Country) {
	session := l.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	for _, c := range countries {
		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			// Tenta encontrar um país existente com o mesmo iso3
			// Se não encontrar, cria um novo país
			query := `
				MERGE (country:Country {iso3: $iso3})
				SET country.name = $name`
			_, err := tx.Run(ctx, query, map[string]interface{}{
				"iso3": c.ISO3,
				"name": c.Name,
			})
			return nil, err
		})

		if err != nil {
			log.Printf("Erro ao inserir país %s: %v", c.ISO3, err)
		}
	}
}
