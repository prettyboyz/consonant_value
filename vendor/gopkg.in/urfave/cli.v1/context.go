package cli

import (
	"errors"
	"flag"
	"reflect"
	"strings"
	"syscall"
)

// Context is a type that is passed through to
// each Ha