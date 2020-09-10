package avoidaccesslen

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strconv"

	"golang.org/x/tools/go/analysis"
)

func hasConstantOnLeft(pass *analysis.Pass, arrObj types.Object, binaryExpr *ast.BinaryExpr) *[]analysis.SuggestedFix {
	if !(binaryExpr.Op == token.ADD) {
		return nil
	}

	x := binaryExpr.X
	y := binaryExpr.Y

	yCall, _ := y.(*ast.CallExpr)
	if yCall == nil {
		return nil
	}
	yFuncObj := extractIndexFuncObj(pass, yCall)
	if lenObj != yFuncObj {
		return nil
	}

	argObj := extractIndexFuncArgObj(pass, yCall)
	if arrObj != argObj {
		return nil
	}

	xVal := pass.TypesInfo.Types[x].Value
	if xVal == nil || xVal.Kind() != constant.Int {
		return nil
	}
	i, err := strconv.Atoi(xVal.String())
	if err != nil {
		return nil
	}

	if i < 0 {
		return nil
	}

	fixes := []analysis.SuggestedFix{}
	return &fixes
}
