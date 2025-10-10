package file

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
)

func GetPerformances() ([]entities.Performance, error) {
	file, err := os.Open("events.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []entities.Performance
	dec := json.NewDecoder(file)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&events); err != nil {
		return nil, err
	}

	// 余剰トークンチェック
	if dec.More() {
		return nil, fmt.Errorf("trailing data after JSON array")
	}

	return events, nil
}
