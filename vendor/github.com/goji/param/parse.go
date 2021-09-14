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
// parser is responsible for, for inst