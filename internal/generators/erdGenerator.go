package generators

import (
	"fmt"
	"strings"

	"contractgen/internal/parser"
)

func GenerateMermaid(contract parser.Contract) string {
	var erd strings.Builder

	erd.WriteString("erDiagram\n")

	// Entities
	for _, table := range contract.Tables {
		tName := strings.ToLower(table.Name)

		erd.WriteString(fmt.Sprintf("\t%s {\n", tName))

		pkCol := make(map[string]bool, len(table.PrimaryKey))
		for _, pk := range table.PrimaryKey {
			pkCol[strings.ToLower(pk)] = true
		}

		for _, col := range table.Columns {
			cName := strings.ToLower(col.Name)
			isPK := pkCol[cName]
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
				erd.WriteString(fmt.Sprintf("\t\t%s %s\n", col.Type, cName))
			} else {
				erd.WriteString(fmt.Sprintf("\t\t%s %s %s\n", col.Type, cName, marker))
			}
		}
		erd.WriteString("\t}\n")
	}

	// Relationships
	for _, table := range contract.Tables {
		tName := strings.ToLower(table.Name)
		for _, col := range table.Columns {
			if col.References == nil {
				continue
			}

			cName := strings.ToLower(col.Name)

			refTable := strings.ToLower(col.References.Table)
			erd.WriteString(fmt.Sprintf(
				"\t%s ||--o{ %s : \"via %s.%s → %s.%s\"\n",
				refTable, tName,
				tName, cName,
				refTable, strings.ToLower(col.References.Column),
			))
		}
	}

	return erd.String()
}
