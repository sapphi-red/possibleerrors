package main

import (
	"github.com/sapphi-red/possibleerrors/mutexscope"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(mutexscope.Analyzer) }

