package validator

import (
	"contractgen/internal/parser"
	"strings"
	"testing"
)

type expectedError struct {
	table           string
	column          string
	messageContains string // for column check set ""
}

// containsErr does set-membership, not multiset. For cases that expect
// N copies of the same error, the count check enforces N; this helper
// only verifies at least one error of the expected shape exists.

func containsErr(errs []ValidationIssue, table, column, msgContains string) bool {
	for _, e := range errs {
		if e.Table == table && e.Column == column {
			if msgContains == "" || strings.Contains(e.Message, msgContains) {
				return true
			}
		}
	}
	return false
}

func TestCheckDuplicateTables(t *testing.T) {
	// Cases cover:
	//   - empty contract
	//   - single table
	//   - multiple tables without duplicates
	//   - one duplicate pair
	//   - one duplicate within a longer list
	//   - case-insensitive duplicates
	cases := []struct {
		name      string
		contract  parser.Contract
		wantErrs  int
		wantTable string
	}{
		{
			name: "multiple tables no duplicates",
			contract: parser.Contract{
				Tables: []parser.Table{
					{Name: "customers"},
					{Name: "orders"},
				},
			},
			wantErrs:  0,
			wantTable: "",
		},
		{
			name: "duplicate table",
			contract: parser.Contract{
				Tables: []parser.Table{
					{Name: "customers"},
					{Name: "customers"},
				},
			},
			wantErrs:  1,
			wantTable: "customers",
		},
		{
			name: "no tables",
			contract: parser.Contract{
				Tables: []parser.Table{},
			},
			wantErrs:  0,
			wantTable: "",
		},
		{
			name: "one table no duplicates",
			contract: parser.Contract{
				Tables: []parser.Table{
					{Name: "customers"},
				},
			},
			wantErrs:  0,
			wantTable: "",
		},
		{
			name: "three tables, two duplicates",
			contract: parser.Contract{
				Tables: []parser.Table{
					{Name: "customers"},
					{Name: "customers"},
					{Name: "customers"},
				},
			},
			wantErrs:  2,
			wantTable: "customers",
		},
		{
			name: "case insensitivity",
			contract: parser.Contract{
				Tables: []parser.Table{
					{Name: "customers"},
					{Name: "Customers"},
				},
			},
			wantErrs:  1,
			wantTable: "customers",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := checkDuplicateTables(tc.contract)
			if len(errs) != tc.wantErrs {
				t.Errorf("got %d errors, expected %d: %v", len(errs), tc.wantErrs, errs)
				return
			}
			if tc.wantTable != "" && len(errs) > 0 && errs[0].Table != tc.wantTable {
				t.Errorf("error refers to table %q, expected %q", errs[0].Table, tc.wantTable)
			}
		})
	}
}

func TestCheckDuplicateColumns(t *testing.T) {
	cases := []struct {
		name     string
		contract parser.Contract
		want     []expectedError
	}{
		{
			name: "no duplicates",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "customers",
						Columns: []parser.Column{
							{Name: "id"},
							{Name: "email"},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "duplicate column one table",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "customers",
						Columns: []parser.Column{
							{Name: "id"},
							{Name: "id"},
						},
					},
				},
			},
			want: []expectedError{
				{table: "customers", column: "id"},
			},
		},
		{
			name: "three columns one duplicate",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "customers",
						Columns: []parser.Column{
							{Name: "id"},
							{Name: "id"},
							{Name: "name"},
						},
					},
				},
			},
			want: []expectedError{
				{table: "customers", column: "id"},
			},
		},
		{
			name: "three columns two duplicates",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "customers",
						Columns: []parser.Column{
							{Name: "id"},
							{Name: "id"},
							{Name: "id"},
						},
					},
				},
			},
			want: []expectedError{
				{table: "customers", column: "id"},
				{table: "customers", column: "id"},
			},
		},
		{
			name: "two tables, each with one duplicate",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "customers",
						Columns: []parser.Column{
							{Name: "id"},
							{Name: "id"},
						},
					},
					{
						Name: "orders",
						Columns: []parser.Column{
							{Name: "customer_id"},
							{Name: "customer_id"},
						},
					},
				},
			},
			want: []expectedError{
				{table: "customers", column: "id"},
				{table: "orders", column: "customer_id"},
			},
		},
		{
			name: "case insensitive",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "customers",
						Columns: []parser.Column{
							{Name: "id"},
							{Name: "ID"},
						},
					},
				},
			},
			want: []expectedError{
				{table: "customers", column: "id"},
			},
		},
		{
			name: "same column name 2 tables",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "customers",
						Columns: []parser.Column{
							{Name: "id"},
						},
					},
					{
						Name: "orders",
						Columns: []parser.Column{
							{Name: "id"},
						},
					},
				},
			},
			want: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := checkDuplicateColumns(tc.contract)
			if len(errs) != len(tc.want) {
				t.Errorf("got %d errors, expected %d: %v", len(errs), len(tc.want), errs)
				return
			}
			for _, exp := range tc.want {
				if !containsErr(errs, exp.table, exp.column, exp.messageContains) {
					t.Errorf("expected error for %s.%s not found, got %v", exp.table, exp.column, errs)
				}
			}
		})
	}

}

