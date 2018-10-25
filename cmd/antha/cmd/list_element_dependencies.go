package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/antha-lang/antha/antha/compile"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listElementDependencies = &cobra.Command{
	Use:   "elementDependencies <files or directories>",
	Short: "List antha element dependencies",
	Long: `List antha element dependencies

To generate correct results, the value of outputPackage should match
that used in "antha compile".

By default, this command will show dependencies between import paths.  If you
want to see how elements in one directory depend on elements in a different
directory, use the "--byPath" option.

If you want to see how elements are related in an arbitrary projected space,
use the "--pathMatch" and "--pathReplace" options.
`,
	RunE: runListElementDependencies,
}

func runListElementDependencies(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	outPackage := viper.GetString("outputPackage")
	if len(outPackage) == 0 {
		return errors.New("outputPackage is not set")
	}

	root := compile.NewElementRoot(outPackage)

	var elements []*compile.Element

	for _, path := range args {
		if err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			} else if f.IsDir() {
				return nil
			} else if !isElementFile(f.Name()) {
				return nil
			} else if elem, err := processFile(root, path); err != nil {
				return err
			} else {
				elements = append(elements, elem)
				return nil
			}
		}); err != nil {
			return err
		}
	}

	var pat *regexp.Regexp
	match := viper.GetString("pathMatch")
	replace := viper.GetString("pathReplace")
	if len(match) > 0 || len(replace) > 0 {
		if len(match) == 0 || len(replace) == 0 {
			return errors.New("both pathMatch and pathReplace must be set")
		}
		var err error
		pat, err = regexp.Compile(match)
		if err != nil {
			return err
		}
	}

	convertTo := make(map[string]string)
	if viper.GetBool("byPath") {
		for _, elem := range elements {
			info := elem.Info()
			convertTo[info.ImportPath] = info.Path
		}
	}

	get := func(v string) string {
		if next, seen := convertTo[v]; seen {
			v = next
		}

		if pat == nil {
			return v
		}

		return pat.ReplaceAllString(v, replace)
	}

	type Pair struct {
		Src string
		Dst string
	}

	edges := make(map[Pair]bool)
	for _, elem := range elements {
		info := elem.Info()
		src := get(info.ImportPath)
		added := 0
		for _, dep := range info.DependsOn {
			// Hacky way to skip normal go packages
			if !strings.HasPrefix(dep, outPackage) {
				continue
			}

			dst := get(dep)

			if src == dst {
				continue
			}

			edges[Pair{Src: src, Dst: dst}] = true
			added++
		}

		// Don't drop a node just because it has no dependencies
		if added == 0 {
			edges[Pair{Src: src}] = true
		}
	}

	var buf bytes.Buffer
	buf.WriteString("digraph {\n") // nolint
	for pair := range edges {
		if len(pair.Dst) == 0 {
			fmt.Fprintf(&buf, "%q\n", pair.Src) // nolint
		} else {
			fmt.Fprintf(&buf, "%q -> %q\n", pair.Src, pair.Dst) // nolint
		}
	}
	buf.WriteString("}\n") // nolint

	io.Copy(os.Stdout, &buf) // nolint

	return nil
}

func init() {
	c := listElementDependencies
	listCmd.AddCommand(c)
	flags := c.Flags()
	flags.String("outputPackage", "", "base package name for generated files")
	flags.Bool("byPath", false, "if set, show dependencies on file paths instead of import paths")
	flags.String("pathMatch", "", "regex substitution to apply to paths, e.g., prefix(\\w+)")
	flags.String("pathReplace", "", "regex substitution to apply to paths, e.g., $1")
}
