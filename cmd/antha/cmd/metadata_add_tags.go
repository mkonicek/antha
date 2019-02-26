// metadata_add_tags.go: Part of the Antha language
// Copyright (C) 2018 The Antha authors. All rights reserved.
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
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	tagsField = "tags"
	nameField = "name"
)

// TODO: move from blind structure to structured data file
type metadata map[string]json.RawMessage

var addTagsMetadataCmd = &cobra.Command{
	Use:   "add-tags tag1 tag2 ...",
	Short: "Add tags to metadata files for elements",
	RunE:  addTagsMetadata,
}

type element struct {
	Dir  string
	Path string
	Name string
}

func (e *element) MetadataPath() string {
	return filepath.Join(e.Dir, "metadata.json")
}

func getProtocolName(path string) (string, error) {
	// NB(ddn): Don't want to implement parser so hack this

	pat, err := regexp.Compile(`(?m:^protocol (\S+))`)
	if err != nil {
		return "", errors.WithStack(err)
	}

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.WithStack(err)
	}

	matches := pat.FindStringSubmatch(string(bs))
	if len(matches) == 0 {
		return "", errors.New("no protocol name found")
	}

	return matches[1], nil
}

type elements struct {
	Elements []*element
	seen     map[string]bool
}

func newElements() *elements {
	return &elements{
		seen: make(map[string]bool),
	}
}

func (e *elements) Walk(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return nil
	}
	if filepath.Ext(path) != ".an" {
		return nil
	}

	dir := filepath.Dir(path)
	if e.seen[dir] {
		return nil
	}

	e.seen[dir] = true

	name, err := getProtocolName(path)
	if err != nil {
		return errors.Errorf("element %s: %s", path, err)
	}

	e.Elements = append(e.Elements, &element{
		Dir:  dir,
		Path: path,
		Name: name,
	})

	return nil
}

func createJSON(path string, obj interface{}) error {
	bs, err := json.Marshal(obj)
	if err != nil {
		return errors.WithStack(err)
	}

	f, err := os.Create(path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close() // nolint

	_, err = io.Copy(f, bytes.NewReader(bs))
	return errors.WithStack(err)
}

func openMetadata(element *element) (metadata, error) {
	_, err := os.Stat(element.MetadataPath())
	if err != nil && !os.IsNotExist(err) {
		return nil, errors.WithStack(err)
	}

	mdata := make(metadata)

	if err == nil {
		bs, err := ioutil.ReadFile(element.MetadataPath())
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if err := json.Unmarshal(bs, &mdata); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return mdata, nil
}

func ensureMetadata(element *element) error {
	mdata, err := openMetadata(element)
	if err != nil {
		return err
	}

	nameBs, err := json.Marshal(element.Name)
	if err != nil {
		return errors.WithStack(err)
	}

	mdata[nameField] = nameBs

	return createJSON(element.MetadataPath(), mdata)
}

func ensureTags(element *element, tags []string) error {
	mdata, err := openMetadata(element)
	if err != nil {
		return err
	}

	var origTags []string
	if bs := mdata[tagsField]; len(bs) != 0 {
		if err := json.Unmarshal(bs, &origTags); err != nil {
			return errors.WithStack(err)
		}
	}

	var nextTags []string
	seen := make(map[string]bool)

	for _, v := range append(origTags, tags...) {
		if seen[v] {
			continue
		}
		seen[v] = true
		nextTags = append(nextTags, v)
	}

	tagsBs, err := json.Marshal(nextTags)
	if err != nil {
		return errors.WithStack(err)
	}

	mdata[tagsField] = tagsBs

	return createJSON(element.MetadataPath(), mdata)
}

func addTagsMetadata(cmd *cobra.Command, args []string) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	if len(args) == 0 {
		return errors.New("no tags to add")
	}

	elements := newElements()
	if err := filepath.Walk(viper.GetString("rootDir"), elements.Walk); err != nil {
		return err
	}

	for _, elem := range elements.Elements {

		if err := ensureMetadata(elem); err != nil {
			return err
		}

		if err := ensureTags(elem, args); err != nil {
			return err
		}
	}

	return nil
}

func init() {
	c := addTagsMetadataCmd
	flags := c.Flags()

	metadataCmd.AddCommand(c)

	flags.String("rootDir", ".", "directory to start finding elements to add tags to")
}
