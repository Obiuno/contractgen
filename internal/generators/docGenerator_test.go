package generators

import (
	"strings"
	"testing"

	"contractgen/internal/parser"
)

func TestGenerateDocs(t *testing.T) {
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
            "# Contents",
            "- [users](#users)",
            "# Tables",
            "## users",
            "| Column | Type | Nullable | Key | References | Description |",
            "| id |",
            "| email |",
        },
    },
    {
        name: "table with descriptions and constraints",
        contract: parser.Contract{
            Tables: []parser.Table{
                {
                    Name:        "users",
                    Description: "User accounts",
                    PrimaryKey:  []string{"id"},
                    Columns: []parser.Column{
                        {Name: "id", Type: "uuid", Description: "Primary key"},
                        {Name: "email", Type: "varchar", Constraints: []string{"not_null"}, Description: "Email address"},
                    },
                },
            },
        },
        wantContains: []string{
            "## users",
            "User accounts",            
            "Primary key",
            "Email address",
            "| PK |",
        },
    },
    {
        name: "multiple tables with FK reference",
        contract: parser.Contract{
            Tables: []parser.Table{
                {Name: "users", Columns: []parser.Column{
                    {Name: "id", Type: "uuid"},
                }, PrimaryKey: []string{"id"}},
                {Name: "orders", Columns: []parser.Column{
                    {Name: "user_id", Type: "uuid", References: &parser.Reference{Table: "users", Column: "id"}},
                }},
            },
        },
        wantContains: []string{
            "- [users](#users)",
            "- [orders](#orders)",
            "## users",
            "## orders",
            "| FK |",
            "users.id",
        },
    },
}


	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := GenerateDocs(tc.contract)
			for _, exp := range tc.wantContains {
				if !strings.Contains(out, exp) {
					t.Errorf("output missing %q\nfull output:\n%s", exp, out)
				}
			}
		})
	}
}