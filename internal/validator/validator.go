package validator

import (
	"contractgen/internal/parser"
	"fmt"
	"strings"
)

type ValidationError struct {
	Message string
	Table   string
	Column  string
}

func (v ValidationError) Error() string {
	location := v.Table
	if v.Column != "" {
		location = v.Table + "." + v.Column
	}

	if location == "" {
		return v.Message
	}
	return fmt.Sprintf("%s: %s", location, v.Message)
}

func Validate(contract parser.Contract) []ValidationError {
	var errs []ValidationError

	errs = append(errs, checkDuplicateTables(contract)...)
	errs = append(errs, checkDuplicateColumns(contract)...)
	errs = append(errs, checkForeignKeys(contract)...)

	return errs
}

func checkDuplicateTables(contract parser.Contract) []ValidationError {
	var errs []ValidationError
	seen := make(map[string]bool)

	for _, table := range contract.Tables {
		name := strings.ToLower(table.Name)

		if seen[name] {
			errs = append(errs, ValidationError{
				Message: "duplicate table definition",
				Table:   name,
			})
			continue
		}

		seen[name] = true

	}
	return errs
}

func checkDuplicateColumns(contract parser.Contract) []ValidationError {
	var errs []ValidationError

	for _, table := range contract.Tables {
		tName := strings.ToLower(table.Name)
		seen := make(map[string]bool)

		for _, col := range table.Columns {
			cName := strings.ToLower(col.Name)

			if seen[cName] {
				errs = append(errs, ValidationError{
					Message: "duplicate column",
					Table:   tName,
					Column:  cName,
				})

				continue
			}
			seen[cName] = true
		}
	}

	return errs
}

func checkForeignKeys(contract parser.Contract) []ValidationError {
	var errs []ValidationError

	// build table and columns as table maps to column names
	tables := make(map[string]map[string]bool)

	for _, table := range contract.Tables {
		tName := strings.ToLower(table.Name)

		//make cols map
		cols := make(map[string]bool)

		for _, col := range table.Columns {
			cols[strings.ToLower(col.Name)] = true
		}
		tables[tName] = cols
	}

	//walk and check
	for _, table := range contract.Tables {
		tName := strings.ToLower(table.Name)

		for _, col := range table.Columns {
			if col.References == nil {
				continue
			}

			cName := strings.ToLower(col.Name)
			// reference tables and columns
			refTable := strings.ToLower(col.References.Table)
			refColumn := strings.ToLower(col.References.Column)

			cols, tableExists := tables[refTable]
			if !tableExists {
				errs = append(errs, ValidationError{
					Message: fmt.Sprintf("references unknown table %q", refTable),
					Table: tName,
					Column: cName,
				})
			continue
			}

			if !cols[refColumn] {
				errs = append(errs, ValidationError{
					Message: fmt.Sprintf("references unknown column %q.%q", refTable, refColumn),
					Table: tName,
					Column: cName,
				})
			}
		}
	}

	return errs
}
