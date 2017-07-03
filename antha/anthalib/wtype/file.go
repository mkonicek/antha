package wtype

import (
	"bytes"

	api "github.com/antha-lang/antha/api/v1"
	"github.com/golang/protobuf/jsonpb"
)

// Representation of a file.
type File struct {
	Name  string
	bytes []byte
}

// ReadAll returns contents of file.
func (a *File) ReadAll() ([]byte, error) {
	return a.bytes, nil
}

// WriteAll replaces contents of file with data.
func (a *File) WriteAll(data []byte) error {
	a.bytes = data
	return nil
}

func (a *File) UnmarshalBlob(blob *api.Blob) error {
	a.Name = blob.Name
	a.bytes = blob.GetBytes().GetBytes()

	return nil
}

func (a *File) UnmarshalJSON(data []byte) error {
	var u jsonpb.Unmarshaler
	var blob api.Blob
	if err := u.Unmarshal(bytes.NewReader(data), &blob); err != nil {
		return err
	}

	return a.UnmarshalBlob(&blob)
}

func (a File) MarshalJSON() ([]byte, error) {
	blob := &api.Blob{
		Name: a.Name,
		From: &api.Blob_Bytes{
			Bytes: &api.FromBytes{
				Bytes: a.bytes,
			},
		},
	}

	var m jsonpb.Marshaler
	var out bytes.Buffer
	if err := m.Marshal(&out, blob); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
