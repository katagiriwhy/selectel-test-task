package main

import (
	"flag"
	"os"

	"golang.org/x/tools/go/analysis/singlechecker"

	"selectel/analyzer"
	"selectel/config"
)

func main() {
	var configPath string

	flag.StringVar(
		&configPath,
		"config",
		"",
		"path to config file (default: .sloglint.yaml)",
	)
	flag.Parse()

	config.Load(configPath)

	singlechecker.Main(analyzer.Analyzer)

	_ = os.Stdout
}
