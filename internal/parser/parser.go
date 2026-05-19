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
	Name        string     `yaml:"name" json:"name"`
	Type        string     `yaml:"type" json:"type"`
	Description string     `yaml:"description,omitempty" json:"description,omitempty"`
	Constraints []string   `yaml:"constraints,omitempty" json:"constraints,omitempty"`
	References  *Reference `yaml:"references,omitempty" json:"references,omitempty"`
}

type Table struct {
	Name        string     `yaml:"name" json:"name"`
	Description string     `yaml:"description,omitempty" json:"description,omitempty"`
	Columns     []Column   `yaml:"columns" json:"columns"`
	PrimaryKey  []string   `yaml:"primary_key,omitempty" json:"primary_key,omitempty"`
	Unique      [][]string `yaml:"unique,omitempty" json:"unique,omitempty"`
}

type Contract struct {
	Version     string  `yaml:"version" json:"version"`
	Description string  `yaml:"description,omitempty" json:"description,omitempty"`
	Tables      []Table `yaml:"tables" json:"tables"`
}

func ParseContract(data []byte) (Contract, error) {
	var contract Contract
	if err := yaml.Unmarshal(data, &contract); err != nil {
		return contract, fmt.Errorf("parsing contract yaml: %w", err)
	}
	return contract, nil
}

func LoadContract(path string) (Contract, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return Contract{}, fmt.Errorf("reading contract file %q: %w", path, err)
	}

	return ParseContract(data)
}

func ContractToJSON(contract Contract) ([]byte, error) {
	return json.MarshalIndent(contract, "", "  ")
}
