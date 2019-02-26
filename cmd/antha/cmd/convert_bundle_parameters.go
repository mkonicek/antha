package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/execute/executeutil"
	"github.com/antha-lang/antha/workflow"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var convertBundleCmd = &cobra.Command{
	Use:   "bundle-parameters",
	Short: "update bundle parameters according to the map of old parameter names to new",
	RunE:  convertBundle,
}

type NewElementMappingDetails struct {
	ConversionDetails `json:"new-element-mapping-details"`
}

type ConversionDetails struct {
	OldElementName       string                     `json:"old-element-name"`
	NewElementName       string                     `json:"new-element-name"`
	NewParameterMapping  map[string]string          `json:"new-parameter-mapping"`
	DeprecatedParameters []string                   `json:"deprecated-parameters"`
	NewParameters        map[string]json.RawMessage `json:"new-parameters"`
	ParameterTypeChanges map[string]json.RawMessage `json:"parameter-type-changes"`
}

func (c NewElementMappingDetails) Empty() bool {
	if c.ConversionDetails.OldElementName != "" || c.ConversionDetails.NewElementName != "" || len(c.ConversionDetails.NewParameterMapping) > 0 {
		return false
	}
	return true
}

// replaces old parameter names with new; this must be done after changing parameter names
func convertProcesses(in map[string]workflow.Process, newElementNames ConversionDetails) map[string]workflow.Process {
	ret := make(map[string]workflow.Process)

	for k, v := range in {
		var comp string

		if v.Component == newElementNames.OldElementName {
			comp = newElementNames.NewElementName
		} else {
			comp = v.Component
		}
		ret[k] = workflow.Process{
			Component: comp,
			Metadata:  v.Metadata,
		}
	}

	return ret
}

func containsSomethingToConvert(in executeutil.Bundle, newElementNames ConversionDetails) bool {
	for _, element := range in.Processes {
		if element.Component == newElementNames.OldElementName {
			return true
		}
	}
	return false
}

func convertParametersAndConnections(in executeutil.Bundle, newElementNames ConversionDetails) (map[string]map[string]json.RawMessage, []workflow.Connection, error) {
	var errs []string
	parameters := make(map[string]map[string]json.RawMessage)
	connections := in.Connections

	for processName, element := range in.Processes {

		// check if element name is to be replaced or is already equal to new element name
		// in order to suupport updating of parameter names.
		if element.Component == newElementNames.OldElementName || element.Component == newElementNames.NewElementName {
			// get existing parameters
			parametersForThisProcess, found := in.Parameters[processName]
			if !found {
				panic(fmt.Errorf("parameters not found for %s", processName))
			}
			newParameters := make(map[string]json.RawMessage)
			// update values of parameters according to map
			for parameterName, value := range parametersForThisProcess {
				// check for replacement parameter names
				if newParameterName, newParameterFound := newElementNames.NewParameterMapping[parameterName]; newParameterFound {
					newParameters[newParameterName] = value
				} else if typeChangeDetails, typeChangeFound := newElementNames.ParameterTypeChanges[parameterName]; typeChangeFound {
					errs = append(errs, fmt.Sprintf("detected a parameter %q of process %q requires type change of %s. Please manually convert this parameter", parameterName, processName, string(typeChangeDetails)))
				} else if !search.InStrings(newElementNames.DeprecatedParameters, parameterName) {
					newParameters[parameterName] = value
				}

			}
			// add new defaults if any specified
			for newParameterName, defaultValue := range newElementNames.NewParameters {
				newParameters[newParameterName] = defaultValue
			}

			// replace connections
			for parameterName, newParameterName := range newElementNames.NewParameterMapping {
				for i := range connections {
					connections[i] = replaceConnection(connections[i], processName, parameterName, newParameterName)
				}
			}
			parameters[processName] = newParameters
		} else {
			parameters[processName] = in.Parameters[processName]
		}
	}

	if len(errs) > 0 {
		return parameters, connections, wtype.NewWarningf(strings.Join(errs, ";"))
	}

	return parameters, connections, nil
}

func replaceConnection(connection workflow.Connection, processToReplace, parameterToReplace, newParameterName string) workflow.Connection {
	var newConnection workflow.Connection
	if connection.Src.Process == processToReplace && connection.Src.Port == parameterToReplace {
		newConnection.Src = workflow.Port{
			Process: connection.Src.Process,
			Port:    newParameterName,
		}
	} else {
		newConnection.Src = connection.Src
	}

	if connection.Tgt.Process == processToReplace && connection.Tgt.Port == parameterToReplace {
		newConnection.Tgt = workflow.Port{
			Process: connection.Tgt.Process,
			Port:    newParameterName,
		}
	} else {
		newConnection.Tgt = connection.Tgt
	}
	return newConnection
}

type nothingToConvert struct {
	ErrMessage string
}

func (err nothingToConvert) Error() string {
	return err.ErrMessage
}

func convertBundleWithArgs(conversionMapFileName, bundleFileName, outPutFileName string) error {

	cFile, err := os.Open(conversionMapFileName)

	if err != nil {
		return err
	}
	defer cFile.Close() // nolint

	var c NewElementMappingDetails
	decConv := json.NewDecoder(cFile)
	if err := decConv.Decode(&c); err != nil {
		return err
	}

	if c.Empty() {
		return nothingToConvert{ErrMessage: "empty conversion map"}
	}

	inFile, err := os.Open(bundleFileName)
	if err != nil {
		return err
	}
	defer inFile.Close() // nolint

	var original executeutil.Bundle
	dec := json.NewDecoder(inFile)
	if err := dec.Decode(&original); err != nil {
		return err
	}

	var bundle executeutil.Bundle = original

	if !containsSomethingToConvert(original, c.ConversionDetails) {
		return nothingToConvert{ErrMessage: "nothing to convert in bundle file"}
	}
	var warning error

	bundle.Parameters, bundle.Connections, warning = convertParametersAndConnections(original, c.ConversionDetails)
	bundle.Processes = convertProcesses(original.Processes, c.ConversionDetails)
	outFile, err := os.Create(outPutFileName)
	if err != nil {
		return err
	}
	defer outFile.Close() // nolint

	enc := json.NewEncoder(outFile)
	enc.SetIndent("", "  ")

	err = enc.Encode(bundle)

	if err != nil {
		return err
	}

	return warning
}

func convertBundle(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}
	return convertBundleWithArgs(viper.GetString("conversionMap"), viper.GetString("bundle"), viper.GetString("output"))
}

func init() {
	c := convertBundleCmd
	convertCmd.AddCommand(c)

	flags := c.Flags()
	flags.String("bundle", "", "Input bundle")
	flags.String("conversionMap", "", "conversion map file")
	flags.String("output", "output.json", "default output bundle name in antha format")
}
