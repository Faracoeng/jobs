package loader

import (
	"context"
	"log"

	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (l *Neo4jLoader) LinkCountryVaccines(ctx context.Context, relations []entity.CountryVaccine, batchSize int) {
	session := l.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	for i := 0; i < len(relations); i += batchSize {
		end := i + batchSize
		if end > len(relations) {
			end = len(relations)
		}
		batch := relations[i:end]

		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			for _, r := range batch {
				query := `
				MERGE (c:Country {iso3: $iso3})
				MERGE (v:Vaccine {name: $vaccine})
				MERGE (c)-[:USES]->(v)
				`
				_, err := tx.Run(ctx, query, map[string]interface{}{
					"iso3":    r.ISO3,
					"vaccine": r.VaccineName,
				})
				if err != nil {
					return nil, err
				}
			}
			return nil, nil
		})

		if err != nil {
			log.Printf("Erro ao criar relações país-vacina: %v", err)
		}
	}
}
func (l *Neo4jLoader) LinkVaccinesToApprovals(ctx context.Context, approvals []entity.VaccineApproval, batchSize int) {
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
				MERGE (v:Vaccine {name: $vaccine})
				MERGE (va:VaccineApproval {date: date($date)})
				MERGE (v)-[:APPROVED_ON]->(va)
				`
				_, err := tx.Run(ctx, query, map[string]interface{}{
					"vaccine": a.VaccineName,
					"date":    a.Date.Format("2006-01-02"),
				})
				if err != nil {
					return nil, err
				}
			}
			return nil, nil
		})

		if err != nil {
			log.Printf("Erro ao criar relações Vaccine -> Approval: %v", err)
		}
	}
}

