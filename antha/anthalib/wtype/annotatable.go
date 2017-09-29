package wtype

type Annotatable interface {
	GetData(k string) (AnthaData, error)
	AddData(k string, d AnthaData) error
	ClearData(k string) error
}

type AnthaData struct {
	Name string
	Data string
}
