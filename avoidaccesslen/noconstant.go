package avoidaccesslen

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func hasNoConstant(pass *analysis.Pass, arrObj types.Object, index ast.Expr) *[]analysis.SuggestedFix {
	indexCall, _ := index.(*ast.CallExpr)
	if indexCall == nil {
		return nil
	}
	indexFuncObj := extractIndexFuncObj(pass, indexCall)
	if lenObj != indexFuncObj {
		return nil
	}

	argObj := extractIndexFuncArgObj(pass, indexCall)
	if arrObj != argObj {
		return nil
	}

	fixes := []analysis.SuggestedFix{
		{
			Message: "Add ` - 1` after `len(arr)`",
			TextEdits: []analysis.TextEdit{{
				Pos:     indexCall.End(),
				End:     indexCall.End(),
				NewText: []byte("-1"),
			}},
		},
	}
	return &fixes
}
