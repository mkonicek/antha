package cmd

import (
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

	for _, elem := range elements.Elements {
		for _, bundle := range bundles.Bundles {
			convertBundleWithArgs(filepath.Join(elem.Dir, "metadata.json"), bundle.Path, bundle.Path)
		}
	}

	return nil
}

func init() {
	c := convertAllBundlesCmd
	convertCmd.AddCommand(c)

	flags := c.Flags()
	flags.String("rootDir", ".", "root directory to search for metadata files with new element mapping and test bundles to update")
}
