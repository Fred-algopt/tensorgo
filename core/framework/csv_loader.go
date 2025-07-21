package dataset

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

// LoadCSV charge un CSV en slice de map[string]interface{}
func LoadCSV(filepath string) (*Dataset[map[string]interface{}], error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("erreur ouverture CSV: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("lecture headers: %v", err)
	}

	var data []map[string]interface{}
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("erreur lecture ligne: %v", err)
		}

		row := make(map[string]interface{})
		for i, val := range record {
			// Tentative de conversion automatique en int ou float
			if i < len(headers) {
				if f, err := strconv.ParseFloat(val, 64); err == nil {
					row[headers[i]] = f
				} else {
					row[headers[i]] = val
				}
			}
		}
		data = append(data, row)
	}
	return New(data), nil
}