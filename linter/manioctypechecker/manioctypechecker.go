package manioctypechecker

import (
	"fmt"
	"go/ast"
	"go/types"
	"reflect"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/imports"
)

const doc = "A static type checker for github.com/fuzmish/manioc"

//nolint:gochecknoglobals
var Analyzer = &analysis.Analyzer{
	Name: "manioctypechecker",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func typeCategoryName(t types.Type) string {
	switch info := t.(type) {
	case *types.Basic:
		return info.Name()
	case *types.Named:
		return typeCategoryName(info.Underlying())
	default:
		return reflect.TypeOf(info).Elem().Name()
	}
}

func checkManiocPackage(pass *analysis.Pass, x ast.Expr) bool {
	//nolint:varnamelen
	ident, ok := x.(*ast.Ident)
	if !ok {
		return false
	}
	use, ok := pass.TypesInfo.Uses[ident]
	if !ok {
		return false
	}
	pkg, ok := use.(*types.PkgName)
	if !ok {
		return false
	}
	path := imports.VendorlessPath(pkg.Imported().Path())
	return path == "github.com/fuzmish/manioc"
}

func checkManiocRegisterConstructorTypeParameters(
	pass *analysis.Pass,
	expr ast.Expr,
	tTInterface types.Type,
	tTConstructor types.Type,
) {
	// check ctor signature
	tSignature, ok := tTConstructor.(*types.Signature)
	if !ok {
		pass.Report(analysis.Diagnostic{
			Pos: expr.Pos(),
			Message: fmt.Sprintf(
				"The argument type should be a function type, but `%v` is given",
				typeCategoryName(tTConstructor),
			),
		})
		return
	}
	// check ctor return type
	funcRet := tSignature.Results()
	rLen := funcRet.Len()
	if rLen < 1 || rLen > 2 {
		pass.Report(analysis.Diagnostic{
			Pos:     expr.Pos(),
			Message: "The number of function return values should be either one or two",
		})
		return
	}
	// check the first return type
	tRet := funcRet.At(0).Type()
	if !types.AssignableTo(tRet, tTInterface) {
		pass.Report(analysis.Diagnostic{
			Pos: expr.Pos(),
			Message: fmt.Sprintf(
				"The type of the first return value `%v` is not assignable to `%v`",
				tRet,
				tTInterface,
			),
		})
		return
	}
	// check the second return type
	if rLen != 1 {
		tRetError := funcRet.At(1).Type()
		tError := types.Universe.Lookup("error").Type()
		if !types.AssignableTo(tRetError, tError) {
			pass.Report(analysis.Diagnostic{
				Pos: expr.Pos(),
				Message: fmt.Sprintf(
					"The type of the second return value should be `error`, but `%v` is given",
					tRetError,
				),
			})
		}
	}
}

func checkManiocRegisterTypeParameters(
	pass *analysis.Pass,
	expr ast.Expr,
	tTInterface types.Type,
	tTImplementation types.Type,
) {
	// if TInterface is an interface type
	tPtrTImplementation := tTImplementation
	if _, ok := tTInterface.Underlying().(*types.Interface); ok {
		// and TImplementation is not a pointer type
		if _, ok := tTImplementation.(*types.Pointer); !ok {
			// then, replace TImplementation with its pointer type implicitly
			tPtrTImplementation = types.NewPointer(tTImplementation)
		}
	}

	// check if TImplementation is not an interface type
	tElmTImplementation := tPtrTImplementation
	if tPtr, ok := tElmTImplementation.(*types.Pointer); ok {
		tElmTImplementation = tPtr.Elem()
	}
	if _, ok := tElmTImplementation.Underlying().(*types.Interface); ok {
		pass.Report(analysis.Diagnostic{
			Pos: expr.Pos(),
			Message: fmt.Sprintf(
				"The implementation type `%v` should not be an interface",
				tElmTImplementation,
			),
		})
		return
	}

	// check if TImplementation is assignable to TInterface
	if !types.AssignableTo(tPtrTImplementation, tTInterface) {
		pass.Report(analysis.Diagnostic{
			Pos: expr.Pos(),
			Message: fmt.Sprintf(
				"`%v` is not assignable to `%v`",
				tTImplementation,
				tTInterface,
			),
		})
	}
}

func checkManiocRegisterConstructorCall(pass *analysis.Pass, call *ast.CallExpr) {
	// skip calls without args since RegisterConstructor takes at least one args (ctor)
	if len(call.Args) < 1 {
		return
	}
	// check selector
	var fun ast.Expr
	var index ast.Expr
	if ile, ok := call.Fun.(*ast.IndexListExpr); ok {
		// RegisterConstructor[TInterface, TConstructor](...)
		fun = ile.X
		if len(ile.Indices) < 1 {
			return
		}
		index = ile.Indices[0]
	} else if ie, ok := call.Fun.(*ast.IndexExpr); ok {
		// RegisterConstructor[TInterface](ctor, ...)
		fun = ie.X
		index = ie.Index
	} else {
		return
	}
	selector, ok := fun.(*ast.SelectorExpr)
	if !ok {
		return
	}
	if !checkManiocPackage(pass, selector.X) {
		return
	}
	switch selector.Sel.Name {
	case "RegisterConstructor",
		"RegisterSingletonConstructor",
		"RegisterScopedConstructor",
		"RegisterTransientConstructor",
		"ResolveFunction",
		"MustResolveFunction":
		tTInterface := pass.TypesInfo.TypeOf(index)
		tTConstructor := pass.TypesInfo.TypeOf(call.Args[0])
		checkManiocRegisterConstructorTypeParameters(pass, call, tTInterface, tTConstructor)
	}
}

func checkManiocRegisterFunction(pass *analysis.Pass, indexList *ast.IndexListExpr) {
	// since the Register function takes two type parameters, skip other cases
	if len(indexList.Indices) != 2 { //nolint:gomnd
		return
	}
	// check selector
	selector, ok := indexList.X.(*ast.SelectorExpr)
	if !ok {
		return
	}
	if !checkManiocPackage(pass, selector.X) {
		return
	}
	switch selector.Sel.Name {
	case "Register",
		"RegisterSingleton",
		"RegisterScoped",
		"RegisterTransient":
		tTInterface := pass.TypesInfo.TypeOf(indexList.Indices[0])
		tTImplementation := pass.TypesInfo.TypeOf(indexList.Indices[1])
		checkManiocRegisterTypeParameters(pass, indexList, tTInterface, tTImplementation)
	case "RegisterConstructor",
		"RegisterSingletonConstructor",
		"RegisterScopedConstructor",
		"RegisterTransientConstructor",
		"ResolveFunction",
		"MustResolveFunction":
		tTInterface := pass.TypesInfo.TypeOf(indexList.Indices[0])
		tTConstructor := pass.TypesInfo.TypeOf(indexList.Indices[1])
		checkManiocRegisterConstructorTypeParameters(pass, indexList, tTInterface, tTConstructor)
	}
}

func run(pass *analysis.Pass) (any, error) {
	//nolint:forcetypeassert
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
		(*ast.IndexListExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
		case *ast.CallExpr:
			checkManiocRegisterConstructorCall(pass, n)
		case *ast.IndexListExpr:
			checkManiocRegisterFunction(pass, n)
		}
	})

	return nil, nil
}
