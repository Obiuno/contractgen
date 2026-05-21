package schema

import (
	"slices"
	"testing"

	"contractgen/internal/parser"
)

func TestBuildSchema(t *testing.T) {
    cases := []struct {
        name     string
        contract parser.Contract
        check    func(*testing.T, *Schema)
    }{
        {
            name: "single table single column",
            contract: parser.Contract{
                Version: "1",
                Tables: []parser.Table{
                    {Name: "users", Columns: []parser.Column{
                        {Name: "id", Type: "uuid"},
                    }},
                },
            },
            check: func(t *testing.T, s *Schema) {
                if s.Version != "1" {
                    t.Errorf("version: got %q, want %q", s.Version, "1")
                }
                if len(s.Tables) != 1 {
                    t.Errorf("table count: got %d, want 1", len(s.Tables))
                }
                if s.Tables["users"] == nil {
                    t.Fatalf("table 'users' missing")
                }
                if s.Tables["users"].Columns["id"] == nil {
                    t.Fatalf("column 'id' missing")
                }
                if s.Tables["users"].Columns["id"].Type != "uuid" {
                    t.Errorf("type: got %q, want 'uuid'", s.Tables["users"].Columns["id"].Type)
                }
            },
        },
        {
            name: "preserves table declaration order",
            contract: parser.Contract{
                Tables: []parser.Table{
                    {Name: "zebra"},
                    {Name: "alpha"},
                    {Name: "mango"},
                },
            },
            check: func(t *testing.T, s *Schema) {
                want := []string{"zebra", "alpha", "mango"}
                if !slices.Equal(s.OrderedTables, want) {
                    t.Errorf("ordered tables: got %v, want %v", s.OrderedTables, want)
                }
            },
        },
        {
            name: "preserves column declaration order",
            contract: parser.Contract{
                Tables: []parser.Table{
                    {Name: "users", Columns: []parser.Column{
                        {Name: "zebra", Type: "text"},
                        {Name: "alpha", Type: "text"},
                        {Name: "mango", Type: "text"},
                    }},
                },
            },
            check: func(t *testing.T, s *Schema) {
                want := []string{"zebra", "alpha", "mango"}
                got := s.Tables["users"].OrderedColumns
                if !slices.Equal(got, want) {
                    t.Errorf("ordered columns: got %v, want %v", got, want)
                }
            },
        },
        {
            name: "lowercases table and column names",
            contract: parser.Contract{
                Tables: []parser.Table{
                    {Name: "Users", Columns: []parser.Column{
                        {Name: "ID", Type: "uuid"},
                        {Name: "Email", Type: "varchar"},
                    }},
                },
            },
            check: func(t *testing.T, s *Schema) {
                if s.Tables["users"] == nil {
                    t.Errorf("expected table key 'users', got keys: %v", mapKeys(s.Tables))
                }
                if s.Tables["users"].Columns["id"] == nil {
                    t.Errorf("expected column key 'id'")
                }
                if s.Tables["users"].Columns["email"] == nil {
                    t.Errorf("expected column key 'email'")
                }
                if !slices.Contains(s.OrderedTables, "users") {
                    t.Errorf("ordered tables missing lowercase 'users': %v", s.OrderedTables)
                }
            },
        },
        {
            name: "lowercases primary key columns",
            contract: parser.Contract{
                Tables: []parser.Table{
                    {Name: "users", PrimaryKey: []string{"ID"}, Columns: []parser.Column{
                        {Name: "id", Type: "uuid"},
                    }},
                },
            },
            check: func(t *testing.T, s *Schema) {
                want := []string{"id"}
                if !slices.Equal(s.Tables["users"].PrimaryKey, want) {
                    t.Errorf("primary key: got %v, want %v", s.Tables["users"].PrimaryKey, want)
                }
            },
        },
        {
            name: "lowercases unique constraints",
            contract: parser.Contract{
                Tables: []parser.Table{
                    {Name: "users", Unique: [][]string{{"First_Name", "Last_Name"}}, Columns: []parser.Column{
                        {Name: "first_name", Type: "varchar"},
                        {Name: "last_name", Type: "varchar"},
                    }},
                },
            },
            check: func(t *testing.T, s *Schema) {
                want := [][]string{{"first_name", "last_name"}}
                got := s.Tables["users"].Unique
                if len(got) != 1 || !slices.Equal(got[0], want[0]) {
                    t.Errorf("unique: got %v, want %v", got, want)
                }
            },
        },
        {
            name: "preserves column references",
            contract: parser.Contract{
                Tables: []parser.Table{
                    {Name: "orders", Columns: []parser.Column{
                        {Name: "user_id", Type: "uuid", References: &parser.Reference{Table: "users", Column: "id"}},
                    }},
                },
            },
            check: func(t *testing.T, s *Schema) {
                ref := s.Tables["orders"].Columns["user_id"].References
                if ref == nil {
                    t.Fatalf("references missing")
                }
                if ref.Table != "users" || ref.Column != "id" {
                    t.Errorf("references: got %+v, want {Table: users, Column: id}", ref)
                }
            },
        },
        {
            name: "preserves column constraints",
            contract: parser.Contract{
                Tables: []parser.Table{
                    {Name: "users", Columns: []parser.Column{
                        {Name: "id", Type: "uuid", Constraints: []string{"not_null"}},
                    }},
                },
            },
            check: func(t *testing.T, s *Schema) {
                want := []string{"not_null"}
                got := s.Tables["users"].Columns["id"].Constraints
                if !slices.Equal(got, want) {
                    t.Errorf("constraints: got %v, want %v", got, want)
                }
            },
        },
        {
            name: "empty contract produces empty schema",
            contract: parser.Contract{},
            check: func(t *testing.T, s *Schema) {
                if s == nil {
                    t.Fatal("schema should not be nil")
                }
                if len(s.Tables) != 0 {
                    t.Errorf("expected empty tables, got %d", len(s.Tables))
                }
                if len(s.OrderedTables) != 0 {
                    t.Errorf("expected empty ordered tables, got %d", len(s.OrderedTables))
                }
            },
        },
        {
            name: "table with no columns",
            contract: parser.Contract{
                Tables: []parser.Table{
                    {Name: "empty_table"},
                },
            },
            check: func(t *testing.T, s *Schema) {
                tbl := s.Tables["empty_table"]
                if tbl == nil {
                    t.Fatal("table missing")
                }
                if len(tbl.Columns) != 0 {
                    t.Errorf("expected no columns, got %d", len(tbl.Columns))
                }
            },
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            schema := BuildSchema(tc.contract)
            tc.check(t, schema)
        })
    }
}

// helper for error messages
func mapKeys[K comparable, V any](m map[K]V) []K {
    keys := make([]K, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}