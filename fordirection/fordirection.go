package fordirection

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	doc = "fordirection is ..."
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

		// インクリメント/デクリメントのみ対象
		incDec, _ := forLoop.Post.(*ast.IncDecStmt)
		if incDec == nil {
			fmt.Print("b")
			return
		}
		counter, _ := incDec.X.(*ast.Ident)
		if counter == nil {
			// ここに入ることなさそうだけど一応
			fmt.Print("c")
			return
		}

		condDirection := getLoopDirection(condBianry, counter)
		if condDirection == noneDirection {
			return
		}
		if condDirection == invalidDirection {
			pass.Reportf(n.Pos(), "loop direction seems to be wrong.")
			return
		}

		incDecDirection := getDirectionFromIncDec(incDec)

		if condDirection != incDecDirection {
			pass.Reportf(n.Pos(), "loop direction seems to be wrong.")
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
