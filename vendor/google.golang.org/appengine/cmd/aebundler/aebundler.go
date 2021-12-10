// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Program aebundler turns a Go app into a fully self-contained tar file.
// The app and its subdirectories (if any) are placed under "."
// and the dependencies from $GOPATH are placed under ./_gopath/src.
// A main func is synthesized if one does not exist.
//
// A sample Dockerfile to be used with this bundler could look like this:
//     FROM gcr.io/google-appengine/go-compat
//     ADD . /app
//     RUN GOPATH=/app/_gopath go build -tags appenginevm -o /app/_ah/exe
package main

import (
	"archive/tar"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	output  = flag.String("o", "", "name of output tar file or '-' for stdout")
	rootDir = flag.String("root", ".", "directory name of application root")
	vm      = flag.Bool("vm", true, `bundle an app for App Engine "flexible environment"`)

	skipFiles = map[string]bool{
		".git":        true,
		".gitconfig":  true,
		".hg":         true,
		".travis.yml": true,
	}
)

const (
	newMain = `package main
import "google.golang.org/appengine"
func main() {
	appengine.Main()
}
`
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\t%s -o <file.tar|->\tBundle app to named tar file or stdout\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\noptional arguments:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	var tags []string
	if *vm {
		tags = append(tags, "appenginevm")
	} else {
		tags = append(tags, "appengine")
	}

	tarFile := *output
	if tarFile == "" {
		usage()
		errorf("Required -o flag not specified.")
	}

	app, err := analyze(tags)
	if err != nil {
		errorf("Error analyzing app: %v", err)
	}
	if err := app.bundle(tarFile); err != nil {
		errorf("Unable to bundle app: %v", err)
	}
}

// errorf prints the error message and exits.
func errorf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "aebundler: "+format+"\n", a...)
	os.Exit(1)
}

type app struct {
	hasMain  bool
	appFiles []string
	imports  map[string]string
}

// analyze checks the app for building with the given build tags and returns hasMain,
// app files, and a map of full directory import names to original import names.
func analyze(tags []string) (*app, error) {
	ctxt := buildContext(tags)
	hasMain, appFiles, err := checkMain(ctxt)
	if err != nil {
		return nil, err
	}
	gopath := filepath.SplitList(ctxt.GOPATH)
	im, err := imports(ctxt, *rootDir, gopath)
	return &app{
		hasMain:  hasMain,
		appFiles: appFiles,
		imports:  im,
	}, err
}

// buildContext returns the context for building the source.
func buildContext(tags []string) *build.Context {
	return &build.Context{
		GOARCH:    build.Default.GOARCH,
		GOOS:      build.Default.GOOS,
		GOROOT:    build.Default.GOROOT,
		GOPATH:    build.Default.GOPATH,
		Compiler:  build.Default.Compiler,
		BuildTags: append(build.Default.BuildTags, tags...),
	}
}

// bundle bundles the app into the named tarFile ("-"==stdout).
func (s *app) bundle(tarFile string) (err error) {
	var out io.Writer
	if tarFile == "-" {
		out = os.Stdout
	} else {
		f, err := os.Create(tarFile)
		if err != nil {
			return err
		}
		defer func() {
			if cerr := f.Close(); err == nil {
				err = cerr
			}
		}()
		out = f
	}
	tw := tar.NewWriter(out)

	for srcDir, importName := range s.imports {
		dstDir := "_gopath/src/" + importName
		if err = copyTree(tw, dstDir, srcDir); err != nil {
			return fmt.Errorf("unable to copy directory %v to %v: %v", srcDir, dstDir, err)
		}
	}
	if err := copyTree(tw, ".", *rootDir); err != nil {
		return fmt.Errorf("unable to copy root directory to /app: %v", err)
	}
	if !s.hasMain {
		if err := synthesizeMain(tw, s.appFiles); err != nil {
			return fmt.Errorf("unable to synthesize new main func: %v", err)
		}
	}

	if err := tw.Close(); err != nil {
		return fmt.Errorf("unable to close tar 