// antha/AnthaStandardLibrary/Packages/enzymes/Find.go: Part of the Antha language
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

// Package search is a utility package providing functions useful for:
// Searching for a target entry in a slice;
// Removing duplicate values from a slice;
// Comparing the Name of two entries of any type with a Name() method returning a string.
// FindAll instances of a target string within a template string.
package search

import "sort"

// BinarySearchInts searchs for an int in a slice of ints using a binary search algorithm.
func BinarySearchInts(list []int, target int) bool {
	sort.Ints(list)
	i := sort.Search(len(list),
		func(i int) bool { return list[i] >= target })
	if i < len(list) && list[i] == target {
		return true
	}
	return false
}

// InInts searchs for an int in a slice of ints
func InInts(list []int, target int) bool {
	for _, num := range list {
		if num == target {
			return true
		}
	}
	return false
}

// PositionsInInts searchs for a target int in a slice of ints and returns all positions found.
func PositionsInInts(list []int, target int) []int {

	var positions []int
	for position, num := range list {
		if num == target {
			positions = append(positions, position)
		}
	}
	return positions
}
