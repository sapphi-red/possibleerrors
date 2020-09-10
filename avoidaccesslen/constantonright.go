package avoidaccesslen

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/analysis"
)

func hasConstantOnRight(pass *analysis.Pass, arrObj types.Object, binaryExpr *ast.BinaryExpr) *[]analysis.SuggestedFix {
	if !(binaryExpr.Op == token.ADD || binaryExpr.Op == token.SUB) {
		return nil
	}

	x := binaryExpr.X
	y := binaryExpr.Y

	xCall, _ := x.(*ast.CallExpr)
	if xCall == nil {
		return nil
	}
	xFuncObj := extractIndexFuncObj(pass, xCall)
	if lenObj != xFuncObj {
		return nil
	}

	argObj := extractIndexFuncArgObj(pass, xCall)
	if arrObj != argObj {
		return nil
	}

	yVal := pass.TypesInfo.Types[y].Value
	if yVal == nil || yVal.Kind() != constant.Int {
		return nil
	}
	i, err := strconv.Atoi(yVal.String())
	if err != nil {
		return nil
	}

	if binaryExpr.Op == token.SUB {
		i *= -1
	}

	if i < 0 {
		return nil
	}

	fixes := []analysis.SuggestedFix{}
	if i == 1 {
		fixes = []analysis.SuggestedFix{
			{
				Message: "Change ` + 1` to ` - 1`",
				TextEdits: []analysis.TextEdit{{
					Pos:     xCall.End(),
					End:     binaryExpr.End(),
					NewText: []byte("-1"),
				}},
			},
		}
	}

	return &fixes
}
