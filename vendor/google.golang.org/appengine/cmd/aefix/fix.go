
// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)

type fix struct {
	name string
	date string // date that fix was introduced, in YYYY-MM-DD format
	f    func(*ast.File) bool
	desc string
}

// main runs sort.Sort(byName(fixes)) before printing list of fixes.
type byName []fix

func (f byName) Len() int           { return len(f) }
func (f byName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f byName) Less(i, j int) bool { return f[i].name < f[j].name }

// main runs sort.Sort(byDate(fixes)) before applying fixes.
type byDate []fix

func (f byDate) Len() int           { return len(f) }
func (f byDate) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }
func (f byDate) Less(i, j int) bool { return f[i].date < f[j].date }

var fixes []fix

func register(f fix) {
	fixes = append(fixes, f)
}

// walk traverses the AST x, calling visit(y) for each node y in the tree but
// also with a pointer to each ast.Expr, ast.Stmt, and *ast.BlockStmt,
// in a bottom-up traversal.
func walk(x interface{}, visit func(interface{})) {
	walkBeforeAfter(x, nop, visit)
}

func nop(interface{}) {}

// walkBeforeAfter is like walk but calls before(x) before traversing
// x's children and after(x) afterward.
func walkBeforeAfter(x interface{}, before, after func(interface{})) {
	before(x)

	switch n := x.(type) {
	default:
		panic(fmt.Errorf("unexpected type %T in walkBeforeAfter", x))

	case nil:

	// pointers to interfaces
	case *ast.Decl:
		walkBeforeAfter(*n, before, after)
	case *ast.Expr:
		walkBeforeAfter(*n, before, after)
	case *ast.Spec:
		walkBeforeAfter(*n, before, after)
	case *ast.Stmt:
		walkBeforeAfter(*n, before, after)

	// pointers to struct pointers
	case **ast.BlockStmt:
		walkBeforeAfter(*n, before, after)
	case **ast.CallExpr:
		walkBeforeAfter(*n, before, after)
	case **ast.FieldList:
		walkBeforeAfter(*n, before, after)
	case **ast.FuncType:
		walkBeforeAfter(*n, before, after)
	case **ast.Ident:
		walkBeforeAfter(*n, before, after)
	case **ast.BasicLit:
		walkBeforeAfter(*n, before, after)

	// pointers to slices
	case *[]ast.Decl:
		walkBeforeAfter(*n, before, after)
	case *[]ast.Expr:
		walkBeforeAfter(*n, before, after)
	case *[]*ast.File:
		walkBeforeAfter(*n, before, after)
	case *[]*ast.Ident:
		walkBeforeAfter(*n, before, after)
	case *[]ast.Spec:
		walkBeforeAfter(*n, before, after)
	case *[]ast.Stmt:
		walkBeforeAfter(*n, before, after)

	// These are ordered and grouped to match ../../pkg/go/ast/ast.go
	case *ast.Field:
		walkBeforeAfter(&n.Names, before, after)
		walkBeforeAfter(&n.Type, before, after)
		walkBeforeAfter(&n.Tag, before, after)
	case *ast.FieldList:
		for _, field := range n.List {
			walkBeforeAfter(field, before, after)
		}
	case *ast.BadExpr:
	case *ast.Ident:
	case *ast.Ellipsis:
		walkBeforeAfter(&n.Elt, before, after)
	case *ast.BasicLit:
	case *ast.FuncLit:
		walkBeforeAfter(&n.Type, before, after)
		walkBeforeAfter(&n.Body, before, after)
	case *ast.CompositeLit:
		walkBeforeAfter(&n.Type, before, after)
		walkBeforeAfter(&n.Elts, before, after)
	case *ast.ParenExpr:
		walkBeforeAfter(&n.X, before, after)
	case *ast.SelectorExpr:
		walkBeforeAfter(&n.X, before, after)
	case *ast.IndexExpr:
		walkBeforeAfter(&n.X, before, after)
		walkBeforeAfter(&n.Index, before, after)
	case *ast.SliceExpr:
		walkBeforeAfter(&n.X, before, after)
		if n.Low != nil {
			walkBeforeAfter(&n.Low, before, after)
		}
		if n.High != nil {
			walkBeforeAfter(&n.High, before, after)
		}
	case *ast.TypeAssertExpr:
		walkBeforeAfter(&n.X, before, after)
		walkBeforeAfter(&n.Type, before, after)
	case *ast.CallExpr:
		walkBeforeAfter(&n.Fun, before, after)
		walkBeforeAfter(&n.Args, before, after)
	case *ast.StarExpr:
		walkBeforeAfter(&n.X, before, after)
	case *ast.UnaryExpr:
		walkBeforeAfter(&n.X, before, after)
	case *ast.BinaryExpr:
		walkBeforeAfter(&n.X, before, after)
		walkBeforeAfter(&n.Y, before, after)
	case *ast.KeyValueExpr:
		walkBeforeAfter(&n.Key, before, after)
		walkBeforeAfter(&n.Value, before, after)

	case *ast.ArrayType:
		walkBeforeAfter(&n.Len, before, after)
		walkBeforeAfter(&n.Elt, before, after)
	case *ast.StructType:
		walkBeforeAfter(&n.Fields, before, after)
	case *ast.FuncType:
		walkBeforeAfter(&n.Params, before, after)
		if n.Results != nil {
			walkBeforeAfter(&n.Results, before, after)
		}
	case *ast.InterfaceType:
		walkBeforeAfter(&n.Methods, before, after)
	case *ast.MapType:
		walkBeforeAfter(&n.Key, before, after)
		walkBeforeAfter(&n.Value, before, after)
	case *ast.ChanType:
		walkBeforeAfter(&n.Value, before, after)

	case *ast.BadStmt:
	case *ast.DeclStmt:
		walkBeforeAfter(&n.Decl, before, after)
	case *ast.EmptyStmt:
	case *ast.LabeledStmt:
		walkBeforeAfter(&n.Stmt, before, after)
	case *ast.ExprStmt:
		walkBeforeAfter(&n.X, before, after)
	case *ast.SendStmt:
		walkBeforeAfter(&n.Chan, before, after)
		walkBeforeAfter(&n.Value, before, after)
	case *ast.IncDecStmt:
		walkBeforeAfter(&n.X, before, after)
	case *ast.AssignStmt:
		walkBeforeAfter(&n.Lhs, before, after)
		walkBeforeAfter(&n.Rhs, before, after)
	case *ast.GoStmt:
		walkBeforeAfter(&n.Call, before, after)
	case *ast.DeferStmt:
		walkBeforeAfter(&n.Call, before, after)
	case *ast.ReturnStmt:
		walkBeforeAfter(&n.Results, before, after)
	case *ast.BranchStmt:
	case *ast.BlockStmt:
		walkBeforeAfter(&n.List, before, after)
	case *ast.IfStmt:
		walkBeforeAfter(&n.Init, before, after)
		walkBeforeAfter(&n.Cond, before, after)
		walkBeforeAfter(&n.Body, before, after)
		walkBeforeAfter(&n.Else, before, after)
	case *ast.CaseClause:
		walkBeforeAfter(&n.List, before, after)
		walkBeforeAfter(&n.Body, before, after)
	case *ast.SwitchStmt:
		walkBeforeAfter(&n.Init, before, after)
		walkBeforeAfter(&n.Tag, before, after)
		walkBeforeAfter(&n.Body, before, after)
	case *ast.TypeSwitchStmt:
		walkBeforeAfter(&n.Init, before, after)
		walkBeforeAfter(&n.Assign, before, after)
		walkBeforeAfter(&n.Body, before, after)
	case *ast.CommClause:
		walkBeforeAfter(&n.Comm, before, after)
		walkBeforeAfter(&n.Body, before, after)
	case *ast.SelectStmt:
		walkBeforeAfter(&n.Body, before, after)
	case *ast.ForStmt:
		walkBeforeAfter(&n.Init, before, after)
		walkBeforeAfter(&n.Cond, before, after)
		walkBeforeAfter(&n.Post, before, after)
		walkBeforeAfter(&n.Body, before, after)
	case *ast.RangeStmt:
		walkBeforeAfter(&n.Key, before, after)
		walkBeforeAfter(&n.Value, before, after)
		walkBeforeAfter(&n.X, before, after)
		walkBeforeAfter(&n.Body, before, after)

	case *ast.ImportSpec:
	case *ast.ValueSpec:
		walkBeforeAfter(&n.Type, before, after)
		walkBeforeAfter(&n.Values, before, after)
		walkBeforeAfter(&n.Names, before, after)
	case *ast.TypeSpec:
		walkBeforeAfter(&n.Type, before, after)

	case *ast.BadDecl:
	case *ast.GenDecl:
		walkBeforeAfter(&n.Specs, before, after)
	case *ast.FuncDecl:
		if n.Recv != nil {
			walkBeforeAfter(&n.Recv, before, after)
		}
		walkBeforeAfter(&n.Type, before, after)
		if n.Body != nil {
			walkBeforeAfter(&n.Body, before, after)
		}

	case *ast.File:
		walkBeforeAfter(&n.Decls, before, after)

	case *ast.Package:
		walkBeforeAfter(&n.Files, before, after)

	case []*ast.File:
		for i := range n {
			walkBeforeAfter(&n[i], before, after)
		}
	case []ast.Decl:
		for i := range n {
			walkBeforeAfter(&n[i], before, after)
		}
	case []ast.Expr:
		for i := range n {
			walkBeforeAfter(&n[i], before, after)
		}
	case []*ast.Ident:
		for i := range n {
			walkBeforeAfter(&n[i], before, after)
		}
	case []ast.Stmt:
		for i := range n {
			walkBeforeAfter(&n[i], before, after)
		}
	case []ast.Spec:
		for i := range n {
			walkBeforeAfter(&n[i], before, after)
		}
	}
	after(x)
}

