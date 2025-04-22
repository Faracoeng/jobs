package neo4j

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"time"
)

type vaccineRepository struct {
	driver neo4j.DriverWithContext
}

func NewVaccineRepository(driver neo4j.DriverWithContext) *vaccineRepository {
	return &vaccineRepository{driver: driver}
}

func (r *vaccineRepository) GetVaccinesByCountry(iso3 string) ([]string, error) {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (c:Country {iso3: $iso3})-[:USES]->(v:Vaccine)
			RETURN v.name
		`
		records, err := tx.Run(ctx, query, map[string]interface{}{
			"iso3": iso3,
		})
		if err != nil {
			return nil, err
		}

		var vaccines []string
		for records.Next(ctx) {
			name, _ := records.Record().Get("v.name")
			vaccines = append(vaccines, name.(string))
		}
		return vaccines, nil
	})

	if err != nil {
		return nil, err
	}
	return result.([]string), nil
}


func (r *vaccineRepository) GetApprovalDates(name string) ([]time.Time, error) {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (:Vaccine {name: $name})-[:APPROVED_ON]->(a:VaccineApproval)
			RETURN a.date ORDER BY a.date
		`
		records, err := tx.Run(ctx, query, map[string]any{"name": name})
		if err != nil {
			return nil, err
		}

		var dates []time.Time
		for records.Next(ctx) {
			raw := records.Record().Values[0].(neo4j.Date)
			dates = append(dates, time.Time(raw.Time()))
		}
		return dates, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]time.Time), nil
}


func (r *vaccineRepository) GetCountriesByVaccine(vaccine string) ([]string, error) {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (c:Country)-[:USES]->(v:Vaccine {name: $vaccine})
			RETURN c.iso3
		`
		records, err := tx.Run(ctx, query, map[string]interface{}{"vaccine": vaccine})
		if err != nil {
			return nil, err
		}

		var countries []string
		for records.Next(ctx) {
			iso, _ := records.Record().Get("c.iso3")
			countries = append(countries, iso.(string))
		}
		return countries, nil
	})
	if err != nil {
		return nil, err
	}
	return result.([]string), nil
}
