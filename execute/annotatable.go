package execute

// An Annotatable is an object that can data attached to it
type Annotatable interface {
	GetData(key string) ([]byte, error)
	SetData(key string, data []byte) error
}
