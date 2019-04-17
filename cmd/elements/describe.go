package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/golang/protobuf/proto"

	"github.com/antha-lang/antha/protobuf"

	element "github.com/Synthace/microservice/cmd/element/protobuf"

	"github.com/antha-lang/antha/antha/compile"
	"github.com/antha-lang/antha/antha/token"
	"github.com/antha-lang/antha/composer"
	"github.com/antha-lang/antha/logger"
	"github.com/antha-lang/antha/workflow"
)

type elementWithMeta struct {
	repoName      workflow.RepositoryName
	repo          *workflow.Repository
	anthaFilePath string
	element       []byte
	meta          []byte
}

const (
	outputFormatHuman    = "human"
	outputFormatJSON     = "json"
	outputFormatProtobuf = "protobuf"
)

func describe(l *logger.Logger, args []string) error {
	flagSet := flag.NewFlagSet(flag.CommandLine.Name()+" describe", flag.ContinueOnError)
	flagSet.Usage = workflow.NewFlagUsage(flagSet, "Show descriptions of elements")

	validOutputFormats := []string{
		outputFormatHuman,
		outputFormatJSON,
		outputFormatProtobuf,
	}

	var regexStr, inDir, outputFormat string
	flagSet.StringVar(&regexStr, "regex", "", "Regular expression to match against element type path (optional)")
	flagSet.StringVar(&inDir, "indir", "", "Directory from which to read files (optional)")
	flagSet.StringVar(&outputFormat, "format", "human", fmt.Sprintf("Format to output data in. One of %v", validOutputFormats))

	if err := flagSet.Parse(args); err != nil {
		return err
	}

	isValidOutputFormat := false
	for _, f := range validOutputFormats {
		if outputFormat == f {
			isValidOutputFormat = true
			break
		}
	}
	if !isValidOutputFormat {
		return fmt.Errorf("'%v' is not a valid output format. Use one of %v", outputFormat, validOutputFormats)
	}

	if wfPaths, err := workflow.GatherPaths(flagSet, inDir); err != nil {
		return err
	} else if rs, err := workflow.ReadersFromPaths(wfPaths); err != nil {
		return err
	} else if wf, err := workflow.WorkflowFromReaders(rs...); err != nil {
		return err
	} else if regex, err := regexp.Compile(regexStr); err != nil {
		return err
	} else {
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
						repo:     repo,
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

		// To stay compatible with the previous code that generated the
		// elements.pb file, we need to serialise everything in one go, as an
		// array of element.Element instances inside a protobuf.Elements
		// wrapper, which means we have to store all elements in memory until
		// we're done. Urgh. Since we're doing that for the protobuf code, we
		// may as well do the same thing for JSON. But ideally we'll switch the
		// elements microservice to accept a JSON payload, kill the
		// outputFormat=protobuf code paths here, then we can write the JSON to
		// STDOUT one element at a time.
		var pbElements []*element.Element

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
			antha, err := tet.EnsureTranspiler(ewm.anthaFilePath, ewm.element, ewm.meta)
			if err != nil {
				return err
			}
			switch outputFormat {
			case outputFormatHuman:
				{
					if err := printHumanReadable(antha, et); err != nil {
						return err
					}

				}
			case outputFormatJSON, outputFormatProtobuf:
				{
					elem, err := getProtobufElement(antha, et, ewm)
					if err != nil {
						return err
					}
					pbElements = append(pbElements, elem)
				}
			}
		}

		w := bufio.NewWriter(os.Stdout)
		defer w.Flush()

		var bs []byte
		var err error

		switch outputFormat {
		case outputFormatProtobuf:
			{
				bs, err = proto.Marshal(&protobuf.Elements{
					Elements: pbElements,
				})
			}
		case outputFormatJSON:
			{
				bs, err = json.Marshal(pbElements)
			}
		default:
			{
				return nil
			}
		}

		if err != nil {
			return err
		}
		_, err = w.Write(bs)
		if err != nil {
			return err
		}

		return nil
	}
}

func protobufPorts(fields []*compile.Field, kind string) ([]*element.Port, error) {
	var result []*element.Port
	for _, field := range fields {
		typeString, err := field.TypeString()
		if err != nil {
			return nil, err
		}
		port := &element.Port{
			Name:        field.Name,
			Type:        typeString,
			Description: field.Meta.Description,
			Kind:        kind,
		}
		result = append(result, port)
	}
	return result, nil
}

func getProtobufElement(antha *compile.Antha, et *workflow.ElementType, ewm *elementWithMeta) (*element.Element, error) {
	e := &element.Element{
		Name:        string(et.Name()),
		Package:     string(ewm.repoName) + "/" + path.Dir(ewm.anthaFilePath),
		Description: antha.Meta.Description,
		Tags:        antha.Meta.Tags,
		Version: &element.Element_GitVersion{
			GitVersion: &element.GitVersion{
				// The repoName field should always be the full URL of the repo
				// minus its scheme. See discussion here:
				// https://synthace.slack.com/archives/CGP8FDL9Z/p1554902911103300
				RepoUrl: "https://" + string(ewm.repoName),
				Sha:     ewm.repo.Commit,
			},
		},
		RequiresAnthaCore: true,
	}

	// Build in ports
	inputs, err := protobufPorts(antha.Meta.Ports[token.INPUTS], "Inputs")
	if err != nil {
		return nil, err
	}
	parameters, err := protobufPorts(antha.Meta.Ports[token.PARAMETERS], "Parameters")
	if err != nil {
		return nil, err
	}
	e.InPorts = append(inputs, parameters...)

	// Build out ports
	outputs, err := protobufPorts(antha.Meta.Ports[token.OUTPUTS], "Outputs")
	if err != nil {
		return nil, err
	}
	data, err := protobufPorts(antha.Meta.Ports[token.DATA], "Data")
	if err != nil {
		return nil, err
	}
	e.OutPorts = append(outputs, data...)

	e.Body = &element.Output{
		Body:     ewm.element,
		Complete: true,
	}

	// TODO: set these?
	// e.BuildOutput =
	// e.CreatedBy =

	return e, nil
}

func printHumanReadable(antha *compile.Antha, et *workflow.ElementType) error {
	const (
		indent  = "\t"
		indent2 = "\t\t"
		indent3 = "\t\t\t"

		fmtStr = `%v
%sRepositoryName: %v
%sElementPath: %v
%sTags: %v
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
			indent, strings.Join(meta.Tags, ", "),
			indent, desc,
			indent,
			indent2, inputs,
			indent2, outputs,
			indent2, params,
			indent2, data,
		)
	}
	return nil
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
