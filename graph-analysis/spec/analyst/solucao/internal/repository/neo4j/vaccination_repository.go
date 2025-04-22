package neo4j

import (
	"context"
	"fmt"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/api/handlers"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type vaccinationRepository struct {
	driver neo4j.DriverWithContext
}

func NewVaccinationRepository(driver neo4j.DriverWithContext) handler.VaccinationRepository {
	return &vaccinationRepository{driver: driver}
}

func (r *vaccinationRepository) FetchVaccinated(iso3 string, date string) (int, error) {
	ctx := context.Background()
	session := r.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (c:Country {iso3: $iso3})-[:VACCINATED_ON]->(v:VaccinationStats {date: date($date)})
			RETURN v.totalVaccinated`
		params := map[string]interface{}{
			"iso3": iso3,
			"date": date,
		}
		record, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, err
		}
		if record.Next(ctx) {
			value := record.Record().Values[0].(int64)
			return int(value), nil
		}
		return nil, fmt.Errorf("dados não encontrados")
	})
	if err != nil {
		return 0, err
	}
	return result.(int), nil
}
