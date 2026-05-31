package main

import (
	"contractgen/internal/cli"
	"contractgen/internal/generators"
	"contractgen/internal/normalise"
	"contractgen/internal/parser"
	"contractgen/internal/validator"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// get instructions from the CLI
	cfg := cli.SetupFlags()

	if cfg.Input == "" {
		fmt.Fprintln(os.Stderr, "Error: --input is required")
		flag.Usage()
		os.Exit(1)
	}

	// send it correct place
	contract, err := parser.LoadContract(cfg.Input)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading contract:", err)
		os.Exit(1)
	}

	normalise.Normalise(&contract)

	issues := validator.Validate(contract)
	errs, warnings := validator.Partition(issues)

	for _, w := range warnings {
		fmt.Fprintln(os.Stderr, "WARNING:", w.Error())
	}

	if len(errs) > 0 {
		fmt.Fprintln(os.Stderr, "Validation failed:")

		for _, e := range errs {
			fmt.Fprintf(os.Stderr, " - %s\n", e.Error())
		}
		os.Exit(1)
	}

	if cfg.JSON {
		jsonContract, err := parser.ContractToJSON(contract)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error converting to JSON:", err)
			os.Exit(1)
		}
		path := filepath.Join(cfg.OutputDir, "schema.json")
		if err := os.WriteFile(path, jsonContract, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to save JSON to %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf("JSON written to %s\n", path)
	}

	if cfg.DDL {
		ddl := generators.GenerateDDL(contract)
		path := filepath.Join(cfg.OutputDir, "schema.sql")
		if err := os.WriteFile(path, []byte(ddl), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to save SQL to %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf("DDL written to %s\n", path)
	}

	if cfg.Doc {
		doc := generators.GenerateDocs(contract)
		path := filepath.Join(cfg.OutputDir, "schema.md")
		if err := os.WriteFile(path, []byte(doc), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to save markdown to %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf("Docs written to %s\n", path)
	}

	if cfg.Mermaid {
		mmd := generators.GenerateMermaid(contract)
		path := filepath.Join(cfg.OutputDir, "schema.mmd")
		if err := os.WriteFile(path, []byte(mmd), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to save mermaid to %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf("Mermaid written to %s\n", path)
	}

}
