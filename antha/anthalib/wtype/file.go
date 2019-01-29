package wtype

import "encoding/json"

// Representation of a file.
type File struct {
	path string
	// Name is entirely for the purpose of displaying a meaningful name
	Name string
}

func NewFile(path string) *File {
	return &File{
		path: path,
	}
}

func (f File) Path() string {
	return f.path
}

type fileJSON struct {
	Path string `json:"Path"`
	Name string `json:"Name"`
}

func (f *File) MarshalJSON() ([]byte, error) {
	fj := &fileJSON{
		Path: f.path,
		Name: f.Name,
	}
	return json.Marshal(fj)
}

func (f *File) UnmarshalJSON(bs []byte) error {
	fj := fileJSON{}
	if err := json.Unmarshal(bs, &fj); err != nil {
		return err
	} else {
		f.path = fj.Path
		f.Name = fj.Name
		return nil
	}
}
