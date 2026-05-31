package schema

import (
	"slices"
	"strings"
	"testing"

	"contractgen/internal/parser"
)

func isTopologicallyValid(order []string, s *Schema) bool {
	position := make(map[string]int, len(order))
	for i, name := range order {
		position[name] = i
	}
	for name, table := range s.Tables {
		for _, colName := range table.OrderedColumns {
			col := table.Columns[colName]
			if col.References == nil {
				continue
			}
			ref := strings.ToLower(col.References.Table)
			if position[ref] >= position[name] {
				return false // referenced table didn't come before
			}
		}
	}
	return true
}

type tableDef struct {
	name    string
	columns []columnDef
}

type columnDef struct {
	name      string
	refTable  string // empty = no reference
	refColumn string
}

func makeSchema(tables ...tableDef) *Schema {
	s := &Schema{
		Tables:        make(map[string]*Table, len(tables)),
		OrderedTables: make([]string, 0, len(tables)),
	}
	for _, td := range tables {
		t := &Table{
			Columns:        make(map[string]*Column),
			OrderedColumns: make([]string, 0, len(td.columns)),
		}
		for _, cd := range td.columns {
			col := &Column{}
			if cd.refTable != "" {
				col.References = &parser.Reference{
					Table:  cd.refTable,
					Column: cd.refColumn,
				}
			}
			t.Columns[cd.name] = col
			t.OrderedColumns = append(t.OrderedColumns, cd.name)
		}
		s.Tables[td.name] = t
		s.OrderedTables = append(s.OrderedTables, td.name)
	}
	return s
}

func TestTopologicalOrder(t *testing.T) {
	cases := []struct {
		name      string
		schema    *Schema
		wantErr   bool
		wantCycle []string
	}{
		{
			name:      "empty schema",
			schema:    &Schema{},
			wantErr:   false,
			wantCycle: nil,
		},
		{
			name: "one table no FKs",
			schema: makeSchema(
				tableDef{name: "customers", columns: []columnDef{{name: "id"}}},
			),
			wantErr: false,
		},
		{
			name: "two tables one FK",
			schema: makeSchema(
				tableDef{name: "customers", columns: []columnDef{{name: "id"}}},
				tableDef{name: "orders", columns: []columnDef{
					{name: "id"},
					{name: "customer_id", refTable: "customers", refColumn: "id"},
				}},
			),
			wantErr: false,
		},
		{
			name: "three tables no FKs",
			schema: makeSchema(
				tableDef{name: "customers", columns: []columnDef{{name: "id"}}},
				tableDef{name: "products", columns: []columnDef{{name: "id"}}},
				tableDef{name: "regions", columns: []columnDef{{name: "id"}}},
			),
			wantErr: false,
		},
		{
			name: "realistic schema with diamond",
			schema: makeSchema(
				tableDef{name: "customers", columns: []columnDef{{name: "id"}}},
				tableDef{name: "products", columns: []columnDef{{name: "id"}}},
				tableDef{name: "orders", columns: []columnDef{
					{name: "id"},
					{name: "customer_id", refTable: "customers", refColumn: "id"},
				}},
				tableDef{name: "order_items", columns: []columnDef{
					{name: "order_id", refTable: "orders", refColumn: "id"},
					{name: "product_id", refTable: "products", refColumn: "id"},
				}},
			),
			wantErr: false,
		},
		{
			name: "self-referential cycle",
			schema: makeSchema(
				tableDef{name: "employee", columns: []columnDef{
					{name: "id"},
					{name: "manager_id", refTable: "employee", refColumn: "id"},
				}},
			),
			wantErr:   true,
			wantCycle: []string{"employee"},
		},
		{
			name: "two-table cycle",
			schema: makeSchema(
				tableDef{name: "a", columns: []columnDef{
					{name: "id"},
					{name: "b_id", refTable: "b", refColumn: "id"},
				}},
				tableDef{name: "b", columns: []columnDef{
					{name: "id"},
					{name: "a_id", refTable: "a", refColumn: "id"},
				}},
			),
			wantErr:   true,
			wantCycle: []string{"a", "b"},
		},
		{
			name: "three-table cycle",
			schema: makeSchema(
				tableDef{name: "a", columns: []columnDef{
					{name: "id"},
					{name: "b_id", refTable: "b", refColumn: "id"},
				}},
				tableDef{name: "b", columns: []columnDef{
					{name: "id"},
					{name: "c_id", refTable: "c", refColumn: "id"},
				}},
				tableDef{name: "c", columns: []columnDef{
					{name: "id"},
					{name: "a_id", refTable: "a", refColumn: "id"},
				}},
			),
			wantErr:   true,
			wantCycle: []string{"a", "b", "c"},
		},
		{
			name: "cycle alongside valid chain",
			schema: makeSchema(
				tableDef{name: "a", columns: []columnDef{
					{name: "id"},
					{name: "b_id", refTable: "b", refColumn: "id"},
				}},
				tableDef{name: "b", columns: []columnDef{
					{name: "id"},
					{name: "a_id", refTable: "a", refColumn: "id"},
				}},
				tableDef{name: "c", columns: []columnDef{{name: "id"}}},
				tableDef{name: "d", columns: []columnDef{
					{name: "id"},
					{name: "c_id", refTable: "c", refColumn: "id"},
				}},
			),
			wantErr:   true,
			wantCycle: []string{"a", "b"}, // c and d are fine, only a/b are stuck
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.schema.TopologicalOrder()

			if tc.wantErr {
				// Cycle case
				if err == nil {
					t.Errorf("expected cycle error, expected: %v, got: %v", tc.wantCycle, result)
					return
				}
				for _, table := range tc.wantCycle {
					if !strings.Contains(err.Error(), table) {
						t.Errorf("error %q does not mention table %q", err.Error(), table)
					}
				}
			} else {
				// Success case
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if !isTopologicallyValid(result, tc.schema) {
					t.Errorf("ordering %v is not topologically valid", result)
				}
			}
		})
	}
}

