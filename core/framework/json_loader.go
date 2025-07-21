
package dataset

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadJSON charge un fichier JSON contenant un tableau d'objets.
func LoadJSON(filepath string) (*Dataset[map[string]interface{}], error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture fichier JSON : %v", err)
	}

	var records []map[string]interface{}
	err = json.Unmarshal(data, &records)
	if err != nil {
		return nil, fmt.Errorf("erreur parsing JSON : %v", err)
	}

	return New(records), nil
}