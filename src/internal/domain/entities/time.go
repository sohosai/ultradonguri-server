package entities

import (
	"encoding/json"
	"fmt"
	"time"
)

type ISOTime struct {
	time.Time
}

func (t *ISOTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) < 2 {
		return fmt.Errorf("invalid time format: %s", s)
	}
	s = s[1 : len(s)-1] // remove quotes

	// ISO 8601 (RFC3339) 形式でパース
	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}
	t.Time = parsed
	return nil
}

func (t ISOTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Format(time.RFC3339))
}

func (t ISOTime) String() string {
	return t.Format(time.RFC3339)
}
