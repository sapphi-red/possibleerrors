package avoidaccesslen

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "avoidaccesslen is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "avoidaccesslen",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		new(ast.IndexExpr),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		indexAccess := n.(*ast.IndexExpr)
		arr := indexAccess.X
		arrType := pass.TypesInfo.Types[arr]

		T := arrType.Type.Underlying()
		_, ok := T.(*types.Slice)
		if !ok {
			return
		}

		arrRightId, _ := arr.(*ast.Ident)
		if arrRightId == nil {
			arrSelector, _ := arr.(*ast.SelectorExpr)
			if arrSelector == nil {
				return
			}
			arrRightId = arrSelector.Sel
		}

		arrObj := pass.TypesInfo.Uses[arrRightId]

		index := indexAccess.Index
		indexCall, _ := index.(*ast.CallExpr)
		if indexCall == nil {
			return
		}

		indexFunc := indexCall.Fun
		indexFuncId, _ := indexFunc.(*ast.Ident)
		if indexFuncId == nil {
			return
		}

		lenObj := types.Universe.Lookup("len")
		indexFuncObj := pass.TypesInfo.Uses[indexFuncId]
		if lenObj != indexFuncObj {
			return
		}

		args := indexCall.Args
		if len(args) != 1 {
			return
		}

		argId := args[0].(*ast.Ident)
		argObj := pass.TypesInfo.Uses[argId]
		if arrObj != argObj {
			return
		}
		pass.Reportf(n.Pos(), "Likely cause index out of range")
	})

	return nil, nil
}
