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

	// Section Headers
	tables.WriteString("-- ==========================================\n")
	tables.WriteString("-- TABLES\n")
	tables.WriteString("-- ==========================================\n\n")

	constraints.WriteString("-- ==========================================\n")
	constraints.WriteString("-- CONSTRAINTS (PK, FK, UNIQUE, NOT NULL)\n")
	constraints.WriteString("-- ==========================================\n\n")

	for _, table := range contract.Tables {
		tableName := strings.ToLower(table.Name)

		// CREATE TABLE — pure structure, no constraints
		tables.WriteString(fmt.Sprintf("-- Table: %s\n", tableName))
		tables.WriteString("CREATE TABLE IF NOT EXISTS ")
		tables.WriteString(tableName)
		tables.WriteString(" (\n")

		columns := make([]string, 0, len(table.Columns))

		constraints.WriteString(fmt.Sprintf("-- Constraints for %s\n", tableName))

		for _, column := range table.Columns {
			colName := strings.ToLower(column.Name)

			// column, just name and type
			columns = append(columns, fmt.Sprintf(
				"\t%s %s",
				colName,
				strings.ToUpper(column.Type),
			))

			// NOT NULL first (must come before PRIMARY KEY)
			if slices.Contains(column.Constraints, "not_null") {
				constraints.WriteString(fmt.Sprintf(
					"ALTER TABLE %s ALTER COLUMN %s SET NOT NULL;\n",
					tableName, colName,
				))
			}

			// then PRIMARY KEY
			if slices.Contains(column.Constraints, "primary_key") {
				constraints.WriteString(fmt.Sprintf(
					"ALTER TABLE %s ADD CONSTRAINT pk_%s PRIMARY KEY (%s);\n",
					tableName, tableName, colName,
				))
			}

			// then UNIQUE
			if slices.Contains(column.Constraints, "unique") {
				constraints.WriteString(fmt.Sprintf(
					"ALTER TABLE %s ADD CONSTRAINT uq_%s_%s UNIQUE (%s);\n",
					tableName, tableName, colName, colName,
				))
			}

			// then FOREIGN KEY
			if column.References != nil {
				refTable := strings.ToLower(column.References.Table)
				refColumn := strings.ToLower(column.References.Column)

				constraints.WriteString(fmt.Sprintf(
					"ALTER TABLE %s ADD CONSTRAINT fk_%s_%s FOREIGN KEY (%s) REFERENCES %s (%s);\n",
					tableName, tableName, colName, colName, refTable, refColumn,
				))
			}

		}

		tables.WriteString(strings.Join(columns, ",\n"))
		tables.WriteString("\n);\n\n")
	}

	return "BEGIN;\n\n" + tables.String() + "\n" + constraints.String() + "\nCOMMIT;"
}
