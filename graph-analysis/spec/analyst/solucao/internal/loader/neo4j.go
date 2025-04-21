// loader/neo4j.go
package loader

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/config"
)

type Neo4jLoader struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jLoader(driver neo4j.DriverWithContext) *Neo4jLoader {
	return &Neo4jLoader{driver: driver}
}

func (l *Neo4jLoader) Close(ctx context.Context) {
	l.driver.Close(ctx)
}

func NewNeo4jDriver(ctx context.Context, cfg *config.Config) (neo4j.DriverWithContext, error) {
	driver, err := neo4j.NewDriverWithContext(
		cfg.URI,
		neo4j.BasicAuth(cfg.Username, cfg.Password, ""),
	)
	if err != nil {
		return nil, err
	}
	if err := driver.VerifyConnectivity(ctx); err != nil {
		return nil, err
	}
	return driver, nil
}