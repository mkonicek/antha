package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/execute/executeutil"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var convertAllBundlesCmd = &cobra.Command{
	Use:   "all-bundle-parameters",
	Short: "update all bundle parameters according to the map of old parameter names to new found in metadata",
	RunE:  convertAllBundles,
}

//
func convertAllBundles(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	// find metadata files
	elements := newElements()
	if err := filepath.Walk(viper.GetString("elementsDir"), elements.Walk); err != nil {
		return err
	}

	var allBundles []*executeutil.TestInput
	var err error

	if path := viper.GetString("specificFile"); path != "" {

		dir, _ := filepath.Split(path)

		allBundles = append(allBundles, &executeutil.TestInput{
			BundlePath: path,
			Dir:        dir,
		})
	} else if allBundles, err = executeutil.FindTestInputs(viper.GetString("bundlesDir")); err != nil {
		return err
	}

	var errs []string

	for _, elem := range elements.Elements {

		metadataFileName := elem.MetadataPath()

		cFile, err := os.Open(metadataFileName)

		if err != nil {
			cFile.Close() //nolint
		} else {
			var c NewElementMappingDetails
			decConv := json.NewDecoder(cFile)
			if err := decConv.Decode(&c); err != nil {
				errs = append(errs, "error decoding to NewElementMappingDetails for "+metadataFileName+": ", err.Error())
			}
			cFile.Close() //nolint

			if !c.Empty() {
				for _, bundle := range allBundles {

					dir, fileName := filepath.Split(bundle.BundlePath)

					if !strings.HasPrefix(fileName, viper.GetString("addPrefix")) {
						fileName = viper.GetString("addPrefix") + fileName
					}
					newPath := filepath.Join(dir, fileName)

					err := convertBundleWithArgs(metadataFileName, bundle.BundlePath, newPath)

					if err != nil {
						switch err.(type) {
						case wtype.Warning:
							errs = append(errs, metadataFileName+" + "+bundle.BundlePath+": "+err.Error())
						default:
							// ignore
						}
					} else {
						// update bundle name, in the case it will be re-modified
						bundle.BundlePath = newPath
					}
				}
			}
		}
	}

	if len(errs) > 0 {
		return errors.Errorf(strings.Join(errs, "\n"))
	}

	return nil
}

func init() {
	c := convertAllBundlesCmd
	convertCmd.AddCommand(c)

	flags := c.Flags()
	flags.String("elementsDir", ".", "root directory to search for metadata files with new element mapping")
	flags.String("bundlesDir", ".", "root directory to search for test bundles to update if specificFile is not specified")
	flags.String("addPrefix", "", "adds a common prefix to the start of all updated bundle files")
	flags.String("specificFile", "", "specify a single bundle file to convert with all metadata files found in rootDir")
}
