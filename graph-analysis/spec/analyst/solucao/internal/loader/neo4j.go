package loader

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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
