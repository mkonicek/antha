package wtype

import "encoding/json"

// Representation of a file.
type File struct {
	path     string
	isOutput bool
	// Name is entirely for the purpose of displaying a meaningful name
	Name string
}

func NewFile(path string) *File {
	return &File{
		path:     path,
		isOutput: true,
	}
}

func (f File) Path() string {
	return f.path
}

func (f File) IsOutput() bool {
	return f.isOutput
}

func (f File) AsInput() *File {
	return &File{
		path: f.path,
		Name: f.Name,
	}
}

type fileJSON struct {
	Path     string `json:"Path"`
	IsOutput bool   `json:"IsOutput"`
	Name     string `json:"Name"`
}

func (f *File) MarshalJSON() ([]byte, error) {
	fj := &fileJSON{
		Path:     f.path,
		IsOutput: f.isOutput,
		Name:     f.Name,
	}
	return json.Marshal(fj)
}

func (f *File) UnmarshalJSON(bs []byte) error {
	fj := fileJSON{}
	if err := json.Unmarshal(bs, &fj); err != nil {
		return err
	} else {
		f.path = fj.Path
		f.isOutput = fj.IsOutput
		f.Name = fj.Name
		return nil
	}
}
