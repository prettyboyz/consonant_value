// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package search

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// ErrFieldMismatch is returned when a field is to be loaded into a different
// than the one it was stored from, or when a field is missing or unexported in
// the destination struct.
type ErrFieldMismatch struct {
	FieldName string
	Reason    string
}

func (e *ErrFieldMismatch) Error() string {
	return fmt.Sprintf("search: cannot load field %q: %s", e.FieldName, e.Reason)
}

// ErrFacetMismatch is returned when a facet is to be loaded into a different
// type than the one it was stored from, or when a field is missing or
// unexported in the destination struct. StructType is the type of the struct
// pointed to by the destination argument passed to Iterator.Next.
type ErrFacetMismatch struct {
	StructType reflect.Type
	FacetName  string
	Reason     string
}

func (e *ErrFacetMismatch) Error() string {
	return fmt.Sprintf("search: cannot load facet %q into a %q: %s", e.FacetName, e.StructType, e.Reason)
}

// structCodec defines how to convert a given struct to/from a search document.
type structCodec struct {
	// byIndex returns the struct tag for the i'th struct field.
	byIndex []structTag

	// fieldByName returns the index of the struct field for the given field name.
	fieldByName map[string]int

	// facetByName returns the index of the struct field for the given facet name,
	facetByName map[string]int
}

// structTag holds a structured version of each struct field's parsed tag.
type structTag struct {
	name   string
	facet  bool
	ignore bool
}

var (
	codecsMu sync.RWMutex
	codecs   = map[reflect.Type]*structCodec{}
)

func loadCodec(t reflect.Type) (*structCodec, error) {
	codecsMu.RLock()
	codec, ok := codecs[t]
	codecsMu.RUnlock()
	if ok {
		return codec, nil
	}

	codecsMu.Lock()
	defer codecsMu.Unlock()
	if codec, ok := codecs[t]; ok {
		return codec, nil
	}

	codec = &structCodec{
		fieldByName: make(map[string]int),
		facetByName: make(map[string]int),
	}

	for i, I := 0, t.NumField(); i < I; i++ {
		f := t.Field(i)
		name, opts := f.Tag.Get("search"),