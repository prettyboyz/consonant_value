package toml

// Struct field handling is adapted from code in encoding/json:
//
// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the Go distribution.

import (
	"reflect"
	"sort"
	"sync"
)

// A field represents a single field found in a struct.
type field struct {
	name  string       // the name of the field (`toml` tag included)
	tag   bool         // whether field has a `toml` tag
	index []int        // represents the depth of an anonymous field
	typ   reflect.Type // the type of the field
}

// byName sorts field by name, breaking ties with depth,
// then breaking ties with "name came from toml tag", then
// breaking ties with index sequence.
type byName []field

func (x byName) Len() int { return len(x) }

func (x byName) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x byName) Less(i, j int) bool {
	if x[i].name != x[j].name {
		return x[i].name < x[j].name
	}
	if len(x[i].index) != len(x[j].index) {
		return len(x[i].index) < len(x[j].index)
	}
	if x[i].tag != x[j].tag {
		return x[i].tag
	}
	return byIndex(x).Less(i, j)
}

// byIndex sorts field by index sequence.
type byIndex []field

func (x byIndex) Len() int { return len(x) }

func (x byIndex) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

func (x byIndex) Less(i, j int) bool {
	for k, xik := range x[i].index {
		if k >= len(x[j].index) {
			return false
		}
		if xik != x[j].index[k] {
			return xik < x[j].index[k]
		}
	}
	return len(x[i].index) < len(x[j].index)
}

// typeFields returns a list of fields that TOML should recognize for the given
// type. The algorithm is breadth-first search over the set of structs to
// include - the top struct and then any reachable anonymous structs.
func typeFields(t reflect.Type) []field {
	// Anonymous fields to explore at the current level and the next.
	current := []field{}
	next := []field{{typ: t}}

	// Count of queued names for current level and the next.
	count := map[reflect.Type]int{}
	nextCount := map[reflect.Type]int{}

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool{}

	// Fields found.
	var fields []field

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true

			// Scan f.typ for fields to include.
			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				if sf.PkgPath != "" && !sf.Anonymous { // unexported
					continue
				}
				opts := getOptions(sf.Tag)
				if opts.skip {
					continue
				}
				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Ptr {
					// Follow pointer.
					ft = ft.Elem()
				}

				// Record found field and index sequence.
				if opts.name != "" || !sf.Anonymous || ft.Kind() != reflect.Struct {
					tagged := opts.name != ""
					name := opts.name
					if name == "" {
						name = sf.Name
					}
					fields = append(fields, field{name, tagged, index, ft})
					if count[f.typ] > 1 {
						// If there were multiple instances, add a second,
						// so that the annihilation code will see a duplicate.
						// It only cares about the distinction between 1 or 2,
						// so don't bother generating any more copies.
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}

				// Record new anonymous struct to explore in next round.
				nextCount[ft]++
				if nextCount[ft] == 1 {
					f := field{name: ft.Name(), index: index, typ: ft}
					next = append(next, f)
				}
			}
		}
	}

	sort.Sort(byName(fields))

	// Delete all fields that are hidden by the Go rules for embedded fields,
	// except that fields with TOML tags are promoted.

	// The fields are sorted in primary order of name, secondary order
	// of field ind