package source

import (
    "github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/etl/util"
)

type CSVSource struct {
    Path string
}

func (c *CSVSource) ReadAll() ([][]string, error) {
    return util.ReadCSVFile(c.Path)
}
