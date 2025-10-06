package spec

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type HM int

func (t *HM) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("time must be a string: %w", err)
	}
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid time %q (want HH:MM)", s)
	}
	h, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || h < 0 || h > 23 || m < 0 || m > 59 {
		return fmt.Errorf("invalid time %q (0<=HH<=23, 0<=MM<=59)", s)
	}
	*t = HM(h*60 + m)
	return nil
}

func (t HM) MarshalJSON() ([]byte, error) {
	min := int(t)
	return json.Marshal(fmt.Sprintf("%02d:%02d", min/60, min%60))
}
