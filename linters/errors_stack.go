package linters

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"go/types"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var ErrorsStackAnalyzer = &analysis.Analyzer{
	Name:     "errors_stack",
	Doc:      "check that errors from external packages are wrapped with stacktrace",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      runErrorsStack,
}

func init() {
	ErrorsStackAnalyzer.Flags.StringVar(&errorsStackAllowedPackagesFlag, "allowed-packages", "github.com/pkg/errors",
		"comma separated list of packages that are allowed to return errors without wrapping them with stacktrace")
	ErrorsStackAnalyzer.Flags.StringVar(&errorsStackIgnorePathRegexFlag, "ignore-path-regex", "", "regex pattern to ignore files from the analysis")
}

var (
	errorsStackAllowedPackagesFlag = ""
	errorsStackIgnorePathRegexFlag = ""
)

var (
	errorsStackAllowedPackages = []string{}
	errorsStackIgnorePathRegex *regexp.Regexp
)

func runErrorsStack(pass *analysis.Pass) (interface{}, error) {
	errorsStackAllowedPackages = strings.Split(errorsStackAllowedPackagesFlag, ",")
	if errorsStackIgnorePathRegexFlag != "" {
		errorsStackIgnorePathRegex = regexp.MustCompile(errorsStackIgnorePathRegexFlag)
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.ReturnStmt)(nil),
		(*ast.IfStmt)(nil),
		(*ast.BlockStmt)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		if errorsStackIgnorePathRegex != nil {
			if errorsStackIgnorePathRegex.MatchString(pass.Fset.Position(n.Pos()).Filename) {
				return
			}
		}

		/*
			Matches the following pattern:
			return errorableFunc()
			return obj.Error
		*/
		if retStmt, ok := n.(*ast.ReturnStmt); ok {
			for _, expr := range retStmt.Results {
				if callExpr, ok := expr.(*ast.CallExpr); ok {
					typ := pass.TypesInfo.TypeOf(callExpr)
					if isErrorType(typ) {
						if !isWrappedError(callExpr, pass) {
							var buf bytes.Buffer
							printer.Fprint(&buf, pass.Fset, callExpr)
							pass.Report(analysis.Diagnostic{
								Pos:     callExpr.Pos(),
								Message: "error returned from external package is not wrapped: return func()",
								SuggestedFixes: []analysis.SuggestedFix{
									{
										Message: "Wrap the error with stacktrace",
										TextEdits: []analysis.TextEdit{
											{
												Pos:     callExpr.Pos(),
												End:     callExpr.End(),
												NewText: []byte("errors.WithStack(" + buf.String() + ")"),
											},
										},
									},
								},
							})
						}
					}
				} else if selectorExpr, ok := expr.(*ast.SelectorExpr); ok {
					typ := pass.TypesInfo.TypeOf(selectorExpr)
					if isErrorType(typ) {
						if !identIsWrappedError(selectorExpr.Sel, pass) {
							if isErrorType(pass.TypesInfo.TypeOf(selectorExpr.X)) {
								// access err.Error is allowed
								continue
							}
							var buf bytes.Buffer
							printer.Fprint(&buf, pass.Fset, selectorExpr)
							pass.Report(analysis.Diagnostic{
								Pos:     selectorExpr.Pos(),
								Message: "error returned from external package is not wrapped: return obj.Error",
								SuggestedFixes: []analysis.SuggestedFix{
									{
										Message: "Wrap the error with stacktrace",
										TextEdits: []analysis.TextEdit{
											{
												Pos:     selectorExpr.Pos(),
												End:     selectorExpr.End(),
												NewText: []byte("errors.WithStack(" + buf.String() + ")"),
											},
										},
									},
								},
							})
						}
					}
				}
			}
		}

		/*
			Matches the following pattern:
			if err := errorableFunc(); err != nil {
				return err
			}
		*/
		if ifStmt, ok := n.(*ast.IfStmt); ok {
			if ifStmt.Init != nil && len(ifStmt.Body.List) == 1 {
				if assignStmt, ok := ifStmt.Init.(*ast.AssignStmt); ok {
					rhs, ok := assignStmt.Rhs[len(assignStmt.Rhs)-1].(*ast.CallExpr)
					if !ok {
						return
					}

					typ := pass.TypesInfo.TypeOf(rhs)
					if isErrorType(typ) {
						if !isWrappedError(rhs, pass) {
							// check if the returned error is wrapped
							if retStmt, ok := ifStmt.Body.List[0].(*ast.ReturnStmt); ok {
								for _, expr := range retStmt.Results {
									if retIdent, ok := expr.(*ast.Ident); ok {
										typ := pass.TypesInfo.TypeOf(retIdent)
										if isErrorType(typ) {
											pass.Report(analysis.Diagnostic{
												Pos:     retIdent.Pos(),
												Message: "error returned from external package is not wrapped: if err := func(); err != nil { return err }",
												SuggestedFixes: []analysis.SuggestedFix{
													{
														Message: "Wrap the error with stacktrace",
														TextEdits: []analysis.TextEdit{
															{
																Pos:     retIdent.Pos(),
																End:     retIdent.End(),
																NewText: []byte("errors.WithStack(" + retIdent.Name + ")"),
															},
														},
													},
												},
											})
										}
									}
								}
							}
						}
					}
				}
			}
		}

		/*
			Matches the following pattern:
			err := errorableFunc()
			if err != nil {
				return err
			}
		*/
		// Iterate through the statements in the block
		if blockStmt, ok := n.(*ast.BlockStmt); ok {
			for i := 1; i < len(blockStmt.List); i++ {
				prevStmt := blockStmt.List[i-1]
				currentStmt := blockStmt.List[i]

				// Step 2: Check if the next statement is the specific 'if' pattern
				ifStmt, ok := currentStmt.(*ast.IfStmt)
				if !ok || ifStmt.Init != nil || len(ifStmt.Body.List) != 1 {
					continue
				}

				// Check for error check in 'if' condition
				binaryExpr, ok := ifStmt.Cond.(*ast.BinaryExpr)
				if !ok || binaryExpr.Op != token.NEQ {
					continue
				}

				// Ensure the 'if' statement condition checks the error variable from the assignment
				ident, ok := binaryExpr.X.(*ast.Ident)
				if !ok { // || ident.Name != lhs.Name {
					continue
				}

				// Check if the 'if' body directly contains a return statement with the error
				retStmt, ok := ifStmt.Body.List[0].(*ast.ReturnStmt)
				if !ok {
					continue
				}

				// void return
				if len(retStmt.Results) == 0 {
					continue
				}

				// check if the error is returned directly
				if retIdent, ok := retStmt.Results[len(retStmt.Results)-1].(*ast.Ident); ok && retIdent.Name == ident.Name {
					// check if the error is wrapped from the last assignment
					prevStmt, ok := prevStmt.(*ast.AssignStmt)
					if !ok {
						continue
					}
					prevErr := prevStmt.Lhs[len(prevStmt.Lhs)-1]
					if prevErrIdent, ok := prevErr.(*ast.Ident); !ok || prevErrIdent.Obj.Decl != retIdent.Obj.Decl {
						// ignore complex assignments for now
						continue
					}

					callExpr, ok := prevStmt.Rhs[len(prevStmt.Rhs)-1].(*ast.CallExpr)
					if !ok {
						continue
					}
					if isWrappedError(callExpr, pass) {
						continue
					}

					// Found the pattern: err assignment followed by if check and return
					pass.Report(analysis.Diagnostic{
						Pos:     retIdent.Pos(),
						Message: "error returned from external package is not wrapped: err := func(); if err != nil { return err }",
						SuggestedFixes: []analysis.SuggestedFix{
							{
								Message: "Wrap the error with stacktrace",
								TextEdits: []analysis.TextEdit{
									{
										Pos:     retIdent.Pos(),
										End:     retIdent.End(),
										NewText: []byte("errors.WithStack(" + retIdent.Name + ")"),
									},
								},
							},
						},
					})
				}
			}
		}
	})

	return nil, nil
}

