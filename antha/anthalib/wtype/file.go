package wtype

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
