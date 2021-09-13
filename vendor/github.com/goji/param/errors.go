package param

import (
	"fmt"
	"reflect"
)

// TypeError is an error type returned when param has difficulty deserializing a
// parameter value.
type TypeError struct {
	// The key that was in error.
	Key string
	// The type that was expected.
	Type reflect.Type
	// The underlying error produced as part of the deserialization process,
	// if one exists.
	Err error
}

func (t TypeError) Error() string {
	return fmt.Sprintf("param: error parsing key %q as %v: %v", t.Key, t.Type,
		t.Err)
}

// SingletonError is an error type returned when a parameter is passed multiple
// times when only a single value is expected. For example, for a struct with
// integer field "foo", "foo=1&foo=2" will return a SingletonError with key
// "foo".
type SingletonError struct {
	// The key that was in error.
	Key string
	// The type that was expected for that key.
	Type reflect.Type
	// The list of values that were provided for that key.
	Values []string
}

func (s SingletonError) Error() string {
	return fmt.Sprintf("param: error parsing key %q: expected single "+
		"value but was given %d: %v", s.Key, len(s.Values), s.Values)
}

// NestingError is an error type returned when a key is nested when the target
// type does not support nesting of the given type. For example, deserializing
// the parameter key "anint[foo]" into a struct that defines an integer param
// "anint" will produce a NestingError with key "anint" and nesting "[foo]".
type NestingError struct {
	// The portion of the key that was correctly parsed 