func TestCheckForeignKeys(t *testing.T) {
	// this populates with the table.column with reference not the target table
	cases := []struct {
		name     string
		contract parser.Contract
		want     []expectedError
	}{
		{
			name: "no tables",
			contract: parser.Contract{
				Tables: []parser.Table{},
			},
			want: []expectedError{},
		},
		{
			name: "two tables no FK",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "customers",
						Columns: []parser.Column{
							{Name: "id"},
						},
					},
					{
						Name: "orders",
						Columns: []parser.Column{
							{Name: "id"},
						},
					},
				},
			},
			want: []expectedError{},
		},
		{
			name: "two tables one valid FK",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "customers",
						Columns: []parser.Column{
							{Name: "id"},
						},
					},
					{
						Name: "orders",
						Columns: []parser.Column{
							{
								Name: "customer_id",
								References: &parser.Reference{
									Table:  "customers",
									Column: "id",
								},
							},
						},
					},
				},
			},
			want: []expectedError{},
		},
		{
			name: "FK to unknown table",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "orders",
						Columns: []parser.Column{
							{
								Name: "customer_id",
								References: &parser.Reference{
									Table:  "customers",
									Column: "id",
								},
							},
						},
					},
				},
			},
			want: []expectedError{
				{
					table:           "orders",
					column:          "customer_id",
					messageContains: "unknown table",
				},
			},
		},
		{
			name: "FK to unknown column",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "customers",
						Columns: []parser.Column{
							{Name: "id"},
						},
					},
					{
						Name: "orders",
						Columns: []parser.Column{
							{
								Name: "customer_id",
								References: &parser.Reference{
									Table:  "customers",
									Column: "customer_id",
								},
							},
						},
					},
				},
			},
			want: []expectedError{
				{
					table:           "orders",
					column:          "customer_id",
					messageContains: "unknown column",
				},
			},
		},
		{
			name: "three tables two valid FKs in one table",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "customers",
						Columns: []parser.Column{
							{Name: "id"},
						},
					},
					{
						Name: "orders",
						Columns: []parser.Column{
							{Name: "id"},
						},
					},
					{
						Name: "orders_items",
						Columns: []parser.Column{
							{
								Name: "customer_id",
								References: &parser.Reference{
									Table:  "customers",
									Column: "id",
								},
							},
							{
								Name: "order_id",
								References: &parser.Reference{
									Table:  "orders",
									Column: "id",
								},
							},
						},
					},
				},
			},
			want: []expectedError{},
		},
		{
			name: "three tables, one valid FK, one invalid FK",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "customers",
						Columns: []parser.Column{
							{Name: "id"},
						},
					},
					{
						Name: "orders",
						Columns: []parser.Column{
							{
								Name: "customer_id",
								References: &parser.Reference{
									Table:  "customers",
									Column: "customer_id",
								},
							},
							{Name: "id"},
						},
					},
					{
						Name: "orders_items",
						Columns: []parser.Column{
							{
								Name: "order_id",
								References: &parser.Reference{
									Table:  "orders",
									Column: "id",
								},
							},
						},
					},
				},
			},
			want: []expectedError{
				{
					table:           "orders",
					column:          "customer_id",
					messageContains: "unknown column",
				},
			},
		},
		{
			name: "three tables, one valid FK, one invalid column FK, one invalid table FK",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "customers",
						Columns: []parser.Column{
							{Name: "id"},
						},
					},
					{
						Name: "orders",
						Columns: []parser.Column{
							{
								Name: "customer_id",
								References: &parser.Reference{
									Table:  "customers",
									Column: "customer_id",
								},
							},
						},
					},
					{
						Name: "orders_items",
						Columns: []parser.Column{
							{
								Name: "order_id",
								References: &parser.Reference{
									Table:  "account",
									Column: "id",
								},
							},
							{
								Name: "customer_id",
								References: &parser.Reference{
									Table:  "customers",
									Column: "id",
								},
							},
						},
					},
				},
			},
			want: []expectedError{
				{
					table:           "orders",
					column:          "customer_id",
					messageContains: "unknown column",
				},
				{
					table:           "orders_items",
					column:          "order_id",
					messageContains: "unknown table",
				},
			},
		},
		{
			name: "Case insensitive",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "customers",
						Columns: []parser.Column{
							{Name: "id"},
						},
					},
					{
						Name: "orders",
						Columns: []parser.Column{
							{
								Name: "customer_id",
								References: &parser.Reference{
									Table:  "Customers",
									Column: "ID",
								},
							},
						},
					},
				},
			},
			want: []expectedError{},
		},
		{
			name: "Self-referential FK",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "employee",
						Columns: []parser.Column{
							{Name: "employee_id"},
							{
								Name: "manager_id",
								References: &parser.Reference{
									Table:  "employee",
									Column: "employee_id",
								},
							},
						},
					},
				},
			},
			want: []expectedError{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := checkForeignKeys(tc.contract)
			if len(errs) != len(tc.want) {
				t.Errorf("got %d errors, want %d: %v", len(errs), len(tc.want), errs)
				return
			}
			for _, exp := range tc.want {
				if !containsErr(errs, exp.table, exp.column, exp.messageContains) {
					t.Errorf("expected error for %s.%s containing %q not found, got %v",
						exp.table, exp.column, exp.messageContains, errs)
				}
			}
		})
	}

}

