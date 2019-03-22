package compile

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/antha-lang/antha/antha/token"
)

type Meta struct {
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Defaults    map[string]json.RawMessage `json:"defaults"`
	Tags        []string                   `json:"tags"`
	Ports       map[token.Token][]*Field   `json:"ports"`
}

func (a *Antha) Meta() (*Meta, error) {
	meta := &Meta{
		Ports: make(map[token.Token][]*Field),
	}

	mdPath := filepath.Join(filepath.Dir(a.fileSet.File(a.file.Package).Name()), "metadata.json")
	if bs, err := ioutil.ReadFile(mdPath); err != nil && !os.IsNotExist(err) {
		return nil, err
	} else if err == nil {
		if err := json.Unmarshal(bs, meta); err != nil {
			return nil, err
		}
	}

	meta.Name = a.protocolName
	meta.Description = a.description

	for _, msg := range a.messages {
		meta.Ports[msg.Kind] = msg.Fields
	}

	return meta, nil
}
