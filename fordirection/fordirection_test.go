package fordirection_test

import (
	"testing"

	"github.com/sapphi-red/possibleerrors/fordirection"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, fordirection.Analyzer, "a")
}
