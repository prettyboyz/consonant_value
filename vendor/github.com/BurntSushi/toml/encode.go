package toml

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type tomlEncodeError struct{ error }

var (
	errArrayMixedElementTypes = errors.New(
		"toml: cannot encode array with mixed element types")
	errArrayNilElement = errors.New(
		"toml: cannot encode array with nil element")
	errNonString = errors.New(
		"toml: cannot encode a map with non-string key type")
	errAnonNonStruct = errors.New(
		"toml: cannot encode an anonymous field that is not a struct")
	errArrayNoTable = errors.New(
		"toml: TOML array element cannot contain a table")
	errNoKey = errors.New(
		"toml: top-level values must be Go maps or structs")
	errAnything = errors.New("") // used in testing
)

var quotedReplacer = strings.NewReplacer(
	"\t", "\\t",
	"\n", "\\n",
	"\r", "\\r",
	"\"", "\\\"",
	"\\", "\\\\",
)

// Encoder controls the encoding of Go values to a TOML document to some
// io.Writer.
//
// The indentation level can be controlled with the Indent field.
type Encoder struct {
	// A single indentation level. By default it is two spaces.
	Indent string

	// hasWritten is whether we have written any output to w yet.
	hasWritten bool
	w          *bufio.Writer
}

// NewEncoder returns a TOML encoder that encodes Go values to the io.Writer
// given. By default, a single indentation level is 2 spaces.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w:      bufio.NewWriter(w),
		Indent: "  ",
	}
}

// Encode writes a TOML representation of the Go value to the underlying
// io.Writer. If the value given cannot be encoded to a valid TOML document,
// then an error is returned.
//
// The mapping between Go values and TOML values should be precisely the same
// as for the Decode* functions. Similarly, the TextMarshaler interface is
// supported by encoding the resulting bytes as strings. (If you want to write
// arbitrary binary data then you will need to use something like base64 since
// TOML does not have any binary types.)
//
// When encoding TOML hashes (i.e., Go maps or structs), keys without any
// sub-hashes are encoded first.
//
// If a Go map is encoded, then its keys are sorted alphabetically for
// deterministic output. More control over this behavior may be provided if
// there is demand for it.
//
// Encoding Go values without a corresponding TOML representation---like map
// types with non-string keys---will cause an error to be returned. Similarly
// for mixed arrays/slices, arrays/slices with nil elements, embedded
// non-struct types and nested slices containing maps or structs.
// (e.g., [][]map[string]string is not allowed but []map[string]string is OK
// and so is []map[string][]string.)
func (enc *Encoder) Encode(v interface{}) error {
	rv := eindirect(reflect.ValueOf(v))
	if err := enc.safeEncode(Key([]string{}), rv); err != nil {
		return err
	}
	return enc.w.Flush()
}

func (enc *Encoder) safeEncode(key Key, rv reflect.Value) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if terr, ok := r.(tomlEncodeError); ok {
				err = terr.error
				return
			}
			panic(r)
		}
	}()
	enc.encode(key, rv)
	return nil
}

func (enc *Encoder) encode(key Key, rv reflect.Value) {
	// Special case. Time needs to be in ISO8601 format.
	// Special case. If we can marshal the type to text, then we used that.
	// Basically, this prevents the encoder for handling these types as
	// generic structs (or whatever the underlying type of a TextMarshaler is).
	switch rv.Interface().(type) {
	case time.Time, TextMarshaler:
		enc.keyEqElement(key, rv)
		return
	}

	k := rv.Kind()
	switch k {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.String, reflect.Bool:
		enc.keyEqElement(key, rv)
	case reflect.Array, reflect.Slice:
		if typeEqual(tomlArrayHash, tomlTypeOfGo(rv)) {
			enc.eArrayOfTables(key, rv)
		} else {
			enc.keyEqElement(key, rv)
		}
	case reflect.Interface:
		if rv.IsNil() {
			return
		}
		enc.encode(key, rv.Elem())
	case reflect.Map:
		if rv.IsNil() {
			return
		}
		enc.eTable(key, rv)
	case reflect