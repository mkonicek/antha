package cmd

import "github.com/spf13/cobra"

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert various file formats",
}

func init() {
	c := convertCmd
	RootCmd.AddCommand(c)
}
