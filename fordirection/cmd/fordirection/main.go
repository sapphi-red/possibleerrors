package main

import (
	"github.com/sapphi-red/possibleerrors/fordirection"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(fordirection.Analyzer) }
