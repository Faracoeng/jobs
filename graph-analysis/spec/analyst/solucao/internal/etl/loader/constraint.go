package loader

import (
	"context"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func (l *Neo4jLoader) CreateConstraints(ctx context.Context) {
	session := l.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	queries := []string{
		`CREATE CONSTRAINT country_iso3_unique IF NOT EXISTS 
		 FOR (c:Country) REQUIRE c.iso3 IS UNIQUE`,

		`CREATE CONSTRAINT vaccine_name_unique IF NOT EXISTS 
		 FOR (v:Vaccine) REQUIRE v.name IS UNIQUE`,

		`CREATE CONSTRAINT covid_case_unique IF NOT EXISTS 
		 FOR (cc:CovidCase) REQUIRE (cc.date, cc.iso3) IS UNIQUE`,

		`CREATE CONSTRAINT vaccination_unique IF NOT EXISTS 
		 FOR (vs:VaccinationStats) REQUIRE (vs.date, vs.iso3) IS UNIQUE`,

		`CREATE CONSTRAINT vaccine_approval_unique IF NOT EXISTS 
		 FOR (va:VaccineApproval) REQUIRE (va.vaccine, va.date) IS UNIQUE`,
	}

	for _, query := range queries {
		_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			return tx.Run(ctx, query, nil)
		})
		if err != nil {
			log.Printf("Erro ao criar constraint: %v", err)
		}
	}
}
