package effects

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/logger"
)

type FileManager struct {
	lock         sync.Mutex
	inDir        string
	outDir       string
	contents     map[string][]byte
	writtenCount uint64
	writtenSet   map[*wtype.File]struct{}
}

func NewFileManager(inDir, outDir string) (*FileManager, error) {
	if inDir, err := filepath.Abs(inDir); err != nil {
		return nil, err
	} else if outDir, err := filepath.Abs(outDir); err != nil {
		return nil, err
	} else {
		return &FileManager{
			inDir:      inDir,
			outDir:     outDir,
			contents:   make(map[string][]byte),
			writtenSet: make(map[*wtype.File]struct{}),
		}, nil
	}
}

func (fm *FileManager) ReadAll(f *wtype.File) ([]byte, error) {
	if f == nil {
		return nil, errors.New("Cannot read nil file")
	}
	fm.lock.Lock()
	defer fm.lock.Unlock()

	p := f.Path()
	if bs, found := fm.contents[p]; found {
		return copyBytes(bs), nil
	}

	pLocal := p

	if u, err := url.Parse(p); err == nil {
		switch u.Scheme {
		case "http", "https":
			f, err := fm.WithWriter(func(w io.Writer) error {
				if resp, err := http.Get(p); err != nil {
					return err
				} else {
					defer resp.Body.Close()
					_, err := io.Copy(w, resp.Body)
					return err
				}
			}, "")
			if err != nil {
				return nil, err
			} else {
				// guaranteed to exist (see WithWriter):
				bs := fm.contents[f.Path()]
				fm.contents[p] = bs
				return copyBytes(bs), nil
			}

		case "file":
			pLocal = filepath.Join(u.Host, filepath.FromSlash(u.Path))
		case "":
			pLocal = filepath.FromSlash(pLocal)
		default:
			return nil, fmt.Errorf("Unrecognised url scheme: %v", u.Scheme)
		}
	}

	// either we couldn't parse it as a url, or we could and it was
	// file:// or we could and there was an empty scheme.

	// we need to be a bit careful here to make sure there is no way to escape from inDir
	pLocal = filepath.Clean(pLocal)
	if !filepath.IsAbs(pLocal) {
		if pLocalAbs, err := filepath.Abs(filepath.Join(fm.inDir, pLocal)); err != nil {
			return nil, err
		} else {
			pLocal = pLocalAbs
		}
	}
	// The use of Clean (Join also calls Clean) should mean that we
	// don't have any .. in the path, so if we have inDir as a prefix,
	// we should be able to be confident that we really are accessing a
	// file within inDir.
	if !strings.HasPrefix(pLocal, fm.inDir+string(os.PathSeparator)) {
		return nil, fmt.Errorf("Invalid path: %v", p)
	}
	if bs, err := ioutil.ReadFile(pLocal); err != nil {
		return nil, err
	} else {
		fm.contents[p] = bs
		return copyBytes(bs), nil
	}
}

func (fm *FileManager) WithWriter(fun func(io.Writer) error, fileName string) (*wtype.File, error) {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	fm.writtenCount++
	pLocal := filepath.Join(fm.outDir, fmt.Sprintf("%d", fm.writtenCount))
	if fh, err := os.OpenFile(pLocal, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0400); err != nil {
		return nil, err
	} else {
		buf := new(bytes.Buffer)
		mfh := io.MultiWriter(fh, buf)
		err := fun(mfh)
		if err != nil {
			return nil, err
		} else if err := fh.Sync(); err != nil {
			return nil, err
		} else if err := fh.Close(); err != nil {
			return nil, err
		} else {
			fm.contents[pLocal] = buf.Bytes()
			f := wtype.NewFile(pLocal)
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
	pLocal := filepath.Join(fm.outDir, fmt.Sprintf("%d", fm.writtenCount))
	if err := ioutil.WriteFile(pLocal, bs, 0400); err != nil {
		return nil, err
	} else {
		bsCopy := make([]byte, len(bs))
		copy(bsCopy, bs)
		fm.contents[pLocal] = bsCopy
		f := wtype.NewFile(pLocal)
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

func copyBytes(bs []byte) []byte {
	if bs != nil {
		bsCopy := make([]byte, len(bs))
		copy(bsCopy, bs)
		return bsCopy
	} else {
		return nil
	}
}