func isErrorType(t types.Type) bool {
	return types.Implements(t, types.Universe.Lookup("error").Type().Underlying().(*types.Interface))
}

func isWrappedError(call *ast.CallExpr, pass *analysis.Pass) bool {
	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		// normal err := func() calls
		if x, ok := fun.X.(*ast.Ident); ok {
			if identIsWrappedError(x, pass) {
				return true
			}
		}
		if identIsWrappedError(fun.Sel, pass) {
			return true
		}
	case *ast.CallExpr:
		// err := namedFunc()() calls
		if fun, ok := call.Fun.(*ast.CallExpr); ok {
			return isWrappedError(fun, pass)
		}
	case *ast.IndexExpr:
		// generic[T]() calls
		if x, ok := fun.X.(*ast.Ident); ok {
			if identIsWrappedError(x, pass) {
				return true
			}
		}
	case *ast.FuncLit:
		// err := func() error { ... }() calls
		return true
	case *ast.Ident:
		// intra-package calls
		if identIsWrappedError(fun, pass) {
			return true
		}
	}
	return false
}

/*
An identifier is considered to be a wrapped error the package it belongs to is in the list of allowed packages
*/
func identIsWrappedError(id *ast.Ident, pass *analysis.Pass) bool {
	if obj := pass.TypesInfo.ObjectOf(id); obj != nil {
		pkgPath := ""
		// in case of named imports
		if pkgName, ok := obj.(*types.PkgName); ok {
			pkgPath = pkgName.Imported().Path()
		} else if pkg := obj.Pkg(); pkg != nil {
			pkgPath = pkg.Path()
		}

		if pkgPath != "" {
			for _, allowedPkg := range errorsStackAllowedPackages {
				if strings.HasPrefix(pkgPath, allowedPkg) {
					return true
				}
			}
		}
	}
	return false
}
