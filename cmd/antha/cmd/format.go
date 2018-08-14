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
	"fmt"
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
	WriteToFile bool
}

const chmodSupported = runtime.GOOS != "windows"

func (w *walker) Walk(path string, fi os.FileInfo, err error) error {
	if err := w.walk(path, fi, err); err != nil {
		return fmt.Errorf("%s: %s", path, err)
	}
	return nil
}

func (w *walker) walk(path string, fi os.FileInfo, err error) error {
	if err != nil || fi.IsDir() || !strings.HasSuffix(path, ".an") {
		return err
	} else if src, err := ioutil.ReadFile(path); err != nil {
		return err
	} else if out, err := format.Source(src); err != nil {
		return err
	} else if !w.WriteToFile {
		os.Stdout.Write(out) // nolint
		return nil
	} else if bytes.Equal(src, out) {
		return nil
	} else if bak, err := ioutil.TempFile(filepath.Dir(path), filepath.Base(path)); err != nil {
		return err
	} else {
		defer func() {
			if err != nil {
				os.Remove(bak.Name()) // nolint
			}
		}()
		defer bak.Close() // nolint: errcheck

		if chmodSupported {
			// maintain the same mode as the read file
			if err = bak.Chmod(fi.Mode()); err != nil {
				return err
			}
		}

		if _, err = io.Copy(bak, bytes.NewReader(out)); err != nil {
			return err
		} else {
			return os.Rename(bak.Name(), path)
		}
	}
}

func runFormat(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	w := walker{
		WriteToFile: viper.GetBool("write"),
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
