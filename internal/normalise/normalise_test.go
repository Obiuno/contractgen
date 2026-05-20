package normalise

import (
	"reflect"
	"testing"
	"contractgen/internal/parser"
)

func TestNormaliseName(t *testing.T) {
    cases := []struct {
        name   string
        input  string
        output string
    }{
        // Identity
        {name: "empty string", input: "", output: ""},
        {name: "single word", input: "foo", output: "foo"},
        {name: "already normalised", input: "bing_zap", output: "bing_zap"},

        // Casing
        {name: "lowercase", input: "Foo", output: "foo"},
        {name: "all caps", input: "BAR", output: "bar"},

        // Whitespace handling
        {name: "trim spaces", input: "  foo  ", output: "foo"},
        {name: "trim other whitespace", input: "\tqux\n", output: "qux"},
        {name: "space to underscore", input: "baz qux", output: "baz_qux"},
        {name: "collapse multiple spaces", input: "baz  qux", output: "baz_qux"},

        // Combinations
        {name: "trim + lowercase", input: "  Bar  ", output: "bar"},
        {name: "wombo combo", input: "  Foo Bar  ", output: "foo_bar"},

        // Pass-through (still invalid for the validator to reject)
        {name: "asterisks survive", input: "**foo**", output: "**foo**"},
        {name: "leading digit survives", input: "123foo", output: "123foo"},
        {name: "hyphens survive", input: "foo-bar", output: "foo-bar"},
        {name: "accented characters survive", input: "café", output: "café"},
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            result := normaliseName(tc.input)
            if result != tc.output {
                t.Errorf("expected %q got %q", tc.output, result)
            }
        })
    }
}

func TestNormalise(t *testing.T) {
    cases := []struct {
        name     string
        input    parser.Contract
        expected parser.Contract
    }{
        {
            name: "normalises table name",
            input: parser.Contract{
                Tables: []parser.Table{{Name: "Customers"}},
            },
            expected: parser.Contract{
                Tables: []parser.Table{{Name: "customers"}},
            },
        },
        {
            name: "normalises column names",
            input: parser.Contract{
                Tables: []parser.Table{
                    {Name: "users", Columns: []parser.Column{{Name: "User ID"}, {Name: "Email"}}},
                },
            },
            expected: parser.Contract{
                Tables: []parser.Table{
                    {Name: "users", Columns: []parser.Column{{Name: "user_id"}, {Name: "email"}}},
                },
            },
        },
        {
            name: "normalises references",
            input: parser.Contract{
                Tables: []parser.Table{
                    {Name: "orders", Columns: []parser.Column{
                        {Name: "customer_id", References: &parser.Reference{Table: "Customers", Column: "ID"}},
                    }},
                },
            },
            expected: parser.Contract{
                Tables: []parser.Table{
                    {Name: "orders", Columns: []parser.Column{
                        {Name: "customer_id", References: &parser.Reference{Table: "customers", Column: "id"}},
                    }},
                },
            },
        },
        {
            name: "normalises primary key columns",
            input: parser.Contract{
                Tables: []parser.Table{
                    {Name: "users", PrimaryKey: []string{"User ID"}},
                },
            },
            expected: parser.Contract{
                Tables: []parser.Table{
                    {Name: "users", PrimaryKey: []string{"user_id"}},
                },
            },
        },
        {
            name: "normalises unique constraints",
            input: parser.Contract{
                Tables: []parser.Table{
                    {Name: "users", Unique: [][]string{{"First Name", "Last Name"}}},
                },
            },
            expected: parser.Contract{
                Tables: []parser.Table{
                    {Name: "users", Unique: [][]string{{"first_name", "last_name"}}},
                },
            },
        },
        {
            name: "leaves descriptions untouched",
            input: parser.Contract{
                Tables: []parser.Table{
                    {Name: "users", Description: "  User Records  ", Columns: []parser.Column{
                        {Name: "id", Description: "Primary Key for the User"},
                    }},
                },
            },
            expected: parser.Contract{
                Tables: []parser.Table{
                    {Name: "users", Description: "  User Records  ", Columns: []parser.Column{
                        {Name: "id", Description: "Primary Key for the User"},
                    }},
                },
            },
        },
        {
            name: "no-op on already normalised",
            input: parser.Contract{
                Tables: []parser.Table{
                    {Name: "customers", Columns: []parser.Column{{Name: "id"}}},
                },
            },
            expected: parser.Contract{
                Tables: []parser.Table{
                    {Name: "customers", Columns: []parser.Column{{Name: "id"}}},
                },
            },
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            input := tc.input
            Normalise(&input)
            if !reflect.DeepEqual(input, tc.expected) {
                t.Errorf("Normalise() = %+v, want %+v", input, tc.expected)
            }
        })
    }
}
