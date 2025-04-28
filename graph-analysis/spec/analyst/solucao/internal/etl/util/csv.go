package util

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
)

// ReadCSVFile lê o conteúdo de um CSV, ignora o cabeçalho e remove linhas duplicadas.
func ReadCSVFile(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo %s: %w", path, err)
	}
	defer file.Close()

	r := csv.NewReader(bufio.NewReader(file))
	lines, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("erro ao ler CSV %s: %w", path, err)
	}

	if len(lines) <= 1 {
		return [][]string{}, nil
	}

	seen := make(map[string]bool)
	var deduped [][]string
	for _, line := range lines[1:] {
		key := fmt.Sprint(line)
		if !seen[key] {
			seen[key] = true
			deduped = append(deduped, line)
		}
	}

	return deduped, nil
}
