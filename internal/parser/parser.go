package parser

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Reference struct {
	Table  string `yaml:"table" json:"table"`
	Column string `yaml:"column" json:"column"`
}

type Column struct {
	Name        string    `yaml:"name" json:"name"`
	Type        string    `yaml:"type" json:"type"`
	Constraints []string  `yaml:"constraints" json:"constraints"`
	References  *Reference `yaml:"references,omitempty" json:"references,omitempty"`
}

type Table struct {
	Name        string   `yaml:"name" json:"name"`
	Description string   `yaml:"description" json:"description"`
	Columns     []Column `yaml:"columns" json:"columns"`
}

type Contract struct {
	Version     string  `yaml:"version" json:"version"`
	Description string  `yaml:"description" json:"description"`
	Tables      []Table `yaml:"tables" json:"tables"`
}

func LoadContract(path string) (Contract, error) {
	var contract Contract

	data, err := os.ReadFile(path)

	// classic error handling
	if err != nil {
		fmt.Println("Error reading file", err)
		return contract, err
	}

	err = yaml.Unmarshal(data, &contract)
	if err != nil {
		fmt.Println("Error parsing YAML:", err)
		return contract, err
	}

	return contract, nil
}

func ContractToJSON(contract Contract) ([]byte, error) {
	return json.MarshalIndent(contract, "", "  ")
}
