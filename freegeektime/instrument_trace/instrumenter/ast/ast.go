package ast

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
)

type instrumenter struct {
	tranceImport string
	trancePkg    string
	tranceFunc   string
}

func New(tranceImport, trancePkg, tranceFunc string) *instrumenter {
	return &instrumenter{
		tranceImport, trancePkg, tranceFunc,
	}
}
func hasFuncDecl(f *ast.File) bool {
	if len(f.Decls) == 0 {
		return false
	}

	for _, decl := range f.Decls {
		if _, ok := decl.(*ast.FuncDecl); ok {
			return true
		}

	}
	return false
}

func (a instrumenter) Instrument(filename string) ([]byte, error) {
	fset := token.NewFileSet()
	curAst, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("error parsing %s:%w", filename, err)
	}
	if !hasFuncDecl(curAst) {
		return nil, nil
	}

	astutil.AddImport(fset, curAst, a.tranceImport)

	a.addDeferTranceIntoFuncDecls(curAst)

	buf := &bytes.Buffer{}

	if err := format.Node(buf, fset, curAst); err != nil {
		return nil, fmt.Errorf("error formatting new code:%w", err)
	}

	return buf.Bytes(), nil

}

func (a instrumenter) addDeferTranceIntoFuncDecls(f *ast.File) {
	for _, decl := range f.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok {
			a.addDeferStmt(fd)
		}

	}
}

func (a instrumenter) addDeferStmt(fd *ast.FuncDecl) (added bool) {
	stmts := fd.Body.List
	for _, stmt := range stmts {
		ds, ok := stmt.(*ast.DeferStmt)
		if !ok {
			continue
		}
		ce, ok := ds.Call.Fun.(*ast.CallExpr)
		if !ok {
			continue
		}
		se, ok := ce.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		x, ok := se.X.(*ast.Ident)
		if !ok {
			continue
		}
		if (x.Name == a.trancePkg) && (se.Sel.Name == a.tranceFunc) {
			return false
		}

	}

	ds := &ast.DeferStmt{
		Call: &ast.CallExpr{
			Fun: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: a.trancePkg,
					},
					Sel: &ast.Ident{
						Name: a.tranceFunc,
					},
				},
			},
		},
	}
	newList := make([]ast.Stmt, len(stmts)+1)
	copy(newList[1:], stmts)
	newList[0] = ds
	fd.Body.List = newList
	return true
}
