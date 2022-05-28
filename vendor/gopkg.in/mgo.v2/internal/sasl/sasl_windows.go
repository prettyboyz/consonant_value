package sasl

// #include "sasl_windows.h"
import "C"

import (
	"fmt"
	"strings"
	"sync"
	"unsafe"
)

type saslStepper interface {
	Step(serverData []byte) (clientData []byte, done bool, err error)
	Close()
}

type saslSession struct {
	// Credentials
	mech          string
	service       string
	host          string
	userPlusRealm st