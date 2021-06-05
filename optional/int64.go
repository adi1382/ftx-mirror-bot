package optional

import "encoding/json"

type Int64 struct {
	Set   bool
	Valid bool
	Value int64
}

func (i *Int64) SetValue(v int64) {
	i.Set = true
	i.Valid = true
	i.Value = v
}

func (i *Int64) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	i.Set = true
	if string(data) == "null" {
		// The key was set to null
		i.Valid = false
		return nil
	}
	// The key isn't set to null
	var temp int64
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.Value = temp
	i.Valid = true
	return nil
}
