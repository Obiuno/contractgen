package main

import (
	"contractgen/internal/cli"
	"contractgen/internal/generators"
	"contractgen/internal/parser"
	"contractgen/internal/schema"
	"contractgen/internal/validator"
	"flag"
	"fmt"
	"os"
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

	if errs := validator.Validate(contract); len(errs) >0 {
		fmt.Fprintln(os.Stderr, "validation failed:")
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, " - %s\n", e.Error())
		}
		os.Exit(1)
	}

	if cfg.ShowJSON {
		jsonContract, err := parser.ContractToJSON(contract)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error converting to JSON:", err)
			os.Exit(1)
		}

		fmt.Println(string(jsonContract))

	}

	//build the schema before generating contract, as refactor for using validated schemas

	sch := schema.BuildSchema(contract)
	fmt.Printf("%+v\n", sch)

	ddl := generators.GenerateDDL(contract)

	//convert the string to byte slice
	data := []byte(ddl)

	// output is where to save it
	err = os.WriteFile(cfg.Output, data, 0644)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save SQL to %s: %v\n", cfg.Output, err)
		os.Exit(1)
	}

	fmt.Printf("Success! SQL generated and saved to: %s\n", cfg.Output)
}
