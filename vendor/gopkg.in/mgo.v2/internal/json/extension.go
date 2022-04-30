package json

import (
	"reflect"
)

// Extension holds a set of additional rules to be used when unmarshaling
// strict JSON or JSON-like content.
type Extension struct {
	funcs  map[string]funcExt
	consts map[string]interface{}
	keyed  map[string]func([]byte) (interface{}, error)
	encode map[reflect.Type]func(v interface{}) ([]byte, error)

	unquotedKeys   bool
	trailingCommas bool
}

type funcExt struct {
	key  string
	args []string
}

// Extend changes the decoder behavior to consider the provided extension.
func (dec *Decoder) Extend(ext *Extension) { dec.d.ext = *ext }

// Extend changes the encoder behavior to consider the provided extension.
func (enc *Encoder) Extend(ext *Extension) { enc.ext = *ext }

// Extend includes in e the extensions defined in ext.
func (e *Extension) Extend(ext *Extension) {
	for name, fext := range ext.funcs {
		e.DecodeFunc(name, fext.key, fext.args...)
	}
	for name, value := range ext.consts {
		e.DecodeConst(name, value)
	}
	for key, decode := range ext.keyed {
		e.DecodeKeyed(key, decode)
	}
	for typ, encode := range ext.encode {
		if e.encode == nil {
			e.encode = make(map[reflect.Type]func(v interface{}) ([]byte, error))
		}
		e.encode[typ] = encode
	}
}

// DecodeFunc defines a function call that may be observed inside JSON content.
// A function with the provided name will be unmarshaled as the document
// {key: {args[0]: ..., args[N]: ...}}.
func (e *Extension) DecodeFunc(name string, key string, args ...string) {
	if e.func