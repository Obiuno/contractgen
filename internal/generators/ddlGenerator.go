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

		// lower the PK list once for lookup
		pkColumns := make([]string, len(table.PrimaryKey))
		for i, col := range table.PrimaryKey {
			pkColumns[i] = strings.ToLower(col)
		}

		tables.WriteString(fmt.Sprintf("-- Table: %s\n", tableName))
		tables.WriteString("CREATE TABLE IF NOT EXISTS ")
		tables.WriteString(tableName)
		tables.WriteString(" (\n")

		columns := make([]string, 0, len(table.Columns))
		var tableConstraints strings.Builder

		// Pass 1: column definitions and column-level constraints
		for _, column := range table.Columns {
			colName := strings.ToLower(column.Name)
			isPK := slices.Contains(pkColumns, colName)

			columns = append(columns, fmt.Sprintf(
				"\t%s %s",
				colName,
				strings.ToUpper(column.Type),
			))

			// NOT NULL: explicit OR inferred from being a PK column
			if slices.Contains(column.Constraints, "not_null") || isPK {
				tableConstraints.WriteString(fmt.Sprintf(
					"ALTER TABLE %s ALTER COLUMN %s SET NOT NULL;\n",
					tableName, colName,
				))
			}

			// FK still column-level (single-column FK shorthand)
			if column.References != nil {
				tableConstraints.WriteString(fmt.Sprintf(
					"ALTER TABLE %s ADD CONSTRAINT fk_%s_%s FOREIGN KEY (%s) REFERENCES %s (%s);\n",
					tableName, tableName, colName, colName,
					strings.ToLower(column.References.Table),
					strings.ToLower(column.References.Column),
				))
			}
		}

		// Pass 2: table-level PRIMARY KEY
		if len(pkColumns) > 0 {
			tableConstraints.WriteString(fmt.Sprintf(
				"ALTER TABLE %s ADD CONSTRAINT pk_%s PRIMARY KEY (%s);\n",
				tableName, tableName, strings.Join(pkColumns, ", "),
			))
		}

		// Pass 3: table-level UNIQUE constraints
		for _, group := range table.Unique {
			// skip if this group, is the same as the PK (already enforced)
			if slices.Equal(pkColumns, lowerAll(group)) {
				continue
			}

			lowered := lowerAll(group)
			tableConstraints.WriteString(fmt.Sprintf(
				"ALTER TABLE %s ADD CONSTRAINT uq_%s_%s UNIQUE (%s);\n",
				tableName, tableName, strings.Join(lowered, "_"),
				strings.Join(lowered, ", "),
			))
		}

		tables.WriteString(strings.Join(columns, ",\n"))
		tables.WriteString("\n);\n\n")

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

func lowerAll(names []string) []string {
	out := make([]string, len(names))
	for i, n := range names {
		out[i] = strings.ToLower(n)
	}
	return out
}
