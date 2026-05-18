package schema

import (
	"fmt"
	"strings"
)

type graph struct {
	adj      map[string][]string
	inDegree map[string]int
}

// buildGraph constructs the dependency graph from the schema's FKs.
// Edges point from referenced tables to referencing tables 
func (s *Schema) buildGraph() *graph {
	g := &graph{
		adj:      make(map[string][]string, len(s.Tables)),
		inDegree: make(map[string]int, len(s.Tables)),
	}
	// intialise each table to in-inDegree 0
	for _, name := range s.OrderedTables {
		g.inDegree[name] = 0
	}

	// walk each table and add the References

	for _, name := range s.OrderedTables {
		table := s.Tables[name]
		for _, colName := range table.OrderedColumns {
			col := table.Columns[colName]
			if col.References == nil {
				continue
			}
			refTable := strings.ToLower(col.References.Table)
			// edge refTable -> name
			g.adj[refTable] = append(g.adj[refTable], name)
			g.inDegree[name]++
		}
	}

	return g
}

// TopologicalOrder returns table names in dependency order:
// referenced tables come before tables that reference them.
// Returns an error if a cycle is detected.
func (s *Schema) TopologicalOrder() ([]string, error) {
	// Kahn's algo
	g := s.buildGraph()

	queue := []string{}
	// fnd nodes with inDegree 0 and add to queue
	for _, name := range s.OrderedTables {
		if g.inDegree[name] == 0 {
			queue = append(queue, name)
		}
	}

	result := []string{}
	// look at queue and pop first element and add to result
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		result = append(result, current)
		// for each table that depends on the popped item reduce inDegree

		for _, neighbor := range g.adj[current] {
			g.inDegree[neighbor]--
			//if inDgree 0 then add to queue
			if g.inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}

	}

	// check length of result against number of tables for cycles
	// collect tables and return error
	if len(result) != len(s.Tables) {
    var cycledTables []string
    for _, name := range s.OrderedTables {
        if g.inDegree[name] > 0 {
            cycledTables = append(cycledTables, name)
        }
    }
    return nil, fmt.Errorf("cycle detected in the following tables:\n%s", strings.Join(cycledTables, ", "))
}


	return result, nil
}
