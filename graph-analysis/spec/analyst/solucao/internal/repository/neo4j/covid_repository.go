package neo4j

import (
	"context"
	"fmt"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/api/handlers"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/entity"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type covidRepository struct {
	driver neo4j.DriverWithContext
}

func NewCovidRepository(driver neo4j.DriverWithContext) handler.CovidRepository {
	return &covidRepository{driver: driver}
}

func (r *covidRepository) FetchStats(iso3 string, date string) (*entity.CovidCase, error) {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (c:Country {iso3: $iso3})-[:HAS_CASE]->(cc:CovidCase {date: date($date)})
			RETURN cc.totalCases, cc.totalDeaths
		`
		params := map[string]interface{}{
			"iso3": iso3,
			"date": date,
		}

		record, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		if record.Next(ctx) {
			values := record.Record().Values
			return &entity.CovidCase{
				ISO3:        iso3,
				TotalCases:  int(values[0].(int64)),
				TotalDeaths: int(values[1].(int64)),
			}, nil
		}
		return nil, fmt.Errorf("sem dados encontrados")
	})
	if err != nil {
		return nil, err
	}
	return result.(*entity.CovidCase), nil
}
