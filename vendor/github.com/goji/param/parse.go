package param

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()

// Generic parse dispatcher. This function's signature is the interface of all
// parse functions. `key` is the entire key that is currently being parsed, such
// as "foo[bar][]". `keytail` is the portion of the string that the current
// parser is responsible for, for instance "[bar][]". `values` is the list of
// values assigned to this key, and `target` is where the resulting typed value
// should be Set() to.
func parse(key, keytail string, values []string, target reflect.Value) {
	t := target.Type()
	if reflect.PtrTo(t).Implements(textUnmarshalerType) {
		parseTextUnmarshaler(key, keytail, values, target)
		return
	}

	switch k := target.Kind(); k {
	case reflect.Bool:
		parseBool(key, keytail, values, target)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parseInt(key, keytail, values, target)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parseUint(key, keytail, values, target)
	case reflect.Float32, reflect.Float64:
		parseFloat(key, keytail, values, target)
	case reflect.Map:
		parseMap(key, keytail, values, target)
	case reflect.Ptr:
		parsePtr(key, keytail, values, target)
	case reflect.Slice:
		parseSlice(key, keytail, values, target)
	case reflect.String:
		parseString(key, keytail, values, target)
	case reflect.Struct:
		parseStruct(key, keytail, values, target)

	default:
		pebkac("unsupported object of type %v and kind %v.",
			target.Type(), k)
	}
}

// We pass down both the full key ("foo[bar][]") and the part the current layer
// is responsible for making sense of ("[bar][]"). This computes the other thing
// you probably want to know, which is the path you took to get here ("foo").
func kpath(key, keytail string) string {
	l, t := len(key), len(keytail)
	return key[:l-t]
}

// Helper for validating that a value has been passed exactly once, and that the
// user is not attempting to nest on the key.
func primitive(key, keytail string, tipe reflect.Type, values []string) {
	if keytail != "" {
		panic(NestingError{
			Key:     kpath(key, keytail),
			Type:    tipe,
			Nesting: keytail,
		})
	}
	if len(values) != 1 {
		panic(SingletonError{
			Key:    kpath(key, keytail),
			Type:   tipe,
			Values: values,
		})
	}
}

func keyed(tipe reflect.Type, key, keytail string) (string, string) {
	if keytail == "" || keytail[0] != '[' {
		panic(SyntaxError{
			Key:       kpath(key, keytail),
			Subtype:   MissingOpeningBracket,
			ErrorPart: keytail,
		})
	}

	idx := strings.IndexRune(keytail, ']')
	if idx == -1 {
		panic(SyntaxError{
			Key:       kpath(key, keytail),
			Subtype:   MissingClosingBracket,
			ErrorPart: keytail[1:],
		})
	}

	return keytail[1:idx], keytail[idx+1:]
}

func parseTextUnmarshaler(key, keytail string, values []string, target reflect.Value) {
	primitive(key, keytail, target.Type(), values)

	tu := target.Addr().Interface().(encoding.TextUnmarshaler)
	err := tu.UnmarshalText([]byte(values[0]))
	if err != nil {
		panic(TypeError{
			Key:  kpath(key, keytail),
			Type: target.Type(),
			Err:  err,
		})
	}
}

func parseBool(key, keytail string, values []string, target reflect.Value) {
	primitive(key, keytail, target.Type(), values)

	switch values[0] {
	case "true", "1", "on":
		target.SetBool(true)
	case "false", "0", "":
		target.SetBool(false)
	default:
		panic(TypeError{
			Key:  kpath(key, keytail),
			Type: target.Type(),
		})
	}
}

func parseInt(key, keytail string, values []string, target reflect.Value) {
	t := target.Type()
	primitive(key, keytail, t, values)

	i, err := strconv.ParseInt(values[0], 10, t.Bits())
	if err != nil {
		panic(TypeError{
			Key:  kpath(key, keytail),
			Type: t,
			Err:  err,
		})
	}
	target.SetInt(i)
}

func parseUint(key, keytail string, values []string, target reflect.Value) {
	t := target.Type()
	primitive(key, keytail, t, values)

	i, err := strconv.ParseUint(values[0], 10, t.Bits())
	if err != nil {
		panic(TypeError{
			Key:  kpath(key, keytail),
			Type: t,
			Err:  err,
		})
	}
	target.SetUint(i)
}

func parseFloat(key, keytail string, values []string, target reflect.Value) {
	t := target.Type()
	primitive(key, keytail, t, values)

	f, err := strconv.ParseFloat(values[0], t.Bits())
	if err != nil {
		panic(TypeError{
			Key:  kpath(key, k