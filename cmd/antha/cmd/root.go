// root.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package cmd

import (
	api "github.com/antha-lang/antha/api/v1"
	"github.com/antha-lang/antha/component"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd is base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:          "antha",
	Short:        "Antha command line tool",
	SilenceUsage: true,
}

// Library of components available to workflows
var library []component.Component

func runComponents() (ret []component.Component) {
	for _, comp := range library {
		if comp.Stage != api.ElementStage_STEPS {
			continue
		}
		ret = append(ret, comp)
	}
	return
}

// Execute adds all child commands to the root command sets flags
// appropriately.  This is called by main.main(). It only needs to happen once
// to the rootCmd.
func Execute(lib []component.Component) error {
	library = lib

	return RootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Config file (default is $HOME/.antharun.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".antharun") // name of config file (without extension)
	viper.AddConfigPath("$HOME")     // adding home directory as first search path
	viper.SetEnvPrefix("antha")      // prefix of environment variables
	viper.AutomaticEnv()             // read in environment variables that match

	// If a config file is found, read it in.
	viper.ReadInConfig() // nolint
}
