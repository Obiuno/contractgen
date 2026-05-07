package cli

import (
	"flag"
)

type Config struct {
	Input    string
	Output   string
	ShowJSON bool
}

func SetupFlags() Config {
	input := flag.String("input", "", "YAML contract Path")
	output := flag.String("output", "schema.sql", "SQL output path")
	showJSON := flag.Bool("json", false, "print JSON")
	flag.Parse()
	return Config{
		Input:    *input,
		Output:   *output,
		ShowJSON: *showJSON,
	}
}
