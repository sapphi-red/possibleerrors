package fordirection

import (
	"errors"
	"go/ast"
	"go/token"
	"log"

	"golang.org/x/tools/go/analysis"
)

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

func getDirectionFromIncDec(incDec *ast.IncDecStmt) uint8 {
	if incDec.Tok == token.INC {
		return upToDirection
	}
	return downToDirection
}
