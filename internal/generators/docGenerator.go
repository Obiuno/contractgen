package generators

import (
	"fmt"
	"slices"
	"strings"

	"contractgen/internal/parser"
)

func GenerateDocs(contract parser.Contract) string {
	var dataDict strings.Builder

	// ToC
	dataDict.WriteString("# Contents\n\n")
	for _, table := range contract.Tables {
		tName := strings.ToLower(table.Name)
		dataDict.WriteString(fmt.Sprintf("- [%s](#%s)\n", tName, tName))
	}
	dataDict.WriteString("\n# Tables\n\n")

	// One section per table
	for _, table := range contract.Tables {
		tName := strings.ToLower(table.Name)

		dataDict.WriteString(fmt.Sprintf("## %s\n\n", tName))
		if table.Description != "" {
			dataDict.WriteString(table.Description + "\n\n")
		}
		dataDict.WriteString("| Column | Type | Nullable | Key | References | Description |\n")
		dataDict.WriteString("|---|---|---|---|---|---|\n")

		pkColumns := make([]string, len(table.PrimaryKey))
		for i, col := range table.PrimaryKey {
			pkColumns[i] = strings.ToLower(col)
		}

		for _, col := range table.Columns {
			cName := strings.ToLower(col.Name)
			isPK := slices.Contains(pkColumns, cName)
			isFK := col.References != nil

			var refs string
			if isFK {
				refs = fmt.Sprintf("%s.%s",
					strings.ToLower(col.References.Table),
					strings.ToLower(col.References.Column))
			}

			var key string
			switch {
			case isPK && isFK:
				key = "PK, FK"
			case isPK:
				key = "PK"
			case isFK:
				key = "FK"
			}

			nullable := "Yes"
			if slices.Contains(col.Constraints, "not_null") || isPK {
				nullable = "No"
			}

			dataDict.WriteString(fmt.Sprintf(
				"| %s | %s | %s | %s | %s | %s |\n",
				cName,
				strings.ToLower(col.Type),
				nullable,
				key,
				refs,
				col.Description,
			))
		}
		dataDict.WriteString("\n")
	}

	return dataDict.String()
}
