package laboratory

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/logger"
)

type FileManager struct {
	lock         sync.Mutex
	outDir       string
	contents     map[string][]byte
	writtenCount uint64
	writtenSet   map[*wtype.File]struct{}
}

func NewFileManager(outDir string) *FileManager {
	return &FileManager{
		outDir:     outDir,
		contents:   make(map[string][]byte),
		writtenSet: make(map[*wtype.File]struct{}),
	}
}

func (fm *FileManager) ReadAll(f *wtype.File) ([]byte, error) {
	if f == nil {
		return nil, errors.New("Cannot read nil file")
	}
	fm.lock.Lock()
	defer fm.lock.Unlock()

	path := f.Path()
	bs, found := fm.contents[path]
	if !found {
		if content, err := ioutil.ReadFile(filepath.FromSlash(path)); err != nil {
			return nil, err
		} else {
			bs = content
			fm.contents[path] = content
		}
	}
	bsCopy := make([]byte, len(bs))
	copy(bsCopy, bs)
	return bsCopy, nil
}

func (fm *FileManager) WithWriter(fun func(io.Writer) error, fileName string) (*wtype.File, error) {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	fm.writtenCount++
	p := filepath.Join(fm.outDir, fmt.Sprintf("%d", fm.writtenCount))
	if fh, err := os.OpenFile(p, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0400); err != nil {
		return nil, err
	} else {
		err := fun(fh)
		fh.Close()
		if err != nil {
			return nil, err
		} else {
			f := wtype.NewFile(p)
			f.Name = fileName
			fm.writtenSet[f] = struct{}{}
			return f, nil
		}
	}
}

func (fm *FileManager) WriteAll(bs []byte, fileName string) (*wtype.File, error) {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	fm.writtenCount++
	p := filepath.Join(fm.outDir, fmt.Sprintf("%d", fm.writtenCount))
	if err := ioutil.WriteFile(p, bs, 0400); err != nil {
		return nil, err
	} else {
		p2 := filepath.ToSlash(p)
		fm.contents[p2] = bs
		f := wtype.NewFile(p)
		f.Name = fileName
		fm.writtenSet[f] = struct{}{}
		return f, nil
	}
}

func (fm *FileManager) WriteString(str string, fileName string) (*wtype.File, error) {
	return fm.WriteAll([]byte(str), fileName)
}

func (fm *FileManager) SummarizeWritten(logger *logger.Logger) {
	logger = logger.With("fileManager", "summary")
	fm.lock.Lock()
	defer fm.lock.Unlock()

	for f := range fm.writtenSet {
		name := "<unnamed>"
		if f.Name != "" {
			name = f.Name
		}
		logger.Log("name", name, "path", f.Path())
	}
}