func TestValidate(t *testing.T) {
	cases := []struct {
		name     string
		contract parser.Contract
		want     []expectedError
	}{
		{
			name: "multiple validation errors",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "customers",
						Columns: []parser.Column{
							{Name: "id"},
							{Name: "id"},
						},
					},
					{
						Name: "orders",
						Columns: []parser.Column{
							{
								Name: "customer_id",
								References: &parser.Reference{
									Table:  "customers",
									Column: "customer_id",
								},
							},
						},
					},
					{
						Name: "orders_items",
						Columns: []parser.Column{
							{
								Name: "order_id",
								References: &parser.Reference{
									Table:  "account",
									Column: "id",
								},
							},
							{
								Name: "customer_id",
								References: &parser.Reference{
									Table:  "customers",
									Column: "id",
								},
							},
						},
					},
					{
						Name: "customers",
						Columns: []parser.Column{
							{Name: "id"},
						},
					},
				},
			},
			want: []expectedError{
				{
					table: "customers",
				},
				{
					table:  "customers",
					column: "id",
				},
				{
					table:           "orders",
					column:          "customer_id",
					messageContains: "unknown column",
				},
				{
					table:           "orders_items",
					column:          "order_id",
					messageContains: "unknown table",
				},
			},
		},
		{
			name: "clean contract produces no errors",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "customers",
						Columns: []parser.Column{
							{Name: "id"},
							{Name: "email"},
						},
					},
				},
			},
			want: []expectedError{},
		},
		{
			name: "invalid identifier short-circuits other checks",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "**bad**", 
						Columns: []parser.Column{
							{Name: "id"},
							{Name: "id"},
						},
					},
				},
			},
			want: []expectedError{
				{table: "**bad**", messageContains: "identifier"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := Validate(tc.contract)
			if len(errs) != len(tc.want) {
				t.Errorf("got %d errors, want %d: %v", len(errs), len(tc.want), errs)
				return
			}
			for _, exp := range tc.want {
				if !containsErr(errs, exp.table, exp.column, exp.messageContains) {
					t.Errorf("expected error for %s.%s containing %q not found, got %v",
						exp.table, exp.column, exp.messageContains, errs)
				}
			}
		})
	}

}

