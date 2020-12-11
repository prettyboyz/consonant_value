package toml

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"reflect"
	"strings"
	"time"
)

func e(format string, args ...interface{}) error {
	return fmt.Errorf("toml: "+format, args...)
}

// Unmarshaler is the interface implemented by objects that can unmarshal a
// TOML description of themselves.
type Unmarshaler interface {
	UnmarshalTOML(interface{}) error
}

// Unmarshal decodes the contents of `p` in TOML format into a pointer `v`.
func Unmarshal(p []byte, v interface{}) error {
	_, err := Decode(string(p), v)
	return err
}

// Primitive is a TOML value that hasn't been decoded into a Go value.
// When using the various `Decode*` functions, the type `Primitive` may
// be given to any value, and its decoding will be delayed.
//
// A `Primitive` value can be decoded using the `PrimitiveDecode` function.
//
// The underlying representation of a `Primitive` value is subject to change.
// Do not rely on it.
//
// N.B. Primitive values are still parsed, so using them will only avoid
// the overhead of reflection. They can be useful when you don't know the
// exact type of TOML data until run time.
type Primitive struct {
	undecoded interface{}
	context   Key
}

// DEPRECATED!
//
// Use MetaData.PrimitiveDecode instead.
func PrimitiveDecode(primValue Primitive, v interface{}) error {
	md := MetaData{decoded: make(map[string]bool)}
	return md.unify(primValue.undecoded, rvalue(v))
}

// PrimitiveDecode is just like the other `Decode*` functions, except it
// decodes a TOML value that has already been parsed. Valid primitive values
// can *only* be obtained from values filled by the decoder functions,
// including this method. (i.e., `v` may contain more `Primitive`
// values.)
//
// Meta data for primitive values is included in the meta data returned by
// the `Decode*` functions with one exception: keys returned by the Undecoded
// method will only reflect keys that were decoded. Namely, any keys hidden
// behind a Primitive will be considered undecoded. Executing this method will
// update the undecoded keys in the meta data. (See the example.)
func (md *MetaData) PrimitiveDecode(primValue Primitive, v interface{}) error {
	md.context = primValue.context
	defer func() { md.context = nil }()
	return md.unify(primValue.undecoded, rvalue(v))
}

// Decode will decode the contents of `data` in TOML format into a pointer
// `v`.
//
// TOML hashes correspond to Go structs or maps. (Dealer's choice. They can be
// used interchangeably.)
//
// TOML arrays of tables correspond to either a slice of structs or a slice
// of maps.
//
// TOML datetimes correspond to Go `time.Time` values.
//
// All other TOML types (float, string, int, bool and array) correspond
// to the obvious Go types.
//
// An exception to the above rules is if a type implements the
// encoding.TextUnmarshaler interface. In this case, any primitive TOML value
// (floats, strings, integers, booleans and datetimes) will be converted to
// a byte string and given to the value's UnmarshalText method. See the
// Unmarshaler example for a demonstration with time duration strings.
//
// Key mapping
//
// TOML keys can map to either keys in a Go map or field names in a Go
// struct. The special `toml` struct tag may be used to map TOML keys to
// struct fields that don't match the key name exactly. (See the example.)
// A case insensitive match to struct names will be tried if an exact match
// can't be found.
//
// The mapping between TOML values and Go values is loose. That is, there
// may exist TOML values that cannot be placed into your representation, and
// there may be parts of your representation that do not correspond to
// TOML values. This loose mapping can be made stricter by using the IsDefined
// and/or Undecoded methods on the MetaData returned.
//
// This decoder will not handle cyclic types. If a cyclic type is passed,
// `Decode` will not terminate.
func Decode(data string, v interface{}) (MetaData, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return MetaData{}, e("Decode of non-pointer %s", reflect.TypeOf(v))
	}
	if rv.IsNil() {
		return MetaData{}, e("Decode of nil %s", reflect.TypeOf(v))
	}
	p, err := parse(data)
	if err != nil {
		return MetaData{}, err
	}
	md := MetaData{
		p.mapping, p.types, p.ordered,
		make(map[string]bool, len(p.ordered)), nil,
	}
	return md, md.unify(p.mapping, indirect(rv))
}

// DecodeFile is just like Decode, except it will automatically read the
// contents of the file at `fpath` and decode it for you.
func DecodeFile(fpath string, v interface{}) (MetaData, error) {
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		return MetaData{}, err
	}
	return Decode(string(bs), v)
}

// DecodeReader is just like Decode, except it will consume all bytes
// from the reader and decode it for you.
func DecodeReader(r io.Reader, v interface{}) (MetaData, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return MetaData{}, err
	}
	return Decode(string(bs), v)
}

