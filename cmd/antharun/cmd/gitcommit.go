// gitcommit
package cmd

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// assumes GOPATH is in home directory if not set as environment variable
func gopath() string {

	// if gopath set return gopath
	if p := os.Getenv("GOPATH"); len(p) != 0 {
		return p
	}
	// if not set assume under user's home directory
	u, err := user.Current()
	if err != nil {
		return ""
	}

	return filepath.Join(u.HomeDir, "go/src")
}

func GitCommit(path string) (string, error) {
	cmdName := "git"
	cmdArgs := []string{"rev-parse", "HEAD"}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = path
	commitID, err := cmd.Output()
	return strings.TrimSpace(string(commitID)), err
}

var reporttemplate string = `
## Aim:



##Status
 


##Next steps:



##Execution instructions:


#### Get required repos

1. branch of antha-lang/antha :

{{.TripleQuote}}bash
cd $GOPATH/src/github.com/antha-lang/antha

git fetch 
git checkout ***ANTHACOMMIT****
cd -
{{.TripleQuote}}


2.  branch of antha-lang/elements


{{.TripleQuote}}bash
cd $GOPATH/src/github.com/antha-lang/elements
git fetch
git checkout ***ELEMENTSCOMMIT****
cd -
{{.TripleQuote}}

3. Other Dependencies:

{{.TripleQuote}}bash
***OTHERDEPENDENCIES***
{{.TripleQuote}}

4. (A) Pipetmaxdriver

{{.TripleQuote}}bash
cd $GOPATH/src/github.com/Synthace/PipetMaxDriver
git fetch
git checkout ***PIPETMAXDRIVERCOMMIT****
{{.TripleQuote}}

Or

4. (B) CybioDriver

{{.TripleQuote}}bash
cd $GOPATH/src/github.com/Synthace/CybioXMLDriver
git fetch
git checkout ***CYBIODRIVERCOMMIT****
cd -
{{.TripleQuote}}

#### Run whenever any source code is changed  (e.g. plate definitions, antha element changes, liquid class changes)

5. Build 

{{.TripleQuote}}bash
make current -C $GOPATH/src/github.com/antha-lang/elements
{{.TripleQuote}}

or

{{.TripleQuote}}
anthabuild
{{.TripleQuote}}


#### Run when parameters or workflow is changed

5. run


PipetMax:


{{.TripleQuote}}bash
antharun --driver  go://github.com/Synthace/PipetMaxDriver/server
{{.TripleQuote}}


Cybio:


{{.TripleQuote}}bash
cd $GOPATH/src/github.com/Synthace/CybioXMLDriver/server
go build ./...
./server -machine felix
{{.TripleQuote}}


{{.TripleQuote}}bash
antharun --driver localhost:50051 --inputPlateType pcrplate_skirted
{{.TripleQuote}}

6. Rename output file

e.g.

{{.TripleQuote}}bash
mv generated.sqlite pipetmaxday1.sqlite
{{.TripleQuote}}

or 

{{.TripleQuote}}bash
mv cybio.xml felixday1.xml
{{.TripleQuote}}



`