func TestValidationErrorMessage(t *testing.T) {
	cases := []struct {
		name string
		err  ValidationIssue
		want string
	}{
		{
			name: "table and column",
			err: ValidationIssue{
				Message: "duplicate column",
				Table:   "customers",
				Column:  "id",
			},
			want: "[error] customers.id: duplicate column",
		},
		{
			name: "table only",
			err: ValidationIssue{
				Message: "duplicate table definition",
				Table:   "customers",
			},
			want: "[error] customers: duplicate table definition",
		},
		{
			name: "neither",
			err: ValidationIssue{
				Message: "something broke",
			},
			want: "[error] something broke",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.err.Error()
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestCheckValidCharacters(t *testing.T) {
	cases := []struct {
		name     string
		contract parser.Contract
		want     []expectedError
	}{
		{
			name: "valid identifier",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "customers",
						Columns: []parser.Column{
							{Name: "id"},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "empty table name",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "",
						Columns: []parser.Column{
							{Name: "id"},
						},
					},
				},
			},
			want: []expectedError{
				{table: "", messageContains: "identifier"},
			},
		},
		{
			name: "starts with digit",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "123foo",
					},
					{
						Name: "a",
						Columns: []parser.Column{
							{
								Name: "456bar",
							},
						},
					},
				},
			},
			want: []expectedError{
				{table: "123foo", messageContains: "identifier"},
				{table: "a", column: "456bar", messageContains: "identifier"},
			},
		},
		{
			name: "valid snake_case",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "user_accounts",
					},
					{
						Name: "a",
						Columns: []parser.Column{
							{
								Name: "first_name",
							},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "leading underscore allowed",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "_private",
					},
					{
						Name: "a",
						Columns: []parser.Column{
							{
								Name: "_id",
							},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "digits at end",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "a2",
					},
					{
						Name: "b",
						Columns: []parser.Column{
							{
								Name: "c2",
							},
						},
					},
				},
			},
			want: nil,
		},
		{
			name: "empty column name",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "foo",
						Columns: []parser.Column{
							{Name: ""},
						},
					},
				},
			},
			want: []expectedError{
				{table: "foo", column: "", messageContains: "identifier"},
			},
		},
		{
			name: "contains space",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "foo bar",
					},
					{
						Name: "a",
						Columns: []parser.Column{
							{
								Name: "baz qux",
							},
						},
					},
				},
			},
			want: []expectedError{
				{table: "foo bar", messageContains: "identifier"},
				{table: "a", column: "baz qux", messageContains: "identifier"},
			},
		},
		{
			name: "contains hyphen",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "foo-bar",
					},
					{
						Name: "a",
						Columns: []parser.Column{
							{
								Name: "baz-qux",
							},
						},
					},
				},
			},
			want: []expectedError{
				{table: "foo-bar", messageContains: "identifier"},
				{table: "a", column: "baz-qux", messageContains: "identifier"},
			},
		},
		{
			name: "uppercase rejected",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "Foo",
					},
					{
						Name: "a",
						Columns: []parser.Column{
							{
								Name: "Bar",
							},
						},
					},
				},
			},
			want: []expectedError{
				{table: "Foo", messageContains: "identifier"},
				{table: "a", column: "Bar", messageContains: "identifier"},
			},
		},
		{
			name: "asterisks rejected",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "**foobar**",
					},
					{
						Name: "a",
						Columns: []parser.Column{
							{
								Name: "**bazqux**",
							},
						},
					},
				},
			},
			want: []expectedError{
				{table: "**foobar**", messageContains: "identifier"},
				{table: "a", column: "**bazqux**", messageContains: "identifier"},
			},
		},
		{
			name: "accents rejected",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "café",
					},
					{
						Name: "a",
						Columns: []parser.Column{
							{
								Name: "café",
							},
						},
					},
				},
			},
			want: []expectedError{
				{table: "café", messageContains: "identifier"},
				{table: "a", column: "café", messageContains: "identifier"},
			},
		},
		{
			name: "multiple invalid tables",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "**foobar**",
					},
					{
						Name: "",
						Columns: []parser.Column{
							{
								Name: "bar",
							},
						},
					},
				},
			},
			want: []expectedError{
				{table: "**foobar**", messageContains: "identifier"},
				{table: "", messageContains: "identifier"},
			},
		},
		{
			name: "table and column name invalid",
			contract: parser.Contract{
				Tables: []parser.Table{
					{
						Name: "Foo",
						Columns: []parser.Column{
							{
								Name: "**bazqux**",
							},
						},
					},
				},
			},
			want: []expectedError{
				{table: "Foo", messageContains: "identifier"},
				{table: "Foo", column: "**bazqux**", messageContains: "identifier"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := checkValidCharacters(tc.contract)
			if len(errs) != len(tc.want) {
				t.Errorf("got %d errors, want %d: %v", len(errs), len(tc.want), errs)
				return
			}
			for _, exp := range tc.want {
				if !containsErr(errs, exp.table, exp.column, exp.messageContains) {
					t.Errorf("expected error for %s.%s containing %q not found, got %v",
						exp.table, exp.column, exp.messageContains, errs)
				}
			}
		})
	}
}

