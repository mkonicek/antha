package pretty

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/meta"
)

type saver struct {
	files []wtype.File
}

func (a *saver) saveFiles(obj interface{}) ([]byte, error) {
	switch obj := obj.(type) {
	case wtype.File:
		if len(obj.Name) == 0 {
			break
		}
		a.files = append(a.files, obj)
	}
	return json.Marshal(obj)
}

// SaveFiles writes out any files in the execute.Result
func SaveFiles(out io.Writer, result *execute.Result) error {
	var s saver
	m := &meta.Marshaler{
		Struct: s.saveFiles,
	}

	for _, output := range result.Workflow.Outputs {
		// Just marshal for the side-effect
		if _, err := m.Marshal(output); err != nil {
			return err
		}
	}

	for _, file := range s.files {
		bs, err := file.ReadAll()
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(file.Name, bs, 0666); err != nil {
			return err
		}
	}

	return nil
}
