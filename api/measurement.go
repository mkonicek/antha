package api

type Measurement struct {
	Value float64 `json:"value"`
	// SI unit
	Unit string `json:"unit"`
}
