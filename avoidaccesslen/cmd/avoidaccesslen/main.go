package main

import (
	"github.com/sapphi-red/possibleerrors/avoidaccesslen"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(avoidaccesslen.Analyzer) }

