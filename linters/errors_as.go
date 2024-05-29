package linters

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/types/typeutil"
)

// ErrorsAsAnalyzer adapted from go/x/tools
// https://cs.opensource.google/go/x/tools/+/refs/tags/v0.14.0:go/analysis/passes/errorsas/errorsas.go
var ErrorsAsAnalyzer = &analysis.Analyzer{
	Name:     "errorsAs",
	Doc:      "reports illegal calls to errors.As and pkg/errors.As",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      runErrorsAs,
}

var errorsPackages = []string{
	"errors",
	"github.com/pkg/errors",
}

func runErrorsAs(pass *analysis.Pass) (interface{}, error) {
	switch pass.Pkg.Path() {
	case "errors", "errors_test":
		// These packages know how to use their own APIs.
		// Sometimes they are testing what happens to incorrect programs.
		return nil, nil
	}

	importsOne := false
	for _, pkg := range errorsPackages {
		if !Imports(pass.Pkg, pkg) {
			continue
		} else {
			importsOne = true
			break
		}
	}
	// does not import an errors.As provider package
	if !importsOne {
		return nil, nil
	}

	passInspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	passInspector.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		fn := typeutil.StaticCallee(pass.TypesInfo, call)
		if fn == nil {
			return // not a static call
		}
		if len(call.Args) < 2 {
			return // not enough arguments, e.g. called with return values of another function
		}
		callsAny := false
		for _, pkg := range errorsPackages {
			if fn.FullName() != fmt.Sprintf("%s.As", pkg) {
				continue
			}
			callsAny = true
			break
		}
		if !callsAny {
			return
		}
		if err := checkAsTarget(pass, call.Args[1]); err != nil {
			pass.ReportRangef(call, "%v", err)
		}
	})
	return nil, nil
}

var errorType = types.Universe.Lookup("error").Type()

// checkAsTarget reports an error if the second argument to errors.As is invalid.
func checkAsTarget(pass *analysis.Pass, e ast.Expr) error {
	t := pass.TypesInfo.Types[e].Type
	if it, ok := t.Underlying().(*types.Interface); ok && it.NumMethods() == 0 {
		// A target of interface{} is always allowed, since it often indicates
		// a value forwarded from another source.
		return nil
	}
	pt, ok := t.Underlying().(*types.Pointer)
	if !ok {
		return errors.New("second argument to errors.As must be a non-nil pointer to either a type that implements error, or to any interface type")
	}
	if pt.Elem() == errorType {
		return errors.New("second argument to errors.As should not be *error")
	}
	_, ok = pt.Elem().Underlying().(*types.Interface)
	if ok || types.Implements(pt.Elem(), errorType.Underlying().(*types.Interface)) {
		return nil
	}
	return errors.New("second argument to errors.As must be a non-nil pointer to either a type that implements error, or to any interface type")
}

// Imports a utility method that returns true if path is imported by pkg.
// borrowed from https://cs.opensource.google/go/x/tools/+/master:go/analysis/passes/internal/analysisutil/util.go;l=102
func Imports(pkg *types.Package, path string) bool {
	for _, imp := range pkg.Imports() {
		if imp.Path() == path {
			return true
		}
	}
	return false
}
