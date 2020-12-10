// Command toml-test-decoder satisfies the toml-test interface for testing
// TOML decoders. Namely, it accepts TOML on stdin and outputs JSON on stdout.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/BurntSushi/toml"
)

func init() {
	log.SetFlags(0)

	flag.Usage = usage
	flag.Parse()
}

func usage() {
	log.Printf("Usage: %s < toml-file\n", path.Base(os.Args[0]))
	flag.PrintDefaults()

	os.Exit(1)
}

func main