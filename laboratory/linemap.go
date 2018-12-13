package laboratory

type lineMapManager struct {
}

func NewLineMapManager() *lineMapManager {
	return &lineMapManager{}
}

func (lmm *lineMapManager) RegisterLineMap(elementImportPath string, lineMap map[int]int) {
}
