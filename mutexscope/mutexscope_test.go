package mutexscope_test

import (
	"testing"

	"github.com/sapphi-red/possibleerrors/mutexscope"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, mutexscope.Analyzer, "a")
	analysistest.Run(t, testdata, mutexscope.Analyzer, "b")
}
