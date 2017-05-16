// compile.go: Part of the Antha language
// Copyright (C) 2017 The Antha authors. All rights reserved.
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
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/antha-lang/antha/antha/format"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var formatCmd = &cobra.Command{
	Use:   "format [path ...]",
	Short: "Format an antha element",
	RunE:  runFormat,
}

type walker struct {
	Write bool
}

const chmodSupported = runtime.GOOS != "windows"

func (w *walker) Walk(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return nil
	}

	if !strings.HasSuffix(path, ".an") {
		return nil
	}

	src, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	out, err := format.Source(src)
	if err != nil {
		return err
	}

	if !w.Write {
		os.Stdout.Write(out)
		return nil
	}

	if bytes.Equal(src, out) {
		return nil
	}

	bak, err := ioutil.TempFile(filepath.Dir(path), filepath.Base(path))
	if err != nil {
		return err
	}
	defer bak.Close()
	defer func() {
		if err != nil {
			os.Remove(bak.Name())
		}
	}()

	if chmodSupported {
		if err = bak.Chmod(fi.Mode()); err != nil {
			return err
		}
	}

	_, err = io.Copy(bak, bytes.NewReader(out))
	if err != nil {
		return err
	}

	err = os.Rename(bak.Name(), path)

	return err
}

func runFormat(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	w := walker{
		Write: viper.GetBool("write"),
	}

	for _, arg := range args {
		if err := filepath.Walk(arg, w.Walk); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	c := formatCmd
	flags := c.Flags()
	RootCmd.AddCommand(c)

	flags.BoolP("write", "w", false, "Write result to (source) file instead of stdout")
}
