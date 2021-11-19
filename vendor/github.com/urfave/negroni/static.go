package negroni

import (
	"net/http"
	"path"
	"strings"
)

// Static is a middleware handler that serves static files in the given
// directory/filesystem. If the file does not exist on the filesystem, it
// passes along to the next middleware in the chain. If you desire "fileserver"
// type behavior where it returns a 404 for unfound files, you should consider
// using http.FileServer from the Go stdlib.
type Static struct {
	// Dir is the directory to serve static files from
	Dir http.FileSystem
	// Prefix is the optiona