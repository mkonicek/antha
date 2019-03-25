package compile

import (
	"encoding/json"

	"github.com/antha-lang/antha/antha/token"
)

type Meta struct {
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Defaults    map[string]json.RawMessage `json:"defaults"`
	Tags        []string                   `json:"tags"`
	Ports       map[token.Token][]*Field   `json:"ports"`
}

func (a *Antha) meta(bs []byte) error {
	meta := &Meta{
		Ports: make(map[token.Token][]*Field),
	}

	if len(bs) != 0 {
		if err := json.Unmarshal(bs, meta); err != nil {
			return err
		}
	}

	meta.Name = a.protocolName
	meta.Description = a.description

	for _, msg := range a.messages {
		meta.Ports[msg.Kind] = msg.Fields
	}

	a.Meta = meta
	return nil
}
