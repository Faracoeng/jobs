package util

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

// ReadCSVFile lê o conteúdo de um CSV, ignora o cabeçalho e remove linhas duplicadas
func ReadCSVFile(path string) [][]string {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Erro ao abrir arquivo %s: %v", path, err)
		return nil
	}
	defer file.Close()

	r := csv.NewReader(bufio.NewReader(file))
	lines, err := r.ReadAll()
	if err != nil {
		log.Printf("Erro ao ler CSV %s: %v", path, err)
		return nil
	}

	if len(lines) <= 1 {
		return [][]string{}
	}

	seen := make(map[string]bool)
	var deduped [][]string
	for _, line := range lines[1:] {
		// chave baseada na linha completa
		key := fmt.Sprint(line) 
		if !seen[key] {
			seen[key] = true
			deduped = append(deduped, line)
		}
	}
	return deduped
}
