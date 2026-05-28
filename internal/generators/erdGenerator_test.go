package generators

import (
	"strings"
	"testing"

	"contractgen/internal/parser"
)

func TestGenerateMermaid(t *testing.T) {
	cases := []struct {
    name         string
    contract     parser.Contract
    wantContains []string
}{
    {
        name: "single table basic structure",
        contract: parser.Contract{
            Tables: []parser.Table{
                {Name: "users", Columns: []parser.Column{
                    {Name: "id", Type: "uuid"},
                    {Name: "email", Type: "varchar"},
                }},
            },
        },
        wantContains: []string{
            "erDiagram",
            "users {",
            "uuid id",
            "varchar email",
        },
    },
    {
        name: "table with primary key marker",
        contract: parser.Contract{
            Tables: []parser.Table{
                {Name: "users", PrimaryKey: []string{"id"}, Columns: []parser.Column{
                    {Name: "id", Type: "uuid"},
                }},
            },
        },
        wantContains: []string{
            "erDiagram",
            "users {",
            "uuid id PK",
        },
    },
    {
        name: "two tables with FK relationship",
        contract: parser.Contract{
            Tables: []parser.Table{
                {Name: "users", PrimaryKey: []string{"id"}, Columns: []parser.Column{
                    {Name: "id", Type: "uuid"},
                }},
                {Name: "orders", Columns: []parser.Column{
                    {Name: "user_id", Type: "uuid", References: &parser.Reference{Table: "users", Column: "id"}},
                }},
            },
        },
        wantContains: []string{
            "users {",
            "orders {",
            "uuid id PK",
            "uuid user_id FK",
            "users ||--o{ orders",
            "via orders.user_id",
        },
    },
    {
        name: "PK and FK on same column",
        contract: parser.Contract{
            Tables: []parser.Table{
                {Name: "users", PrimaryKey: []string{"id"}, Columns: []parser.Column{
                    {Name: "id", Type: "uuid"},
                }},
                {Name: "user_settings", PrimaryKey: []string{"user_id"}, Columns: []parser.Column{
                    {Name: "user_id", Type: "uuid", References: &parser.Reference{Table: "users", Column: "id"}},
                }},
            },
        },
        wantContains: []string{
            "uuid user_id PK,FK",
            "users ||--o{ user_settings",
        },
    },
}
	
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := GenerateMermaid(tc.contract)
			for _, exp := range tc.wantContains {
				if !strings.Contains(out, exp) {
					t.Errorf("output missing %q\nfull output:\n%s", exp, out)
				}
			}
		})
	}
}