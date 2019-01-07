// wutil/biotree.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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

package wtype

import (
	"io/ioutil"
	"regexp"
	"strings"
)

// we use the open tree of life to define taxonomic relationships
type TOL struct {
	UID      string
	Name     string
	Taxid    string
	Parent   *TOL
	Depth    int
	Children []*TOL
}

// read the data from the opentree file
func Load_TOL(filename string) (*TOL, *map[string]*TOL) {
	cnts, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	// the TOL file conveniently declares parent/child relationships, we can use these to simply index things

	index := make(map[string]*TOL)

	lines := strings.Split(string(cnts), "\n")

	var top *TOL

	for i, line := range lines {
		if i == 0 {
			continue
		}
		tx := strings.Split(line, "|")

		if len(tx) != 8 {
			continue
		}

		uid := strings.TrimSpace(tx[0])
		parent := strings.TrimSpace(tx[1])
		name := strings.TrimSpace(tx[2])
		taxid := parse_out_taxid(strings.TrimSpace(tx[4]))

		node := TOL{uid, name, taxid, index[parent], 0, make([]*TOL, 0)}

		index[name] = &node
		index[uid] = &node

		taxtox := strings.Split(strings.TrimSpace(tx[4]), ",")

		for _, tax := range taxtox {
			index[tax] = &node
		}

		parent_node := index[parent]

		if parent_node != nil {
			node.Depth = parent_node.Depth + 1
			parent_node.Children = append(parent_node.Children, &node)
		} else {
			top = &node
		}
	}
	return top, &index
}

// extract the NCBI taxid from the tree file
func parse_out_taxid(tok string) string {
	taxid := ""

	stx := strings.Split(tok, ",")
	re, err := regexp.Compile("ncbi")
	if err != nil {
		panic(err)
	}

	for _, tk2 := range stx {
		mt := re.MatchString(tk2)
		if mt {
			tx2 := strings.Split(tk2, ":")
			taxid = tx2[1]
			break
		}
	}
	return taxid
}

// returns the string name of the LCA if t is the ancestor of t2
func (t TOL) IsAncestorOf(t2 *TOL) string {
	if t.Depth == 0 {
		// the root
		return t.Name
	}

	tax := make([]string, 0, (*t2).Depth+1)

	tax = (*t2).Get_taxonomy(tax)

	ix := IndexOfString(t.Name, &tax)

	if ix == -1 {
		return ""
	} else {
		return tax[ix]
	}
}

// extract the lineage of one particular node
func (t TOL) Get_taxonomy(arr []string) []string {
	// put your name in arr, then get your ancestors to do the same

	arr = append(arr, t.Name)
	if t.Parent != nil {
		arr = (*t.Parent).Get_taxonomy(arr)
	}
	return arr
}

// find a node by name
func (t TOL) Find_string(name string) *TOL {
	if t.Name == name {
		return &t
	} else if t.Children == nil {
		return nil
	} else {
		for _, np := range t.Children {
			f := (*np).Find_string(name)
			if f != nil {
				return f
			}
		}

		return nil
	}
}

func IndexOfString(query string, pa *[]string) int {
	for i := 0; i < len(*pa); i++ {
		if (*pa)[i] == query {
			return i
		}
	}

	return -1
}
