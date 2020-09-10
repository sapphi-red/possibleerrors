package fordirection

import (
	"errors"
	"go/ast"
	"go/constant"
	"go/token"
	"log"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func extractCounterAndCreateSuggestionFromAssign(pass *analysis.Pass, assign *ast.AssignStmt) (*ast.Ident, uint8, analysis.SuggestedFix, error) {
	if isCompoundAssign(assign.Tok) {
		return extractCounterAndCreateSuggestionFromCompoundAssign(pass, assign)
	}
	if isSimpleAssign(assign.Tok) {
		return extractCounterAndCreateSuggestionFromSimpleAssign(pass, assign)
	}
	return nil, 0, analysis.SuggestedFix{}, errors.New("Not an expected assignment.")
}

func extractCounterAndCreateSuggestionFromCompoundAssign(pass *analysis.Pass, assign *ast.AssignStmt) (*ast.Ident, uint8, analysis.SuggestedFix, error) {
	assignDirection := getDirectionAssign(assign)
	if assignDirection == noneDirection {
		return nil, 0, analysis.SuggestedFix{}, errors.New("Not a simple assignment.")
	}

	if len(assign.Lhs) > 1 || len(assign.Rhs) > 1 {
		return nil, 0, analysis.SuggestedFix{}, errors.New("Not a simple assignment.")
	}

	rhsType := pass.TypesInfo.Types[assign.Rhs[0]]
	invertDirection := false
	// 値がわかっているとき
	if rhsType.Value != nil {
		if rhsType.Value.Kind() != constant.Int && rhsType.Value.Kind() != constant.Float {
			return nil, 0, analysis.SuggestedFix{}, errors.New("Assigned value is not constant.")
		}
		// 値が負のときは方向を反転
		if strings.HasPrefix(rhsType.Value.ExactString(), "-") {
			invertDirection = true
		}
	}

	counter, _ := assign.Lhs[0].(*ast.Ident)
	if counter == nil {
		// ここに入ることなさそうだけど一応
		return nil, 0, analysis.SuggestedFix{}, errors.New("Not a simple assignment.")
	}

	if invertDirection {
		if assignDirection == upToDirection {
			assignDirection = downToDirection
		} else if assignDirection == downToDirection {
			assignDirection = upToDirection
		}
	}

	assignFix := analysis.SuggestedFix{
		Message: "Reverse assign (+= to -=, -= to +=)",
		TextEdits: []analysis.TextEdit{{
			Pos:     assign.TokPos,
			End:     assign.Rhs[0].Pos(),
			NewText: []byte(getReversedAssignTokenString(assign.Tok)),
		}},
	}

	return counter, assignDirection, assignFix, nil
}

func extractCounterAndCreateSuggestionFromSimpleAssign(pass *analysis.Pass, assign *ast.AssignStmt) (*ast.Ident, uint8, analysis.SuggestedFix, error) {
	if len(assign.Lhs) > 1 || len(assign.Rhs) > 1 {
		return nil, 0, analysis.SuggestedFix{}, errors.New("Not a simple assignment.")
	}

	rightBinaryExpr, _ := assign.Rhs[0].(*ast.BinaryExpr)
	if rightBinaryExpr == nil {
		return nil, 0, analysis.SuggestedFix{}, errors.New("Not a simple assignment.")
	}

	counter, _ := assign.Lhs[0].(*ast.Ident)
	if counter == nil {
		// ここに入ることなさそうだけど一応
		return nil, 0, analysis.SuggestedFix{}, errors.New("Not a simple assignment.")
	}

	rightXIdent, _ := rightBinaryExpr.X.(*ast.Ident)
	rightYIdent, _ := rightBinaryExpr.Y.(*ast.Ident)
	if rightXIdent == nil && rightYIdent == nil {
		return nil, 0, analysis.SuggestedFix{}, errors.New("Not a simple assignment.")
	}

	var diffNumExpr ast.Expr = nil
	if rightXIdent != nil && rightXIdent.Name == counter.Name {
		diffNumExpr = rightYIdent
	}
	if rightYIdent != nil && rightYIdent.Name == counter.Name {
		diffNumExpr = rightXIdent
	}
	if diffNumExpr == nil {
		return nil, 0, analysis.SuggestedFix{}, errors.New("Not a simple assignment.")
	}

	rhsType := pass.TypesInfo.Types[diffNumExpr]
	invertDirection := false
	// 値がわかっているとき
	if rhsType.Value != nil {
		if rhsType.Value.Kind() != constant.Int && rhsType.Value.Kind() != constant.Float {
			return nil, 0, analysis.SuggestedFix{}, errors.New("Assigned value is not constant.")
		}
		// 値が負のときは方向を反転
		if strings.HasPrefix(rhsType.Value.ExactString(), "-") {
			invertDirection = true
		}
	}

	assignDirection := getDirectionBinary(rightBinaryExpr)

	if invertDirection {
		if assignDirection == upToDirection {
			assignDirection = downToDirection
		} else if assignDirection == downToDirection {
			assignDirection = upToDirection
		}
	}

	assignFix := analysis.SuggestedFix{
		Message: "Reverse assign (+ to -, - to +)",
		TextEdits: []analysis.TextEdit{{
			Pos:     rightBinaryExpr.OpPos,
			End:     rightBinaryExpr.Y.Pos(),
			NewText: []byte(getReversedBinaryTokenString(rightBinaryExpr.Op)),
		}},
	}

	return counter, assignDirection, assignFix, nil
}

func isCompoundAssign(tok token.Token) bool {
	return tok == token.ADD_ASSIGN || tok == token.SUB_ASSIGN
}

func isSimpleAssign(tok token.Token) bool {
	return tok == token.ASSIGN
}

func getDirectionAssign(assign *ast.AssignStmt) uint8 {
	if assign.Tok == token.ADD_ASSIGN {
		return upToDirection
	}
	if assign.Tok == token.SUB_ASSIGN {
		return downToDirection
	}
	return noneDirection
}

func getDirectionBinary(binary *ast.BinaryExpr) uint8 {
	if binary.Op == token.ADD {
		return upToDirection
	}
	if binary.Op == token.SUB {
		return downToDirection
	}
	return noneDirection
}

func getReversedAssignTokenString(t token.Token) string {
	switch t {
	case token.ADD_ASSIGN:
		return token.SUB_ASSIGN.String()
	case token.SUB_ASSIGN:
		return token.ADD_ASSIGN.String()
	}

	log.Fatalf("Unexpected token passed to getReversedAssignTokenString: %#v", t)
	return ""
}

func getReversedBinaryTokenString(t token.Token) string {
	switch t {
	case token.ADD:
		return token.SUB.String()
	case token.SUB:
		return token.ADD.String()
	}

	log.Fatalf("Unexpected token passed to getReversedBinaryTokenString: %#v", t)
	return ""
}
