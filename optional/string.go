package optional

import "encoding/json"

type String struct {
	Set   bool
	Valid bool
	Value string
}

func (i *String) SetValue(v string) {
	i.Set = true
	i.Valid = true
	i.Value = v
}

func (i *String) UnmarshalJSON(data []byte) error {
	// If this method was called, the value was set.
	i.Set = true
	if string(data) == "null" {
		// The key was set to null
		i.Valid = false
		return nil
	}
	// The key isn't set to null
	var temp string
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	i.Value = temp
	i.Valid = true
	return nil
}
