package util_test

import (
	"os"
	"testing"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/etl/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadCSVFile_Deduplication(t *testing.T) {
	content := "col1,col2\nA,B\nA,B\nC,D"
	tmp := createTempCSV(t, content)
	defer os.Remove(tmp)

	lines, err := util.ReadCSVFile(tmp)
	require.NoError(t, err)
	require.Len(t, lines, 2) // sem cabeçalho e sem duplicação
	assert.Equal(t, []string{"A", "B"}, lines[0])
	assert.Equal(t, []string{"C", "D"}, lines[1])
}

func createTempCSV(t *testing.T, content string) string {
    tmpFile, err := os.CreateTemp("", "*.csv")
    require.NoError(t, err)
    _, err = tmpFile.WriteString(content)
    require.NoError(t, err)
    err = tmpFile.Close()
    require.NoError(t, err)
    return tmpFile.Name()
}
