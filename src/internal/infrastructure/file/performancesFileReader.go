package file

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sohosai/ultradonguri-server/internal/domain/entities"
)

func GetPerformances() ([]entities.PerformanceForPerformances, error) {
	//ファイル名は後々変更
	file, err := os.Open("events.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events PerformancesRecords
	dec := json.NewDecoder(file)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&events); err != nil {
		return nil, err
	}

	// 余剰トークンチェック
	if dec.More() {
		return nil, fmt.Errorf("trailing data after JSON array")
	}

	performanceEntities := events.ToDomainPerformanceForPerformances()
	return performanceEntities, nil
}
