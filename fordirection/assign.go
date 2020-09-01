package fordirection

import (
	"errors"
	"go/ast"
	"go/token"
	"log"

	"golang.org/x/tools/go/analysis"
)

func extractCounterAndCreateSuggestionFromAssign(assign *ast.AssignStmt) (*ast.Ident, uint8, analysis.SuggestedFix, error) {
	if len(assign.Lhs) > 1 {
		return nil, 0, analysis.SuggestedFix{}, errors.New("Not a simple assignment.")
	}
	counter, _ := assign.Lhs[0].(*ast.Ident)
	if counter == nil {
		// ここに入ることなさそうだけど一応
		return nil, 0, analysis.SuggestedFix{}, errors.New("Not a simple assignment.")
	}

	assignDirection := getDirectionAssign(assign)
	if assignDirection == noneDirection {
		return nil, 0, analysis.SuggestedFix{}, errors.New("Not a simple assignment.")
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

func getDirectionAssign(assign *ast.AssignStmt) uint8 {
	if assign.Tok == token.ADD_ASSIGN {
		return upToDirection
	}
	if assign.Tok == token.SUB_ASSIGN {
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
