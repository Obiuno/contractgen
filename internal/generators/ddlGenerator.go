package generators

import (
	"contractgen/internal/parser"
	"slices"
	"strings"
)

// create table and name

// loop through columns
// columns data type, a struct

// column type

func GenerateDDL(contract parser.Contract) string {

	var allTables []string

	for _, table := range contract.Tables {

		var ddl strings.Builder
		ddl.WriteString("CREATE TABLE IF NOT EXISTS ")
		ddl.WriteString(strings.ToLower(table.Name))
		ddl.WriteString(" (\n")

		var columns []string
		for _, column := range table.Columns {
			var col strings.Builder
			col.WriteString("\t")
			col.WriteString(strings.ToLower(column.Name))
			col.WriteString(" ")
			col.WriteString(strings.ToUpper(column.Type))

			if slices.Contains(column.Constraints, "primary_key") {
				col.WriteString(" PRIMARY KEY")
			}

			if slices.Contains(column.Constraints, "unique") {
				col.WriteString(" UNIQUE")
			}

			if slices.Contains(column.Constraints, "not_null") {
				col.WriteString(" NOT NULL")
			}

			columns = append(columns, col.String())
		}

		ddl.WriteString(strings.Join(columns, ",\n"))
		ddl.WriteString("\n);")
		allTables = append(allTables, ddl.String())
	}

	return strings.Join(allTables, "\n\n")

}
