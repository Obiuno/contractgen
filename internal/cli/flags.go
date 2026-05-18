package cli

import (
	"flag"
)

type Config struct {
	Input     string
	OutputDir string
	DDL       bool
	Doc       bool
	JSON      bool
	Mermaid   bool
}

func SetupFlags() Config {
	input := flag.String("input", "", "YAML contract path")
	outputDir := flag.String("output-dir", ".", "Directory to write output files")
	ddl := flag.Bool("ddl", false, "produce SQL DDL")
	doc := flag.Bool("doc", false, "produce markdown docs")
	json := flag.Bool("json", false, "produce JSON")
	mermaid := flag.Bool("mermaid", false, "produce mermiad ERD")
	flag.Parse()

	cfg := Config{
		Input:     *input,
		OutputDir: *outputDir,
		DDL:       *ddl,
		Doc:       *doc,
		JSON:      *json,
		Mermaid:   *mermaid,
	}

	// If no output flag is set, produce everything
	if !cfg.DDL && !cfg.Doc && !cfg.JSON && !cfg.Mermaid {
		cfg.DDL = true
		cfg.Doc = true
		cfg.JSON = true
		cfg.Mermaid = true
	}

	return cfg
}