// imports returns true if f imports path.
func imports(f *ast.File, path string) bool {
	return importSpec(f, path) != nil
}

// importSpec returns the import spec if f imports path,
// or nil otherwise.
func importSpec(f *ast.File, path string) *ast.ImportSpec {
	for _, s := range f.Imports {
		if importPath(s) == path {
			return s
		}
	}
	return nil
}

// importPath returns the unquoted import path of s,
// or "" if the path is not properly quoted.
func importPath(s *ast.ImportSpec) string {
	t, err := strconv.Unquote(s.Path.Value)
	if err == nil {
		return t
	}
	return ""
}

// declImports reports whether gen contains an import of path.
func declImports(gen *ast.GenDecl, path string) bool {
	if gen.Tok != token.IMPORT {
		return false
	}
	for _, spec := range gen.Specs {
		impspec := spec.(*ast.ImportSpec)
		if importPath(impspec) == path {
			return true
		}
	}
	return false
}

// isPkgDot returns true if t is the expression "pkg.name"
// where pkg is an imported identifier.
func isPkgDot(t ast.Expr, pkg, name string) bool {
	sel, ok := t.(*ast.SelectorExpr)
	return ok && isTopName(sel.X, pkg) && sel.Sel.String() == name
}

// isPtrPkgDot returns true if f is the expression "*pkg.name"
// where pkg is an imported identifier.
func isPtrPkgDot(t ast.Expr, pkg, name string) bool {
	ptr, ok := t.(*ast.StarExpr)
	return ok && isPkgDot(ptr.X, pkg, name)
}

// isTopName returns true if n is a top-level unresolved identifier with the given name.
func isTopName(n ast.Expr, name string) bool {
	id, ok := n.(*ast.Ident)
	return ok && id.Name == name && id.Obj == nil
}

// isName returns true if n is an identifier with the given name.
func isName(n ast.Expr, name string) bool {
	id, ok := n.(*ast.Ident)
	return ok && id.String() == name
}

// isCall returns true if t is a call to pkg.name.
func isCall(t ast.Expr, pkg, name string) bool {
	call, ok := t.(*ast.CallExpr)
	return ok && isPkgDot(call.Fun, pkg, name)
}