func TestSeverityPerCheck(t *testing.T) {
	cases := []struct {
		name string
		run  func() []ValidationIssue
		want Severity
	}{
		{
			name: "duplicate tables produce errors",
			run: func() []ValidationIssue {
				return checkDuplicateTables(parser.Contract{
					Tables: []parser.Table{{Name: "x"}, {Name: "x"}},
				})
			},
			want: SeverityError,
		},
		{
			name: "duplicate columns produce errors",
			run: func() []ValidationIssue {
				return checkDuplicateColumns(parser.Contract{
					Tables: []parser.Table{{Name: "x", Columns: []parser.Column{{Name: "id"}, {Name: "id"}}}},
				})
			},
			want: SeverityError,
		},
		{
			name: "foreign key issues produce errors",
			run: func() []ValidationIssue {
				return checkForeignKeys(parser.Contract{
					Tables: []parser.Table{
						{Name: "orders", Columns: []parser.Column{
							{Name: "customer_id", References: &parser.Reference{Table: "ghost", Column: "id"}},
						}},
					},
				})
			},
			want: SeverityError,
		},
		{
			name: "invalid characters produce errors",
			run: func() []ValidationIssue {
				return checkValidCharacters(parser.Contract{
					Tables: []parser.Table{{Name: "**bad**"}},
				})
			},
			want: SeverityError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			issues := tc.run()
			if len(issues) == 0 {
				t.Fatal("expected at least one issue, got none")
			}
			for _, issue := range issues {
				if issue.Severity != tc.want {
					t.Errorf("got severity %v, want %v in issue: %v", issue.Severity, tc.want, issue)
				}
			}
		})
	}
}

func TestPartition_UnknownSeverity(t *testing.T) {
    issues := []ValidationIssue{
        {Severity: Severity(99), Message: "unknown"},
    }
    errs, warns := Partition(issues)
    if len(errs) != 0 || len(warns) != 0 {
        t.Errorf("expected unknown severity to be dropped, got %d errors and %d warnings", len(errs), len(warns))
    }
}

func TestPartition(t *testing.T) {
	cases := []ValidationIssue{
		{Severity: SeverityError, Message: "e1"},
		{Severity: SeverityWarning, Message: "w1"},
		{Severity: SeverityError, Message: "e2"},
		{Severity: SeverityWarning, Message: "w2"},
	}
	errs, warns := Partition(cases)
	if len(errs) != 2 || len(warns) != 2 {
		t.Errorf("got %d errors and %d warnings, expected 2 and 2", len(errs), len(warns))
	}
}
