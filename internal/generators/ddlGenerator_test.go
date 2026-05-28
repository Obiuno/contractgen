package generators

import (
	"strings"
	"testing"

	"contractgen/internal/parser"
)

func TestGenerateDDL(t *testing.T) {
	cases := []struct {
		name         string
		contract     parser.Contract
		wantContains []string
	}{
		{
			name: "minimal table",
			contract: parser.Contract{
				Tables: []parser.Table{
					{Name: "users", Columns: []parser.Column{
						{Name: "id", Type: "uuid"},
					}},
				},
			},
			wantContains: []string{
				"CREATE TABLE IF NOT EXISTS users",
				"id UUID",
				"BEGIN;",
				"COMMIT;",
			},
		},
		{
			name: "table with primary key",
			contract: parser.Contract{
				Tables: []parser.Table{
					{Name: "users", Columns: []parser.Column{
						{Name: "id", Type: "uuid"},
					}, PrimaryKey: []string{"id"}},
				},
			},
			wantContains: []string{
				"CREATE TABLE IF NOT EXISTS users",
				"PRIMARY KEY (id)",
				"SET NOT NULL", // PK implies NOT NULL in your generator
			},
		},
		{
			name: "table with not null constraint",
			contract: parser.Contract{
				Tables: []parser.Table{
					{Name: "users", Columns: []parser.Column{
						{Name: "id", Type: "uuid", Constraints: []string{"not_null"}},
					}},
				},
			},
			wantContains: []string{
				"CREATE TABLE IF NOT EXISTS users",
				"id UUID",
				"SET NOT NULL",
			},
		},
		{
			name: "two tables with FK",
			contract: parser.Contract{
				Tables: []parser.Table{
					{Name: "users", Columns: []parser.Column{
						{Name: "id", Type: "uuid"},
					}},
					{Name: "accounts", Columns: []parser.Column{
						{Name: "user_id", Type: "uuid", References: &parser.Reference{Table: "users", Column: "id"}},
					}},
				},
			},
			wantContains: []string{
				"CREATE TABLE IF NOT EXISTS users",
				"CREATE TABLE IF NOT EXISTS accounts",
				"ALTER TABLE accounts ADD CONSTRAINT",
				"FOREIGN KEY (user_id)",
				"REFERENCES users (id)",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := GenerateDDL(tc.contract)
			for _, exp := range tc.wantContains {
				if !strings.Contains(out, exp) {
					t.Errorf("output missing %q\nfull output:\n%s", exp, out)
				}
			}
		})
	}
}


