package parser

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Column struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Constraints []string `yaml:"constraints"`
}

type Table struct {
	Name string `yaml:"name"`
	Columns []Column `yaml:"columns"`
	Description string `yaml:"description"`
}

type Contract struct {
	Version string `yaml:"version"`
	Tables []Table `yaml:"tables"`
	Description string `yaml:"description"`
}

func LoadContract(path string) (Contract, error) {
	var contract Contract

	data, err := os.ReadFile(path)

	// classic error handling
	if err != nil {
		return contract, err
	}

	err = yaml.Unmarshal(data, &contract)
    if err != nil {
		fmt.Println("Error:", err)
        return contract, err
    }

	return contract, nil
}