package optional

import (
	"encoding/json"
	"time"
)

type Time struct {
	Set   bool
	Valid bool
	Value time.Time
}

func (i *Time) SetValue(v time.Time) {
	i.Set = true
	i.Valid = true
	i.Value = v
}

func (i *Time) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	i.Set = true
	if string(data) == "null" {
		// The key was set to null
		i.Valid = false
		return nil
	}
	// The key isn't set to null
	var temp time.Time
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.Value = temp
	i.Valid = true
	return nil
}
