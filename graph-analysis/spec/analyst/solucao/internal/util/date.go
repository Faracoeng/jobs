package util

import (
	"time"
	"fmt"
)

var dateFormats = []string{
	"2006-01-02",           // padrão ISO
	"02/01/2006",           // formato brasileiro
	"2006-01-02 15:04:05",  // ISO com horário
	"02/01/2006 15:04:05",  // BR com horário
}

// ParseDate tenta converter uma string em time.Time com formatos conhecidos
func ParseDate(input string) (time.Time, error) {
	for _, layout := range dateFormats {
		if t, err := time.Parse(layout, input); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("formato de data inválido: %s", input)
}
