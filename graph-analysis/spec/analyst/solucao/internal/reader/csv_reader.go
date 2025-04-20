package reader

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"github.com/Faracoeng/jobs/graph-analysis/spec/analyst/solucao/internal/model"
)
func ReadeCountries(path string) []model.Country{
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Erro ao abrir arquivo %s: %v", path, err)
		return nil
	}
	defer f.Close()

	// Adiciona buffer para leitura, bom para arquivos CSV grandes
	r := csv.NewReader(bufio.NewReader(f))
	countries, err := r.ReadAll()
	if err != nil {
		log.Printf("Erro ao ler CSV %s: %v", path, err)
		return nil
	}
	var countryList []model.Country
	for i, row := range countries {
		// Pula a primeira linha do csv
		if i == 0 {
			continue
		}
		if len(row) < 2 {
			log.Printf("Linha %d inválida em %s: %v", i, path, row)
			continue
		}
		country := model.Country{
			ISO3: row[0],
			Name: row[1],
		}
		countryList = append(countryList, country)
	}
	return countryList
}