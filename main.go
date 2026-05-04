package main

import (
	"contractgen/internal/generators"
	"contractgen/internal/parser"
	"fmt"
)

func main() {
	contract, err := parser.LoadContract("contracts/example.yml")
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	//fmt.Printf("%+v\n", contract)

	jsonContract, err := parser.ContractToJSON(contract)
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Println(string(jsonContract))

	ddl := generators.GenerateDDL(contract)

	fmt.Println(ddl)
}
