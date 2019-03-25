package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/antha-lang/antha/antha/compile"
	"github.com/antha-lang/antha/antha/token"
	"github.com/antha-lang/antha/composer"
	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow"
)

func describe(l *logger.Logger, args []string) error {
	const (
		indent  = "\t"
		indent2 = "\t\t"
		indent3 = "\t\t\t"

		fmtStr = `%v
%sRepositoryName: %v
%sElementPath: %v
%sDescription:
%v
%sPorts:
%sInputs:
%v
%sParameters:
%v
%sOutputs:
%v
%sData:
%v
`
	)

	flagSet := flag.NewFlagSet(flag.CommandLine.Name()+" describe", flag.ContinueOnError)
	flagSet.Usage = workflow.NewFlagUsage(flagSet, "Show descriptions of elements")

	var regexStr, inDir string
	flagSet.StringVar(&regexStr, "regex", "", "Regular expression to match against element type path (optional)")
	flagSet.StringVar(&inDir, "indir", "", "Directory from which to read files (optional)")

	if err := flagSet.Parse(args); err != nil {
		return err
	} else if wfPaths, err := workflow.GatherPaths(flagSet, inDir); err != nil {
		return err
	} else if rs, err := workflow.ReadersFromPaths(wfPaths); err != nil {
		return err
	} else if wf, err := workflow.WorkflowFromReaders(rs...); err != nil {
		return err
	} else if regex, err := regexp.Compile(regexStr); err != nil {
		return err
	} else {
		type elementWithMeta struct {
			repoName      workflow.RepositoryName
			anthaFilePath string
			element       []byte
			meta          []byte
		}
		// the map keys are the dir paths of the element so that it's the same for the antha file and the metadata
		elements := make(map[string]*elementWithMeta)
		elementNames := []string{}

		for repoName, repo := range wf.Repositories {
			err := repo.Walk(func(f *workflow.File) error {
				dir := filepath.Dir(f.Name)
				if (!workflow.IsAnthaFile(f.Name) && !workflow.IsAnthaMetadata(f.Name)) || !regex.MatchString(dir) {
					return nil
				}

				ewm, found := elements[dir]
				if !found {
					ewm = &elementWithMeta{
						repoName: repoName,
					}
					elements[dir] = ewm
					elementNames = append(elementNames, dir)
				}

				if rc, err := f.Contents(); err != nil {
					return err
				} else {
					defer rc.Close()
					if bs, err := ioutil.ReadAll(rc); err != nil {
						return err
					} else if workflow.IsAnthaFile(f.Name) {
						ewm.anthaFilePath = f.Name
						ewm.element = bs
					} else if workflow.IsAnthaMetadata(f.Name) {
						ewm.meta = bs
					}
					return nil
				}
			})
			if err != nil {
				return err
			}
		}

		sort.Strings(elementNames)
		for _, name := range elementNames {
			ewm := elements[name]
			if ewm.element == nil { // we cope with meta being nil
				continue
			}

			et := &workflow.ElementType{
				ElementPath:    workflow.ElementPath(filepath.ToSlash(filepath.Dir(ewm.anthaFilePath))),
				RepositoryName: ewm.repoName,
			}
			tet := composer.NewTranspilableElementType(et)
			if antha, err := tet.EnsureTranspiler(ewm.anthaFilePath, ewm.element, ewm.meta); err != nil {
				return err
			} else {
				meta := antha.Meta
				desc := indent2 + strings.Replace(strings.Trim(meta.Description, "\n"), "\n", "\n"+indent2, -1)
				if inputs, err := formatFields(meta.Defaults, meta.Ports[token.INPUTS], indent3, indent); err != nil {
					return err
				} else if outputs, err := formatFields(meta.Defaults, meta.Ports[token.OUTPUTS], indent3, indent); err != nil {
					return err
				} else if params, err := formatFields(meta.Defaults, meta.Ports[token.PARAMETERS], indent3, indent); err != nil {
					return err
				} else if data, err := formatFields(meta.Defaults, meta.Ports[token.DATA], indent3, indent); err != nil {
					return err
				} else {
					fmt.Printf(fmtStr,
						et.Name(),
						indent, et.RepositoryName,
						indent, et.ElementPath,
						indent, desc,
						indent,
						indent2, inputs,
						indent2, outputs,
						indent2, params,
						indent2, data,
					)
				}
			}
		}
		return nil
	}
}

func formatFields(defaults map[string]json.RawMessage, fields []*compile.Field, prefix, indent string) (string, error) {
	if len(fields) == 0 {
		return prefix + "None", nil
	}
	acc := make([]string, 0, 2*len(fields))
	for _, field := range fields {
		// If the type is an inline type declaration, this formatting
		// will go wrong. But life would be bad already if that sort of
		// thing was going on... so we just hope for the best.
		if typeStr, err := field.TypeString(); err != nil {
			return "", err
		} else {
			// the default can be a multiline thing, eg a map. So we have to be careful:
			def := ""
			if v, found := defaults[field.Name]; found {
				if bs, err := json.MarshalIndent(v, prefix+indent+indent, indent); err != nil {
					return "", err
				} else if bytes.ContainsRune(bs, '\n') {
					def = fmt.Sprintf("\n%s%sdefault:\n%s%s", prefix, indent, prefix+indent+indent, bs)
				} else {
					def = fmt.Sprintf(" (default: %s)", v)
				}
			}
			acc = append(acc, fmt.Sprintf("%s%s: %s%s", prefix, field.Name, typeStr, def))
			doc := strings.Trim(field.Doc, "\n")
			if len(doc) != 0 {
				acc = append(acc, prefix+indent+strings.Replace(doc, "\n", "\n"+prefix+indent, -1))
			}
		}
	}
	return strings.Join(acc, "\n"), nil
}
