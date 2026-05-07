package generators

import (
	"fmt"
	"slices"
	"strings"

	"contractgen/internal/parser"
)

func GenerateDDL(contract parser.Contract) string {
	var tables strings.Builder
	var constraints strings.Builder

	tables.WriteString("-- ==========================================\n")
	tables.WriteString("-- TABLES\n")
	tables.WriteString("-- ==========================================\n\n")

	constraints.WriteString("-- ==========================================\n")
	constraints.WriteString("-- CONSTRAINTS (PK, FK, UNIQUE, NOT NULL)\n")
	constraints.WriteString("-- ==========================================\n\n")

	for _, table := range contract.Tables {
		tableName := strings.ToLower(table.Name)

		tables.WriteString(fmt.Sprintf("-- Table: %s\n", tableName))
		tables.WriteString("CREATE TABLE IF NOT EXISTS ")
		tables.WriteString(tableName)
		tables.WriteString(" (\n")

		columns := make([]string, 0, len(table.Columns))
		var tableConstraints strings.Builder

		for _, column := range table.Columns {
			colName := strings.ToLower(column.Name)

			columns = append(columns, fmt.Sprintf(
				"\t%s %s",
				colName,
				strings.ToUpper(column.Type),
			))

			if slices.Contains(column.Constraints, "not_null") {
				tableConstraints.WriteString(fmt.Sprintf(
					"ALTER TABLE %s ALTER COLUMN %s SET NOT NULL;\n",
					tableName, colName,
				))
			}
			if slices.Contains(column.Constraints, "primary_key") {
				tableConstraints.WriteString(fmt.Sprintf(
					"ALTER TABLE %s ADD CONSTRAINT pk_%s PRIMARY KEY (%s);\n",
					tableName, tableName, colName,
				))
			}
			if slices.Contains(column.Constraints, "unique") {
				tableConstraints.WriteString(fmt.Sprintf(
					"ALTER TABLE %s ADD CONSTRAINT uq_%s_%s UNIQUE (%s);\n",
					tableName, tableName, colName, colName,
				))
			}
			if column.References != nil {
				tableConstraints.WriteString(fmt.Sprintf(
					"ALTER TABLE %s ADD CONSTRAINT fk_%s_%s FOREIGN KEY (%s) REFERENCES %s (%s);\n",
					tableName, tableName, colName, colName,
					strings.ToLower(column.References.Table),
					strings.ToLower(column.References.Column),
				))
			}
		}

		tables.WriteString(strings.Join(columns, ",\n"))
		tables.WriteString("\n);\n\n")

		// only write constraints section if there are any
		if tableConstraints.Len() > 0 {
			constraints.WriteString(fmt.Sprintf("-- Constraints for %s\n", tableName))
			constraints.WriteString(tableConstraints.String())
			constraints.WriteString("\n")
		}
	}

	var output strings.Builder
	output.WriteString("BEGIN;\n\n")
	output.WriteString(tables.String())
	output.WriteString(constraints.String())
	output.WriteString("COMMIT;\n")
	return output.String()
}
