package web

import (
	"log"
	"net/http"
	"strings"

	"contractgen/internal/generators"
	"contractgen/internal/parser"
	"contractgen/internal/validator"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if err := IndexTmpl.Execute(w, nil); err != nil {
		log.Printf("index template: %v", err)
	}
}

func GenerateHandler(w http.ResponseWriter, r *http.Request) {
	yamlText := r.FormValue("yaml")

	contract, err := parser.ParseContract([]byte(yamlText))

	if err != nil {
		http.Error(w, "invalid YAML: "+err.Error(), http.StatusBadRequest)
		return
	}

	if issues := validator.Validate(contract); len(issues) > 0 {
		msgs := make([]string, len(issues))
		for i, e := range issues {
			msgs[i] = e.Error()
		}
		http.Error(w, "validation failed:\n"+strings.Join(msgs, "\n"), http.StatusBadRequest)
		return
	}

	data := struct {
		DDL      string
		DataDict string
		ERD      string
	}{
		DDL:      generators.GenerateDDL(contract),
		DataDict: generators.GenerateDocs(contract),
		ERD:      generators.GenerateMermaid(contract),
	}

	if err := ResultsTmpl.Execute(w, data); err != nil {
		log.Printf("results template: %v", err)
	}
}
