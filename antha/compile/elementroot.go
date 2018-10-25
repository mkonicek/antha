package compile

import (
	"bytes"
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

type protocolDir struct {
	ProtocolName string
	Dir          string
}

// An ElementRoot collects data from multiple Element passes
type ElementRoot struct {
	outputPackageBase string
	// Map from protocol name to directory path
	protocolDirs []protocolDir
}

// NewElementRoot creates a new ElementRoot
func NewElementRoot(basePackage string) *ElementRoot {
	return &ElementRoot{
		outputPackageBase: basePackage,
	}
}

func (r *ElementRoot) addProtocolDirectory(protocolName, dir string) {
	r.protocolDirs = append(r.protocolDirs, protocolDir{
		ProtocolName: protocolName,
		Dir:          dir,
	})
}

// copyGoFiles returns the files to copy to output directory
func (r *ElementRoot) copyGoFiles() ([]*ElementFile, error) {
	seen := make(map[string]bool)
	var files []*ElementFile
	for _, pdir := range r.protocolDirs {
		if seen[pdir.Dir] {
			continue
		}
		seen[pdir.Dir] = true

		fis, err := ioutil.ReadDir(pdir.Dir)
		if err != nil {
			return nil, err
		}
		for _, fi := range fis {
			if fi.IsDir() {
				continue
			}
			if !strings.HasSuffix(fi.Name(), ".go") {
				continue
			}

			bs, err := ioutil.ReadFile(filepath.Join(pdir.Dir, fi.Name()))
			if err != nil {
				return nil, err
			}
			filename := path.Join(pdir.ProtocolName, fi.Name())
			files = append(files, &ElementFile{
				Name: filename,
				Data: bs,
			})
		}
	}

	return files, nil
}

func (r *ElementRoot) generateLib() ([]byte, error) {
	const tmpl = `
package _lib

import (
	"github.com/antha-lang/antha/component"
	"fmt"
	{{ range .Packages }}{{ .Name }} {{ .Path }}
	{{ end }}
)

func GetComponents() (comps []component.Component, err error) {
	defer func() {
		if res := recover(); res != nil {
			err = fmt.Errorf("cannot update component: %s", res)
		}
	}()

	add := func(da []*component.Component) {
		for _, d := range da {
			if err := component.UpdateParamTypes(d); err != nil {
				panic(err)
			}
			comps = append(comps, *d)
		}
	}
	{{ range .Packages }}add({{ .Name }}.GetComponent())
	{{ end }}
	return
}
`
	type Package struct {
		Name string
		Path string
	}

	type TVars struct {
		Packages []Package
	}
	tv := TVars{}

	for _, pdir := range r.protocolDirs {
		pkg := path.Join(r.outputPackageBase, pdir.ProtocolName, elementPackage)
		pkgName := manglePackageName(pkg)
		tv.Packages = append(tv.Packages, Package{
			Name: pkgName,
			Path: strconv.Quote(pkg),
		})
	}

	var out bytes.Buffer
	if err := template.Must(template.New("").Parse(tmpl)).Execute(&out, tv); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// Generate generates additional files not stric
func (r *ElementRoot) Generate() ([]*ElementFile, error) {
	files, err := r.copyGoFiles()
	if err != nil {
		return nil, err
	}

	libBs, err := r.generateLib()
	if err != nil {
		return nil, err
	}

	files = append(files, &ElementFile{
		Name: "_lib/lib.go",
		Data: libBs,
	})

	return files, nil
}
