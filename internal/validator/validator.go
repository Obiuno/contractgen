package validator

import (
	"contractgen/internal/parser"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type Severity int

const (
	SeverityError Severity = iota
	SeverityWarning
)

func (s Severity) String() string {
	switch s {
	case SeverityError:
		return "error"
	case SeverityWarning:
		return "warning"
	default:
		return "unknown"
	}

}

type ValidationIssue struct {
	Severity Severity
	Message  string
	Table    string
	Column   string
}

func (v ValidationIssue) Error() string {
	location := v.Table
	if v.Column != "" {
		location = v.Table + "." + v.Column
	}
	if location == "" {
		return fmt.Sprintf("[%s] %s", v.Severity, v.Message)
	}
	return fmt.Sprintf("[%s] %s: %s", v.Severity, location, v.Message)
}

func Partition(issues []ValidationIssue) (errs, warnings []ValidationIssue) {
	for _, i := range issues {
		switch i.Severity {
		case SeverityError:
			errs = append(errs, i)
		case SeverityWarning:
			warnings = append(warnings, i)
		default:
			log.Printf("unknown severity %v on issue %v", i.Severity, i)
		}
	}
	return
}

func Validate(contract parser.Contract) []ValidationIssue {

	if errs := checkValidCharacters(contract); len(errs) > 0 {
		return errs
	}
	var errs []ValidationIssue

	errs = append(errs, checkDuplicateTables(contract)...)
	errs = append(errs, checkDuplicateColumns(contract)...)
	errs = append(errs, checkForeignKeys(contract)...)

	return errs
}

var identifierRegex = regexp.MustCompile(`^[a-z_][a-z0-9_]*$`)

func checkValidCharacters(contract parser.Contract) []ValidationIssue {
	var errs []ValidationIssue

	// iterate through tables
	for _, table := range contract.Tables {
		if !identifierRegex.MatchString(table.Name) {
			errs = append(errs, ValidationIssue{
				Severity: SeverityError,
				Message:  "table name is not a valid snake_case identifier",
				Table:    table.Name,
			})
		}

		for _, col := range table.Columns {
			if !identifierRegex.MatchString(col.Name) {
				errs = append(errs, ValidationIssue{
					Severity: SeverityError,
					Message:  "column name is not a valid snake_case identifier",
					Table:    table.Name,
					Column:   col.Name,
				})
			}
		}
	}
	return errs
}

func checkDuplicateTables(contract parser.Contract) []ValidationIssue {
	var errs []ValidationIssue
	seen := make(map[string]bool)

	for _, table := range contract.Tables {
		name := strings.ToLower(table.Name)

		if seen[name] {
			errs = append(errs, ValidationIssue{
				Severity: SeverityError,
				Message:  "duplicate table definition",
				Table:    name,
			})
			continue
		}

		seen[name] = true

	}
	return errs
}

func checkDuplicateColumns(contract parser.Contract) []ValidationIssue {
	var errs []ValidationIssue

	for _, table := range contract.Tables {
		tName := strings.ToLower(table.Name)
		seen := make(map[string]bool)

		for _, col := range table.Columns {
			cName := strings.ToLower(col.Name)

			if seen[cName] {
				errs = append(errs, ValidationIssue{
					Severity: SeverityError,
					Message:  "duplicate column",
					Table:    tName,
					Column:   cName,
				})

				continue
			}
			seen[cName] = true
		}
	}

	return errs
}

func checkForeignKeys(contract parser.Contract) []ValidationIssue {
	var errs []ValidationIssue

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
				errs = append(errs, ValidationIssue{
					Severity: SeverityError,
					Message:  fmt.Sprintf("references unknown table %q", refTable),
					Table:    tName,
					Column:   cName,
				})
				continue
			}

			if !cols[refColumn] {
				errs = append(errs, ValidationIssue{
					Severity: SeverityError,
					Message:  fmt.Sprintf("references unknown column %q.%q", refTable, refColumn),
					Table:    tName,
					Column:   cName,
				})
			}
		}
	}

	return errs
}
