package wutil

import (
	"encoding/json"
	"fmt"
)

type Func1Prm interface {
	F(x float64) float64
	Name() string
}
type InvertibleFunc1Prm interface {
	Func1Prm
	I(x float64) float64
}

func UnmarshalFunc(b []byte) (Func1Prm, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal(b, &m)

	if err != nil {
		return nil, err
	}

	if _, ok := m["Quadratic"]; ok {
		var q Quadratic
		err = json.Unmarshal(b, &q)
		if err != nil {
			return nil, err
		}
		return Func1Prm(&q), nil
	} else if _, ok := m["Cubic"]; ok {
		var c Cubic
		err = json.Unmarshal(b, &c)
		if err != nil {
			return nil, err
		}
		return Func1Prm(&c), nil
	} else if _, ok := m["Quartic"]; ok {
		var q Quartic
		err = json.Unmarshal(b, &q)
		if err != nil {
			return nil, err
		}
		return Func1Prm(&q), nil
	}

	return nil, fmt.Errorf("Not a wutil function")
}