func TestTopologicalOrderDeterminism(t *testing.T) {
	s := makeSchema(
		tableDef{name: "customers", columns: []columnDef{{name: "id"}}},
		tableDef{name: "products", columns: []columnDef{{name: "id"}}},
		tableDef{name: "orders", columns: []columnDef{
			{name: "id"},
			{name: "customer_id", refTable: "customers", refColumn: "id"},
		}},
		tableDef{name: "order_items", columns: []columnDef{
			{name: "order_id", refTable: "orders", refColumn: "id"},
			{name: "product_id", refTable: "products", refColumn: "id"},
		}},
	)
	// apperently map iterations are random on purpose so don't have to do it myself

	first, err := s.TopologicalOrder()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i := 0; i < 20; i++ {
		result, err := s.TopologicalOrder()
		if err != nil {
			t.Fatalf("run %d: unexpected error: %v", i, err)
		}
		if !slices.Equal(result, first) {
			t.Errorf("run %d: got %v, want %v (determinism violated)", i, result, first)
		}
	}

}

func TestDetectCycles(t *testing.T) {
    cases := []struct {
        name      string
        schema    *Schema
        wantCycle []string
    }{
        {
            name:      "empty schema",
            schema:    &Schema{},
            wantCycle: nil,
        },
        {
            name: "no cycles",
            schema: makeSchema(
                tableDef{name: "customers", columns: []columnDef{{name: "id"}}},
                tableDef{name: "orders", columns: []columnDef{
                    {name: "id"},
                    {name: "customer_id", refTable: "customers", refColumn: "id"},
                }},
            ),
            wantCycle: nil,
        },
        {
            name: "self-reference",
            schema: makeSchema(
                tableDef{name: "employee", columns: []columnDef{
                    {name: "id"},
                    {name: "manager_id", refTable: "employee", refColumn: "id"},
                }},
            ),
            wantCycle: []string{"employee"},
        },
        {
            name: "two-table cycle",
            schema: makeSchema(
                tableDef{name: "a", columns: []columnDef{
                    {name: "id"},
                    {name: "b_id", refTable: "b", refColumn: "id"},
                }},
                tableDef{name: "b", columns: []columnDef{
                    {name: "id"},
                    {name: "a_id", refTable: "a", refColumn: "id"},
                }},
            ),
            wantCycle: []string{"a", "b"},
        },
        {
            name: "broken FK reference is not a cycle",
            schema: makeSchema(
                tableDef{name: "orders", columns: []columnDef{
                    {name: "id"},
                    {name: "customer_id", refTable: "nonexistent", refColumn: "id"},
                }},
            ),
            wantCycle: nil,
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            cycled := tc.schema.DetectCycles()
            
            if len(cycled) != len(tc.wantCycle) {
                t.Errorf("got cycled tables %v, want %v", cycled, tc.wantCycle)
                return
            }
            
            for _, expected := range tc.wantCycle {
                if !slices.Contains(cycled, expected) {
                    t.Errorf("expected %q in cycled tables, got %v", expected, cycled)
                }
            }
        })
    }
}