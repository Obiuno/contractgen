package generators

import (
	"fmt"
	"strings"

	"contractgen/internal/schema"
)

func GenerateMermaid(s *schema.Schema) string {
	var erd strings.Builder

	erd.WriteString("erDiagram\n")

	// Entities
	for _, name := range s.OrderedTables {
		table := s.Tables[name]

		erd.WriteString(fmt.Sprintf("\t%s {\n", name))

		pkCol := make(map[string]bool, len(table.PrimaryKey))
		for _, pk := range table.PrimaryKey {
			pkCol[pk] = true
		}

		for _, colName := range table.OrderedColumns {
			col := table.Columns[colName]
			isPK := pkCol[colName]
			isFK := col.References != nil

			var marker string
			switch {
			case isPK && isFK:
				marker = "PK,FK"
			case isPK:
				marker = "PK"
			case isFK:
				marker = "FK"
			}

			if marker == "" {
				erd.WriteString(fmt.Sprintf("\t\t%s %s\n", col.Type, colName))
			} else {
				erd.WriteString(fmt.Sprintf("\t\t%s %s %s\n", col.Type, colName, marker))
			}
		}
		erd.WriteString("\t}\n")
	}

	// Relationships
	for _, name := range s.OrderedTables {
		table := s.Tables[name]
		for _, colName := range table.OrderedColumns {
			col := table.Columns[colName]
			if col.References == nil {
				continue
			}
			refTable := strings.ToLower(col.References.Table)
			erd.WriteString(fmt.Sprintf(
				"\t%s ||--o{ %s : \"via %s.%s → %s.%s\"\n",
				refTable, name,
				name, colName,
				refTable, strings.ToLower(col.References.Column),
			))
		}
	}

	return erd.String()
}
