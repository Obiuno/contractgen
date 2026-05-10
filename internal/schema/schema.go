package schema

import (
	"contractgen/internal/parser"
	"slices"
	"strings"
)

type Schema struct {
	Tables     map[string]*Table
	Version    string
	Duplicates []DuplicateError
}

type Table struct {
	Columns    map[string]*Column
	PKs        []string
	UNIQUEs    []string
	Duplicates []DuplicateError
}

type Column struct {
	Type        string
	Constraints []string
	References  *parser.Reference
	Duplicates  []DuplicateError
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

		// detect duplicate table
		if _, exists := schema.Tables[tName]; exists {
			schema.Duplicates = append(schema.Duplicates, DuplicateError{
				Kind: "table",
				Name: tName,
			})
			continue // note and continue
		}

		meta := &Table{
			Columns: make(map[string]*Column, len(table.Columns)),
		}

		for _, col := range table.Columns {
			cName := strings.ToLower(col.Name)

			// handle column duplicates
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

			if slices.Contains(col.Constraints, "primary_key") {
				meta.PKs = append(meta.PKs, cName)
			}
			if slices.Contains(col.Constraints, "unique") {
				meta.UNIQUEs = append(meta.UNIQUEs, cName)
			}

		}
		schema.Tables[tName] = meta
	}

	return schema
}
