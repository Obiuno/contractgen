package main

import (
	"contractgen/internal/parser"
	"fmt"
)

func main()  {
	contract, err := parser.LoadContract("contracts/example.yml")
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	fmt.Printf("%+v\n", contract)
}