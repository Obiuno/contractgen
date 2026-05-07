package main

import (
	"contractgen/internal/cli"
	"contractgen/internal/generators"
	"contractgen/internal/parser"
	"flag"
	"fmt"
	"os"
)

func main() {
	// get instructions from the CLI
	config := cli.SetupFlags()

	if config.Input == "" {
		fmt.Fprintln(os.Stderr, "Error: --input is required")
		flag.Usage()
		os.Exit(1)
	}

	// send it correct place
	contract, err := parser.LoadContract(config.Input)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading contract:", err)
		os.Exit(1)
	}

	if config.ShowJSON {
		jsonContract, err := parser.ContractToJSON(contract)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error converting to JSON:", err)
			os.Exit(1)
		}

		fmt.Println(string(jsonContract))

	}

	ddl := generators.GenerateDDL(contract)

	//convert the string to byte slice
	data := []byte(ddl)

	// output is where to save it
	err = os.WriteFile(config.Output, data, 0644)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to save SQL to %s: %v\n", config.Output, err)
		os.Exit(1)
	}

	fmt.Printf("Success! SQL generated and saved to: %s\n", config.Output)
}
