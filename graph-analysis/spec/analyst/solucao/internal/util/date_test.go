package util_test

import (
    "testing"
    "time"

    "github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/util"
    "github.com/stretchr/testify/require"
)

func TestParseDate(t *testing.T) {
    input := "2021-03-01"
    expected := time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC)

    result, err := util.ParseDate(input)

    require.NoError(t, err)
    require.Equal(t, expected, result)
}

func TestParseDateInvalid(t *testing.T) {
    _, err := util.ParseDate("not-a-date")
    require.Error(t, err)
}
