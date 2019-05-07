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
	"github.com/antha-lang/antha/utils"
)

type FileManager struct {
	lock         sync.Mutex
	inDir        string
	outDir       string
	inCache      map[string][]byte
	outCache     map[string][]byte
	writtenCount uint64
}

func NewFileManager(inDir, outDir string) (*FileManager, error) {
	if inDir, err := filepath.Abs(inDir); err != nil {
		return nil, err
	} else if outDir, err := filepath.Abs(outDir); err != nil {
		return nil, err
	} else {
		return &FileManager{
			inDir:    inDir,
			outDir:   outDir,
			inCache:  make(map[string][]byte),
			outCache: make(map[string][]byte),
		}, nil
	}
}

func (fm *FileManager) ReadAll(f *wtype.File) ([]byte, error) {
	if f == nil || f.Path() == "" {
		return nil, errors.New("Cannot read nil file")
	}
	fm.lock.Lock()
	defer fm.lock.Unlock()

	p := f.Path()
	if f.IsOutput() {
		if bs, found := fm.outCache[p]; found {
			return copyBytes(bs), nil
		} else {
			return nil, fmt.Errorf("Attempt to read unknown but written file: %#v", f)
		}
	}

	// input only from here on
	if bs, found := fm.inCache[p]; found {
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
				// guaranteed to exist now (see WithWriter):
				bs := fm.outCache[f.Path()]
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
	pLocal = filepath.Join(fm.inDir, pLocal)

	// The use of Clean (Join calls Clean) should mean that we don't
	// have any .. in the path, so if we have inDir as a prefix, we are
	// confident that we really are accessing a file within inDir.
	if !strings.HasPrefix(pLocal, fm.inDir+string(os.PathSeparator)) {
		return nil, fmt.Errorf("Invalid file: %#v", f)
	}
	if bs, err := ioutil.ReadFile(pLocal); err != nil {
		return nil, err
	} else {
		fm.inCache[p] = bs
		return copyBytes(bs), nil
	}
}

func (fm *FileManager) WithWriter(fun func(io.Writer) error, fileName string) (*wtype.File, error) {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	fm.writtenCount++
	leaf := fmt.Sprintf("%d", fm.writtenCount)
	pLocal := filepath.Join(fm.outDir, leaf)
	if fh, err := utils.CreateFile(pLocal, utils.ReadWrite); err != nil {
		return nil, err
	} else {
		buf := new(bytes.Buffer)
		mfh := io.MultiWriter(fh, buf)
		if err := fun(mfh); err != nil {
			return nil, err
		} else if err := fh.Sync(); err != nil {
			return nil, err
		} else if err := fh.Close(); err != nil {
			return nil, err
		} else {
			fm.outCache[leaf] = buf.Bytes()
			f := wtype.NewFile(leaf)
			f.Name = fileName
			return f, nil
		}
	}
}

func (fm *FileManager) WriteAll(bs []byte, fileName string) (*wtype.File, error) {
	fm.lock.Lock()
	defer fm.lock.Unlock()

	fm.writtenCount++
	leaf := fmt.Sprintf("%d", fm.writtenCount)
	pLocal := filepath.Join(fm.outDir, leaf)
	if err := utils.CreateAndWriteFile(pLocal, bs, utils.ReadWrite); err != nil {
		return nil, err
	} else {
		fm.outCache[leaf] = copyBytes(bs)
		f := wtype.NewFile(leaf)
		f.Name = fileName
		return f, nil
	}
}

func (fm *FileManager) WriteString(str string, fileName string) (*wtype.File, error) {
	return fm.WriteAll([]byte(str), fileName)
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
