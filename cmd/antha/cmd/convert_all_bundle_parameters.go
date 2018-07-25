package cmd

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var convertAllBundlesCmd = &cobra.Command{
	Use:   "all-bundle-parameters",
	Short: "update all bundle parameters according to the map of old parameter names to new found in metadata",
	RunE:  convertAllBundles,
}

type bundle struct {
	Dir      string
	Path     string
	FileName string
}

type bundles struct {
	Bundles []*bundle
	seen    map[string]bool
}

func newBundles() *bundles {
	return &bundles{
		seen: make(map[string]bool),
	}
}

func (b *bundles) Walk(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return nil
	}
	if !strings.HasSuffix(path, ".bundle.json") {
		return nil
	}

	dir, fileName := filepath.Split(path)
	if b.seen[dir] {
		return nil
	}

	b.seen[dir] = true

	b.Bundles = append(b.Bundles, &bundle{
		Dir:      dir,
		Path:     path,
		FileName: fileName,
	})

	return nil
}

//
func convertAllBundles(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	// find metadata files
	elements := newElements()
	if err := filepath.Walk(viper.GetString("rootDir"), elements.Walk); err != nil {
		return err
	}

	bundles := newBundles()
	if err := filepath.Walk(viper.GetString("rootDir"), bundles.Walk); err != nil {
		return err
	}

	var errs []string

	for _, elem := range elements.Elements {

		metadataFileName := filepath.Join(elem.Dir, "metadata.json")

		cFile, err := os.Open(metadataFileName)

		if err != nil {
			errs = append(errs, metadataFileName+": ", err.Error())
			cFile.Close() //nolint
		} else {
			var c NewElementMappingDetails
			decConv := json.NewDecoder(cFile)
			if err := decConv.Decode(&c); err != nil {
				errs = append(errs, "error decoding to NewElementMappingDetails for "+metadataFileName+": ", err.Error())
			}
			cFile.Close() //nolint

			if !c.Empty() {
				for _, bundle := range bundles.Bundles {
					err := convertBundleWithArgs(metadataFileName, bundle.Path, bundle.Path)
					if err != nil {
						errs = append(errs, metadataFileName+" + "+bundle.Path+": "+err.Error())
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
	flags.String("rootDir", ".", "root directory to search for metadata files with new element mapping and test bundles to update")
}
