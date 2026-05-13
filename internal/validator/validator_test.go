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

func containsErr(errs []ValidationError, table, column, msgContains string) bool {
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
			name: "errors for all 3 checks",
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
        err  ValidationError
        want string
    }{
        {
            name: "table and column",
            err: ValidationError{
                Message: "duplicate column",
                Table:   "customers",
                Column:  "id",
            },
            want: "customers.id: duplicate column",
        },
        {
            name: "table only",
            err: ValidationError{
                Message: "duplicate table definition",
                Table:   "customers",
            },
            want: "customers: duplicate table definition",
        },
        {
            name: "neither",
            err: ValidationError{
                Message: "something broke",
            },
            want: "something broke",
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