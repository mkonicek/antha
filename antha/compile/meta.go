package compile

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Meta struct {
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Defaults    map[string]json.RawMessage `json:"defaults"`
	Tags        []string                   `json:"tags"`
}

func (a *Antha) Meta() (*Meta, error) {
	meta := &Meta{}

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
	return meta, nil
}
