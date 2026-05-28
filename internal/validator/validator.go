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

	if errs := checkContractStructure(contract); len(errs) > 0 {
		return errs
	}

	if errs := checkValidCharacters(contract); len(errs) > 0 {
		return errs
	}
	var errs []ValidationIssue

	errs = append(errs, checkCompleteness(contract)...)
	errs = append(errs, checkDuplicateTables(contract)...)
	errs = append(errs, checkDuplicateColumns(contract)...)
	errs = append(errs, checkForeignKeys(contract)...)
	errs = append(errs, checkConstraintReferences(contract)...)

	return errs
}

var identRegex = regexp.MustCompile(`^[a-z_][a-z0-9_]*$`)

func checkValidCharacters(contract parser.Contract) []ValidationIssue {
	var errs []ValidationIssue

	for _, table := range contract.Tables {
		if !identRegex.MatchString(table.Name) {
			errs = append(errs, ValidationIssue{
				Severity: SeverityError,
				Message:  "table name is not a valid snake_case identifier",
				Table:    table.Name,
			})
		}

		for _, col := range table.Columns {
			if !identRegex.MatchString(col.Name) {
				errs = append(errs, ValidationIssue{
					Severity: SeverityError,
					Message:  "column name is not a valid snake_case identifier",
					Table:    table.Name,
					Column:   col.Name,
				})
			}

			// Reference fields, if present
			if col.References != nil {
				if col.References.Table != "" && !identRegex.MatchString(col.References.Table) {
					errs = append(errs, ValidationIssue{
						Severity: SeverityError,
						Message:  fmt.Sprintf("reference table %q is not a valid identifier", col.References.Table),
						Table:    table.Name,
						Column:   col.Name,
					})
				}
				if col.References.Column != "" && !identRegex.MatchString(col.References.Column) {
					errs = append(errs, ValidationIssue{
						Severity: SeverityError,
						Message:  fmt.Sprintf("reference column %q is not a valid identifier", col.References.Column),
						Table:    table.Name,
						Column:   col.Name,
					})
				}
			}
		}

		// Primary key columns
		for _, pk := range table.PrimaryKey {
			if !identRegex.MatchString(pk) {
				errs = append(errs, ValidationIssue{
					Severity: SeverityError,
					Message:  fmt.Sprintf("primary key entry %q is not a valid identifier", pk),
					Table:    table.Name,
				})
			}
		}

		// Unique constraint columns
		for _, group := range table.Unique {
			for _, col := range group {
				if !identRegex.MatchString(col) {
					errs = append(errs, ValidationIssue{
						Severity: SeverityError,
						Message:  fmt.Sprintf("unique constraint entry %q is not a valid identifier", col),
						Table:    table.Name,
					})
				}
			}
		}
	}

	return errs
}

func checkContractStructure(contract parser.Contract) []ValidationIssue {
	if len(contract.Tables) == 0 {
		return []ValidationIssue{{
			Severity: SeverityError,
			Message:  "contract has no tables",
		}}
	}
	return nil
}

func checkCompleteness(contract parser.Contract) []ValidationIssue {
	var errs []ValidationIssue

	for _, table := range contract.Tables {
		tName := table.Name

		// Table-level: must have at least one column
		if len(table.Columns) == 0 {
			errs = append(errs, ValidationIssue{
				Severity: SeverityError,
				Message:  "table has no columns",
				Table:    tName,
			})
			continue // no point checking individual columns
		}

		for _, col := range table.Columns {
			// Column name (might already be caught by checkValidCharacters
			if strings.TrimSpace(col.Name) == "" {
				errs = append(errs, ValidationIssue{
					Severity: SeverityError,
					Message:  "column has no name",
					Table:    tName,
				})
				continue // can't check type without a name
			}

			// Column type
			if strings.TrimSpace(col.Type) == "" {
				errs = append(errs, ValidationIssue{
					Severity: SeverityError,
					Message:  "column has no type",
					Table:    tName,
					Column:   col.Name,
				})
			}

			// If references block exists, both fields must be specified
			if col.References != nil {
				if strings.TrimSpace(col.References.Table) == "" {
					errs = append(errs, ValidationIssue{
						Severity: SeverityError,
						Message:  "reference is missing table",
						Table:    tName,
						Column:   col.Name,
					})
				}
				if strings.TrimSpace(col.References.Column) == "" {
					errs = append(errs, ValidationIssue{
						Severity: SeverityError,
						Message:  "reference is missing column",
						Table:    tName,
						Column:   col.Name,
					})
				}
			}
		}
	}

	return errs
}

func checkConstraintReferences(contract parser.Contract) []ValidationIssue {
	var errs []ValidationIssue

	for _, table := range contract.Tables {
		tName := strings.ToLower(table.Name)

		// Build the set of valid column names for this table
		cols := make(map[string]bool)
		for _, col := range table.Columns {
			cols[strings.ToLower(col.Name)] = true
		}

		// Check primary key references
		for _, pk := range table.PrimaryKey {
			pkLower := strings.ToLower(pk)
			if !cols[pkLower] {
				errs = append(errs, ValidationIssue{
					Severity: SeverityError,
					Message:  fmt.Sprintf("primary key references unknown column %q", pk),
					Table:    tName,
				})
			}
		}

		// Check unique constraint references
		for _, group := range table.Unique {
			for _, col := range group {
				colLower := strings.ToLower(col)
				if !cols[colLower] {
					errs = append(errs, ValidationIssue{
						Severity: SeverityError,
						Message:  fmt.Sprintf("unique constraint references unknown column %q", col),
						Table:    tName,
					})
				}
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
