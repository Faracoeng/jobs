package reader

import (
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/etl/source"
)

// Line representa uma linha de dados lida do DataSource
type Line struct {
	Data []string
	Err  error
}

// StreamData lê os dados de um DataSource e envia as linhas através de um canal
func StreamData(ds source.DataSource) <-chan Line {
	out := make(chan Line)

	go func() {
		defer close(out)

		rows, err := ds.ReadAll()
		if err != nil {
			out <- Line{Err: err}
			return
		}

		for _, row := range rows {
			out <- Line{Data: row}
		}
	}()

	return out
}
