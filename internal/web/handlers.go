package web

import (
	"log"
	"net/http"

	"contractgen/internal/generators"
	"contractgen/internal/parser"
	"contractgen/internal/schema"
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

	sch := schema.BuildSchema(contract)

	if _, err := sch.TopologicalOrder(); err != nil {
		log.Printf("topological order warning: %v", err)
	}

	data := struct {
		DDL string
		DataDict string
		ERD string
	}{
		DDL: generators.GenerateDDL(contract),
		DataDict: generators.GenerateDocs(contract),
		ERD: generators.GenerateMermaid(sch),
	}

	if err := ResultsTmpl.Execute(w, data); err != nil {
		log.Printf("results template: %v", err)
	}
}