// unify performs a sort of type unification based on the structure of `rv`,
// which is the client representation.
//
// Any type mismatch produces an error. Finding a type that we don't know
// how to handle produces an unsupported type error.
func (md *MetaData) unify(data interface{}, rv reflect.Value) error {

	// Special case. Look for a `Primitive` value.
	if rv.Type() == reflect.TypeOf((*Primitive)(nil)).Elem() {
		// Save the undecoded data and the key context into the primitive
		// value.
		context := make(Key, len(md.context))
		copy(context, md.context)
		rv.Set(reflect.ValueOf(Primitive{
			undecoded: data,
			context:   context,
		}))
		return nil
	}

	// Special case. Unmarshaler Interface support.
	if rv.CanAddr() {
		if v, ok := rv.Addr().Interface().(Unmarshaler); ok {
			return v.UnmarshalTOML(data)
		}
	}

	// Special case. Handle time.Time values specifically.
	// TODO: Remove this code when we decide to drop support for Go 1.1.
	// This isn't necessary in Go 1.2 because time.Time satisfies the encoding
	// interfaces.
	if rv.Type().AssignableTo(rvalue(time.Time{}).Type()) {
		return md.unifyDatetime(data, rv)
	}

	// Special case. Look for a value satisfying the TextUnmarshaler interface.
	if v, ok := rv.Interface().(TextUnmarshaler); ok {
		return md.unifyText(data, v)
	}
	// BUG(burntsushi)
	// The behavior here is incorrect whenever a Go type satisfies the
	// encoding.TextUnmarshaler interface but also corresponds to a TOML
	// hash or array. In particular, the unmarshaler should only be applied
	// to primitive TOML values. But at this point, it will be applied to
	// all kinds of values and produce an incorrect error whenever those values
	// are hashes or arrays (including arrays of tables).

	k := rv.Kind()

	// laziness
	if k >= reflect.Int && k <= reflect.Uint64 {
		return md.unifyInt(data, rv)
	}
	switch k {
	case reflect.Ptr:
		elem := reflect.New(rv.Type().Elem())
		err := md.unify(data, reflect.Indirect(elem))
		if err != nil {
			return err
		}
		rv.Set(elem)
		return nil
	case reflect.Struct:
		return md.unifyStruct(data, rv)
	case reflect.Map:
		return md.unifyMap(data, rv)
	case reflect.Array:
		return md.unifyArray(data, rv)
	case reflect.Slice:
		return md.unifySlice(data, rv)
	case reflect.String:
		return md.unifyString(data, rv)
	case reflect.Bool:
		return md.unifyBool(data, rv)
	case reflect.Interface:
		// we only support empty interfaces.
		if rv.NumMethod() > 0 {
			return e("unsupported type %s", rv.Type())
		}
		return md.unifyAnything(data, rv)
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		return md.unifyFloat64(data, rv)
	}
	return e("unsupported type %s", rv.Kind())
}

func (md *MetaData) unifyStruct(mapping interface{}, rv reflect.Value) error {
	tmap, ok := mapping.(map[string]interface{})
	if !ok {
		if mapping == nil {
			return nil
		}
		return e("type mismatch for %s: expected table but found %T",
			rv.Type().String(), mapping)
	}

	for key, datum := range tmap {
		var f *field
		fields := cachedTypeFields(rv.Type())
		for i := range fields {
			ff := &fields[i]
			if ff.name == key {
				f = ff
				break
			}
			if f == nil && strings.EqualFold(ff.name, key) {
				f = ff
			}
		}
		if f != nil {
			subv := rv
			for _, i := range f.index {
				subv = indirect(subv.Field(i))
			}
			if isUnifiable(subv) {
				md.decoded[md.context.add(key).String()] = true
				md.context = append(md.context, key)
				if err := md.unify(datum, subv); err != nil {
					return err
				}
				md.context = md.context[0 : len(md.context)-1]
			} else if f.name != "" {
				// Bad user! No soup for you!
				return e("cannot write unexported field %s.%s",
					rv.Type().String(), f.name)
			}
		}
	}
	return nil
}

func (md *MetaData) unifyMap(mapping interface{}, rv reflect.Value) error {
	tmap, ok := mapping.(map[string]interface{})
	if !ok {
		if tmap == nil {
			return nil
		}
		return badtype("map", mapping)
	}
	if rv.IsNil() {
		rv.Set(reflect.MakeMap(rv.Type()))
	}
	for k, v := range tmap {
		md.decoded[md.context.add(k).String()] = true
		md.context = append(md.context, k)

		rvkey := indirect(reflect.New(rv.Type().Key()))
		rvval := reflect.Indirect(reflect.New(rv.Type().Elem()))
		if err := md.unify(v, rvval); err != nil {
			return err
		}
		md.context = md.context[0 : len(md.context)-1]

		rvkey.SetString(k)
		rv.SetMapIndex(rvkey, rvval)
	}
	return nil
}

func (md *MetaData) unifyArray(data interface{}, rv reflect.Value) error {
	datav := reflect.ValueOf(data)
	if datav.Kind() != reflect.Slice {
		if !datav.IsValid() {
			return nil
		}
		return badtype("slice", data)
	}
	sliceLen := datav.Len()
	if sliceLen != rv.Len() {
		return e("expected array length %d; got TOML array of length %d",
			rv.Len(), sliceLen)
	}
	return md.unifySliceArray(datav, rv)
}

func (md *MetaData) unifySlice(data interface{}, rv reflect.Value) error {
	datav := reflect.ValueOf(data)
	if datav.Kind() != reflect.Slice {
		if !datav.IsValid() {
			return nil
		}
		return badtype("slice", data)
	}
	n := datav.Len()
	if rv.IsNil() || rv.Cap() < n {
		rv.Set(reflect.MakeSlice(rv.Type(), n, n))
	}
	rv.SetLen(n)
	return md.unifySliceArray(datav, rv)
}

func (md *MetaData) unifySliceArray(data, rv reflect.Value) error {
	sliceLen := data.Len()
	for i := 0; i < sliceLen; i++ {
		v := data.Index(i).Interface()
		sliceval := indirect(rv.Index(i))
		if err := md.unify(v, sliceval); err != nil {
			return err
		}
	}
	return nil
}

func (md *MetaData) unifyDatetime(data interface{}, rv reflect.Value) error {
	if _, ok := data.(time.Time); ok {
		rv.Set(reflect.ValueOf(data))
		return nil
	}
	return badtype("time.Time", data)
}

func (md *MetaData) unifyString(data interface{}, rv reflect.Valu