package cmd

import (
	"encoding/json"
	"os"
	"fmt"
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
	OldElementName string `json:"old-element-name"`
	NewElementName string `json:"new-element-name"`
	NewParameterMapping map[string]string `json:"new-parameter-mapping"`
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
		}
	}

	return ret
}

func convertParametersAndConnections(in anthaBundle, newElementNames ConversionDetails) (map[string]map[string]json.RawMessage, []workflow.Connection) {
	
	parameters := make(map[string]map[string]json.RawMessage)
	connections := in.Connections

	for processName, element := range in.Processes {
		
		// check if element name is to be replaced
		if element.Component == newElementNames.OldElementName {
			// get existing parameters
			parametersForThisProcess, found := in.Parameters[processName]
			if !found{
				panic(fmt.Errorf("parameters not found for %s",processName))
			}
			newParameters := make(map[string]json.RawMessage)
			// update values of parameters according to map
			for parameterName, value := range parametersForThisProcess{
				// check for replacement parameter names
				if newParameterName, newParameterFound := newElementNames.NewParameterMapping[parameterName]; newParameterFound{
					newParameters[newParameterName] = value
				} else {
					newParameters[parameterName] = value
				}
			}
			// replace connections
			for parameterName, newParameterName := range newElementNames.NewParameterMapping {
				for i := range connections{
					connections[i] = replaceConnection(connections[i], processName, parameterName, newParameterName)
				}
			}
			parameters[processName] = newParameters
		} else{
			parameters[processName] = in.Parameters[processName]
		}
	}

	return parameters, connections
}

func replaceConnection(connection workflow.Connection, processToReplace, parameterToReplace, newParameterName string)workflow.Connection{
	var newConnection workflow.Connection	
	if connection.Src.Process == processToReplace && connection.Src.Port == parameterToReplace{
		newConnection.Src = workflow.Port{
			Process: connection.Src.Process,
			Port: newParameterName,
		}
	} else {
		newConnection.Src = connection.Src
	}
	
	if connection.Tgt.Process == processToReplace && connection.Tgt.Port == parameterToReplace{
		newConnection.Tgt = workflow.Port{
			Process: connection.Tgt.Process,
			Port: newParameterName,
		}
		panic(fmt.Sprintln(newConnection))
	} else {
		newConnection.Tgt = connection.Tgt
	}
	return newConnection
}
//
func convertBundle(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	inFile, err := os.Open(viper.GetString("bundle"))
	if err != nil {
		return err
	}
	defer inFile.Close() // nolint

	var original anthaBundle
	dec := json.NewDecoder(inFile)
	if err := dec.Decode(&original); err != nil {
		return err
	}
		
	cFile, err := os.Open(viper.GetString("conversionMap"))
	
	if err != nil {
		return err
	}
	defer cFile.Close() // nolint
	
	var c NewElementMappingDetails
	decConv := json.NewDecoder(cFile)
	if err := decConv.Decode(&c); err != nil {
		return err
	}
	

	var bundle anthaBundle
	bundle.Parameters, bundle.Connections = convertParametersAndConnections(original, c.ConversionDetails)
	bundle.Processes = convertProcesses(original.Processes, c.ConversionDetails)
	bundle.Config = original.Config
	
	outFile, err := os.Create(viper.GetString("output"))
	if err != nil {
		return err
	}
	defer outFile.Close() // nolint

	enc := json.NewEncoder(outFile)
	enc.SetIndent("", "  ")
	return enc.Encode(bundle)
}

func init() {
	c := convertBundleCmd
	convertCmd.AddCommand(c)

	flags := c.Flags()
	flags.String("bundle", "", "Input bundle")
	flags.String("conversionMap","", "conversion map file")
	flags.String("output", "output.json", "default output bundle name in antha format")
}
