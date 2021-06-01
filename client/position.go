package client

type position struct {
	Market string  `json:"market"`
	Size   float64 `json:"size"`
	Side   string  `json:"side"`
}
