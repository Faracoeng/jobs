
package reader_test

import (
    "os"
    "testing"
    "github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/reader"
    "github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/util"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func createTempCSV(t *testing.T, content string) string {
    tmpFile, err := os.CreateTemp("", "*.csv")
    require.NoError(t, err)
    _, err = tmpFile.WriteString(content)
    require.NoError(t, err)
    err = tmpFile.Close()
    require.NoError(t, err)
    return tmpFile.Name()
}

func TestReadCountries(t *testing.T) {
    content := "iso3,name\nBRA,Brazil\nUSA,United States"
    file := createTempCSV(t, content)
    defer os.Remove(file)

    result := reader.ReadCountries(file)
    require.NotNil(t, result)
    require.Len(t, result, 2)

    assert.Equal(t, "BRA", result[0].ISO3)
    assert.Equal(t, "Brazil", result[0].Name)
    assert.Equal(t, "USA", result[1].ISO3)
    assert.Equal(t, "United States", result[1].Name)
}

func TestReadVaccines(t *testing.T) {
    content := "name\nPfizer\nModerna"
    file := createTempCSV(t, content)
    defer os.Remove(file)

    result := reader.ReadVaccines(file)
    require.NotNil(t, result)
    require.Len(t, result, 2)

    assert.Equal(t, "Pfizer", result[0].Name)
    assert.Equal(t, "Moderna", result[1].Name)
}

func TestReadCovidCases(t *testing.T) {
    content := "iso3,date,total_cases,total_deaths\nBRA,2021-01-15,1000,50"
    file := createTempCSV(t, content)
    defer os.Remove(file)

    cases := reader.ReadCovidCases(file)
    require.Len(t, cases, 1)

    expectedDate, _ := util.ParseDate("2021-01-15")
    assert.Equal(t, "BRA", cases[0].ISO3)
    assert.Equal(t, 1000, cases[0].TotalCases)
    assert.Equal(t, 50, cases[0].TotalDeaths)
    assert.Equal(t, expectedDate, cases[0].Date)
}

func TestReadVaccinations(t *testing.T) {
    content := "iso3,date,total_vaccinated\nBRA,2021-01-22,15000"
    file := createTempCSV(t, content)
    defer os.Remove(file)

    records := reader.ReadVaccinations(file)
    require.Len(t, records, 1)

    expectedDate, _ := util.ParseDate("2021-01-22")
    assert.Equal(t, "BRA", records[0].ISO3)
    assert.Equal(t, 15000, records[0].TotalVaccinated)
    assert.Equal(t, expectedDate, records[0].Date)
}

func TestReadVaccineApprovals(t *testing.T) {
    content := "vaccine,date\nPfizer,2020-12-10"
    file := createTempCSV(t, content)
    defer os.Remove(file)

    approvals := reader.ReadVaccineApprovals(file)
    require.Len(t, approvals, 1)

    expectedDate, _ := util.ParseDate("2020-12-10")
    assert.Equal(t, "Pfizer", approvals[0].VaccineName)
    assert.Equal(t, expectedDate, approvals[0].Date)
}

func TestReadCountryVaccines(t *testing.T) {
    content := "iso3,vaccine\nBRA,Pfizer"
    file := createTempCSV(t, content)
    defer os.Remove(file)

    cv := reader.ReadCountryVaccines(file)
    require.Len(t, cv, 1)

    assert.Equal(t, "BRA", cv[0].ISO3)
    assert.Equal(t, "Pfizer", cv[0].VaccineName)
}