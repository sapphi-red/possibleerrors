package avoidaccesslen_test

import (
	"testing"

	"github.com/sapphi-red/possibleerrors/avoidaccesslen"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, avoidaccesslen.Analyzer, "a")
}
