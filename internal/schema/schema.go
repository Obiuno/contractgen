package schema

import (
	"contractgen/internal/parser"
	"strings"
)

type Schema struct {
	Tables        map[string]*Table
	OrderedTables []string
	Version       string
}

type Table struct {
	Columns        map[string]*Column
	OrderedColumns []string
	PrimaryKey     []string
	Unique         [][]string
}

type Column struct {
	Type        string
	Constraints []string
	References  *parser.Reference
}

func BuildSchema(contract parser.Contract) *Schema {
	schema := &Schema{
		Tables:  make(map[string]*Table, len(contract.Tables)),
		OrderedTables: make([]string, 0, len(contract.Tables)),
		Version: contract.Version,
	}

	for _, table := range contract.Tables {
		tName := strings.ToLower(table.Name)

		meta := &Table{
			Columns:    make(map[string]*Column, len(table.Columns)),
			OrderedColumns: make([]string, 0, len(table.Columns)),
			PrimaryKey: lowerAll(table.PrimaryKey),
			Unique:     lowerNested(table.Unique),
		}

		for _, col := range table.Columns {
			cName := strings.ToLower(col.Name)

			meta.Columns[cName] = &Column{
				Type:        col.Type,
				Constraints: col.Constraints,
				References:  col.References,
			}
			meta.OrderedColumns = append(meta.OrderedColumns, cName)
		}

		schema.Tables[tName] = meta
		schema.OrderedTables =append(schema.OrderedTables, tName)
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
