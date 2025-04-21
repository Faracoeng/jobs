// internal/loader/relationship_loader.go
package loader

import (
	"context"
	"log"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/model"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (l *Neo4jLoader) LinkVaccineApprovals(ctx context.Context, approvals []model.VaccineApproval) {
	session := l.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	for _, a := range approvals {
		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			query := `
			MATCH (v:Vaccine {name: $vaccine})
			MATCH (a:VaccineApproval {vaccine: $vaccine, date: date($date)})
			MERGE (v)-[:APPROVED_ON]->(a)`
			params := map[string]interface{}{
				"vaccine": a.VaccineName,
				"date":    a.Date.Format("2006-01-02"),
			}
			return tx.Run(ctx, query, params)
		})
		if err != nil {
			log.Printf("Erro ao criar relação de aprovação para %s: %v", a.VaccineName, err)
		}
	}
}

func (l *Neo4jLoader) LinkCountryVaccines(ctx context.Context, relations []model.CountryVaccine) {
	session := l.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	for _, r := range relations {
		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			query := `
			MATCH (c:Country {iso3: $iso3})
			MATCH (v:Vaccine {name: $vaccine})
			MERGE (c)-[:USES]->(v)`
			params := map[string]interface{}{
				"iso3":    r.ISO3,
				"vaccine": r.VaccineName,
			}
			return tx.Run(ctx, query, params)
		})
		if err != nil {
			log.Printf("Erro ao criar relação de uso de vacina (%s → %s): %v", r.ISO3, r.VaccineName, err)
		}
	}
}
