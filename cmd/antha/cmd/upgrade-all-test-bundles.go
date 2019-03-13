package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"

	"github.com/antha-lang/antha/execute/executeutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var convertBundlesWithTestsCmd = &cobra.Command{
	Use:   "all-bundles-with-tests",
	Short: "update all test bundle parameters",
	RunE:  updateTestBundles,
}

func updateTestBundles(cmd *cobra.Command, args []string) error {

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	allBundles, err := executeutil.FindTestInputs(viper.GetString("bundlesDir"))

	if err != nil {
		return err
	}

	if viper.GetString("driver") == "" {
		return errors.New("--driver must be secified when converting test bundles; please specify one; e.g. go://github.com/Synthace/instruction-plugins/PipetMax")
	}

	for _, bundle := range allBundles {

		fileName := bundle.BundlePath

		if bundle.Expected.CompareOutputs || viper.GetBool("all") {

			cmd := exec.Command("antha", "run", "--driver", viper.GetString("driver"), "--bundle", fileName, "--makeTestBundle", fileName)
			var out bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &out
			cmd.Stderr = &stderr
			err := cmd.Run()
			fmt.Println(fileName)
			if err != nil {
				fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			}

		}

	}
	return nil
}

func init() {

	c := convertBundlesWithTestsCmd

	convertCmd.AddCommand(c)

	flags := c.Flags()
	flags.String("bundlesDir", ".", "root directory to search for bundle files with test data to update")
	flags.String("driver", "", "driver flag to add to antha run command, e.g. go://github.com/Synthace/instruction-plugins/PipetMax")
	flags.Bool("all", false, "by default, only bundles with existing test data will be updated; use this flag to attempt to add test information to all bundles.")
}
