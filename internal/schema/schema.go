package schema

import (
	"contractgen/internal/parser"
	"strings"
)

type Schema struct {
	Tables     map[string]*Table
	Version    string
	Duplicates []DuplicateError
}

type Table struct {
	Columns    map[string]*Column
	PrimaryKey []string
	Unique     [][]string
	Duplicates []DuplicateError
}

type Column struct {
	Type        string
	Constraints []string
	References  *parser.Reference
}

type DuplicateError struct {
	Kind    string
	Name    string
	Context string
}

func BuildSchema(contract parser.Contract) *Schema {
	schema := &Schema{
		Tables:  make(map[string]*Table, len(contract.Tables)),
		Version: contract.Version,
	}

	for _, table := range contract.Tables {
		tName := strings.ToLower(table.Name)

		if _, exists := schema.Tables[tName]; exists {
			schema.Duplicates = append(schema.Duplicates, DuplicateError{
				Kind: "table",
				Name: tName,
			})
			continue
		}

		meta := &Table{
			Columns:    make(map[string]*Column, len(table.Columns)),
			PrimaryKey: lowerAll(table.PrimaryKey),
			Unique:     lowerNested(table.Unique),
		}

		for _, col := range table.Columns {
			cName := strings.ToLower(col.Name)

			if _, exists := meta.Columns[cName]; exists {
				meta.Duplicates = append(meta.Duplicates, DuplicateError{
					Kind:    "column",
					Name:    cName,
					Context: tName,
				})
				continue
			}

			meta.Columns[cName] = &Column{
				Type:        col.Type,
				Constraints: col.Constraints,
				References:  col.References,
			}
		}

		schema.Tables[tName] = meta
	}

	return schema
}

// helper to lowercase a slice of column names
func lowerAll(names []string) []string {
	out := make([]string, len(names))
	for i, n := range names {
		out[i] = strings.ToLower(n)
	}
	return out
}

// helper to lowercase nested slices (for Unique constraints)
func lowerNested(groups [][]string) [][]string {
	out := make([][]string, len(groups))
	for i, g := range groups {
		out[i] = lowerAll(g)
	}
	return out
}
