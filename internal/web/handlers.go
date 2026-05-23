package web

import (
	"fmt"
	"net/http"
	
	"contractgen/internal/generators"
	"contractgen/internal/normalise"
	"contractgen/internal/parser"
	"contractgen/internal/validator"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Example string
	}{
		Example: exampleYAML,
	}
	IndexTmpl.Execute(w, data)
}

const exampleYAML = `version: "1"
description: "Example yaml"
tables:
  - name: users
    description: "User accounts"
    columns:
      - name: id
        type: uuid
        constraints: [not_null]
      - name: email
        type: varchar
        constraints: [not_null]
    primary_key: [id]
    unique: [[email]]

  - name: posts
    description: "Blog posts"
    columns:
      - name: id
        type: uuid
        constraints: [not_null]
      - name: user_id
        type: uuid
        constraints: [not_null]
        references:
          table: users
          column: id
      - name: title
        type: varchar
        constraints: [not_null]
    primary_key: [id]
`

type ResultsData struct {
	DDL      string
	DataDict string
	ERD      string
	Issues   []validator.ValidationIssue
}

func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	yamlText := r.FormValue("yaml")

	data := ResultsData{}

	contract, err := parser.ParseContract([]byte(yamlText))

	// handle parsing issues
	if err != nil {
		data.Issues = []validator.ValidationIssue{
			{
				Severity: validator.SeverityError,
				Message:  fmt.Sprintf("invalid YAML %s", err.Error()),
			},
		}
		ResultsTmpl.Execute(w, data)
		return
	}

	normalise.Normalise(&contract)

	// validation for a working contract
	issues := validator.Validate(contract)
	data.Issues = issues

  // partition the issues to their severity
  errs, _ := validator.Partition(issues)
  if len(errs) > 0 {
    // errors block, warnings ignored
    ResultsTmpl.Execute(w, data)
    return
  }

  // no error issues, so can still generate
  data.DDL = generators.GenerateDDL(contract)
  data.DataDict = generators.GenerateDocs(contract)
  data.ERD = generators.GenerateMermaid(contract)

  ResultsTmpl.Execute(w, data)

}
