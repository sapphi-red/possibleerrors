package avoidaccesslen

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "avoidaccesslen is finds `arr[len(arr)]` which occurs index out of range"

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "avoidaccesslen",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

var lenObj = types.Universe.Lookup("len")

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		new(ast.IndexExpr),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		indexAccess := n.(*ast.IndexExpr)
		arrObj := extractArrObj(pass, indexAccess)
		if arrObj == nil {
			return
		}

		index := indexAccess.Index
		indexCall, _ := index.(*ast.CallExpr)
		if indexCall == nil {
			return
		}

		indexFuncObj := extractIndexFuncObj(pass, indexCall)
		if lenObj != indexFuncObj {
			return
		}

		argObj := extractIndexFuncArgObj(pass, indexCall)
		if arrObj != argObj {
			return
		}

		fix := analysis.SuggestedFix{
			Message: "Add ` - 1` after `len(arr)`",
			TextEdits: []analysis.TextEdit{{
				Pos:     indexCall.End(),
				End:     indexCall.End(),
				NewText: []byte("-1"),
			}},
		}

		pass.Report(analysis.Diagnostic{
			Pos:            n.Pos(),
			End:            n.End(),
			Message:        "Will occur index out of range",
			SuggestedFixes: []analysis.SuggestedFix{fix},
		})
	})

	return nil, nil
}

func extractArrObj(pass *analysis.Pass, indexAccess *ast.IndexExpr) types.Object {
	arr := indexAccess.X
	arrType := pass.TypesInfo.Types[arr]

	T := arrType.Type.Underlying()
	_, ok := T.(*types.Slice)
	if !ok {
		return nil
	}

	arrRightId, _ := arr.(*ast.Ident)
	if arrRightId == nil {
		arrSelector, _ := arr.(*ast.SelectorExpr)
		if arrSelector == nil {
			return nil
		}
		arrRightId = arrSelector.Sel
	}

	return pass.TypesInfo.Uses[arrRightId]
}

func extractIndexFuncObj(pass *analysis.Pass, indexCall *ast.CallExpr) types.Object {
	indexFunc := indexCall.Fun
	indexFuncId, _ := indexFunc.(*ast.Ident)
	if indexFuncId == nil {
		return nil
	}

	return pass.TypesInfo.Uses[indexFuncId]
}

func extractIndexFuncArgObj(pass *analysis.Pass, indexCall *ast.CallExpr) types.Object {
	args := indexCall.Args
	if len(args) != 1 {
		return nil
	}

	argId := args[0].(*ast.Ident)
	return pass.TypesInfo.Uses[argId]
}
