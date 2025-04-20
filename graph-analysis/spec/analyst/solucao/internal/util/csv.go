package util

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
)

// ReadCSVFile lê o conteúdo de um CSV e retorna todas as linhas, ignorando o cabeçalho
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
	// Retornar as linhas, ignorando o cabeçalho
	return lines[1:] 
}
