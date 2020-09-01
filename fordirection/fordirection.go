package fordirection

import (
	"errors"
	"go/ast"
	"go/token"
	"log"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	doc = "fordirection finds for-loops which likely has a wrong direction"
)

const (
	upToDirection = iota
	downToDirection
	noneDirection
	invalidDirection
)

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "fordirection",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		new(ast.ForStmt),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		forLoop := n.(*ast.ForStmt)

		// 条件が大小比較の二項演算のみ対象
		condBianry, _ := forLoop.Cond.(*ast.BinaryExpr)
		if condBianry == nil || !isComparationToken(condBianry.Op) {
			return
		}

		counter, assignDirection, postFix, err := extractCounterAndCreateSuggestion(forLoop.Post)
		if err != nil {
			return
		}

		condDirection := getLoopDirection(condBianry, counter)
		if condDirection == noneDirection {
			return
		}
		if condDirection == invalidDirection {
			pass.Reportf(n.Pos(), "Loop direction seems to be wrong.")
			return
		}

		if condDirection != assignDirection {
			// TODO: Auto detect
			conditionFix := analysis.SuggestedFix{
				Message: "Reverse condition (> to <, < to >)",
				TextEdits: []analysis.TextEdit{{
					Pos:     condBianry.OpPos,
					End:     condBianry.Y.Pos(),
					NewText: []byte(getReversedComparationTokenString(condBianry.Op)),
				}},
			}

			pass.Report(analysis.Diagnostic{
				Pos:            forLoop.Pos(),
				End:            forLoop.Post.End(),
				Message:        "Loop direction seems to be wrong.",
				SuggestedFixes: []analysis.SuggestedFix{conditionFix, postFix},
			})
			return
		}
	})

	return nil, nil
}

func isComparationToken(t token.Token) bool {
	return (t == token.LSS) || (t == token.GTR) || (t == token.LEQ) || (t == token.GEQ)
}

func getLoopDirection(cond *ast.BinaryExpr, counter *ast.Ident) uint8 {
	xIdent, _ := cond.X.(*ast.Ident)
	yIdent, _ := cond.Y.(*ast.Ident)

	if xIdent == nil && yIdent == nil {
		return noneDirection
	}
	if (xIdent != nil && yIdent != nil) && (xIdent.Name == counter.Name && yIdent.Name == counter.Name) {
		return invalidDirection
	}

	if xIdent != nil && xIdent.Name == counter.Name {
		if cond.Op == token.LSS || cond.Op == token.LEQ {
			return upToDirection
		}
		return downToDirection
	}
	if yIdent != nil && yIdent.Name == counter.Name {
		if cond.Op == token.LSS || cond.Op == token.LEQ {
			return downToDirection
		}
		return upToDirection
	}
	return noneDirection
}

func getDirectionFromIncDec(incDec *ast.IncDecStmt) uint8 {
	if incDec.Tok == token.INC {
		return upToDirection
	}
	return downToDirection
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

func getReversedComparationTokenString(t token.Token) string {
	switch t {
	case token.LSS:
		return ">"
	case token.GTR:
		return "<"
	case token.LEQ:
		return ">="
	case token.GEQ:
		return "<="
	}

	log.Fatalf("Unexpected token passed to getReversedComparationTokenString: %#v", t)
	return ""
}

func getReversedIncDecTokenString(t token.Token) string {
	switch t {
	case token.INC:
		return "--"
	case token.DEC:
		return "++"
	}

	log.Fatalf("Unexpected token passed to getReversedIncDecTokenString: %#v", t)
	return ""
}

func getReversedAssignTokenString(t token.Token) string {
	switch t {
	case token.ADD_ASSIGN:
		return "-="
	case token.SUB_ASSIGN:
		return "+="
	}

	log.Fatalf("Unexpected token passed to getReversedAssignTokenString: %#v", t)
	return ""
}

func extractCounterAndCreateSuggestion(post ast.Stmt) (*ast.Ident, uint8, analysis.SuggestedFix, error) {
	switch post := post.(type) {
	case *ast.IncDecStmt:
		return extractCounterAndCreateSuggestionFromIncDec(post)
	case *ast.AssignStmt:
		return extractCounterAndCreateSuggestionFromAssign(post)
	}
	// TODO: i = i + 5
	return nil, 0, analysis.SuggestedFix{}, errors.New("Not increment/descriment.")
}

func extractCounterAndCreateSuggestionFromIncDec(incDec *ast.IncDecStmt) (*ast.Ident, uint8, analysis.SuggestedFix, error) {
	counter, _ := incDec.X.(*ast.Ident)
	if counter == nil {
		// ここに入ることなさそうだけど一応
		return nil, 0, analysis.SuggestedFix{}, errors.New("Missing identifier.")
	}

	incDecDirection := getDirectionFromIncDec(incDec)

	incDecFix := analysis.SuggestedFix{
		Message: "Reverse increment/decrement (++ to --, -- to ++)",
		TextEdits: []analysis.TextEdit{{
			Pos:     incDec.TokPos,
			End:     incDec.End(),
			NewText: []byte(getReversedIncDecTokenString(incDec.Tok)),
		}},
	}

	return counter, incDecDirection, incDecFix, nil
}

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
