package optional

import "encoding/json"

type Bool struct {
	Set   bool
	Value bool
	Valid bool
}

func (i *Bool) SetValue(v bool) {
	i.Set = true
	i.Valid = true
	i.Value = v
}

func (i *Bool) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	i.Set = true
	if string(data) == "null" {
		// The key was set to null
		i.Valid = false
		return nil
	}
	// The key isn't set to null
	var temp bool
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.Value = temp
	i.Valid = true
	return nil
}
