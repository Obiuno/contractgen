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

			// skip references to no-existent tables
			if _, exists := s.Tables[refTable]; !exists {
            continue
			}
			
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

	order, unprocessed := g.kahn(s.OrderedTables)
	if len(unprocessed) > 0 {
		return nil, fmt.Errorf("Cycle detected %s", strings.Join(unprocessed, ", "))
	}

	return order, nil
}

func (s *Schema) DetectCycles() ([]string) {
	g := s.buildGraph()
	_, unprocessed := g.kahn(s.OrderedTables)
	return unprocessed
}

// kahn's externally for topological order and cycle detection
// if unprocessed is empty, the order is complete i.e. there are no cycles
func (g *graph) kahn(orderedNodes []string) (order, unprocessed []string) {

	n := len(orderedNodes)
    queue := make([]string, 0, n)
    order = make([]string, 0, n)
    unprocessed = make([]string, 0, n)

	// fnd nodes with inDegree 0 and add to queue
	for _, name := range orderedNodes {
		if g.inDegree[name] == 0 {
			queue = append(queue, name)
		}
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		order = append(order, current)
		// for each table that depends on the popped item reduce inDegree

		for _, neighbor := range g.adj[current] {
			g.inDegree[neighbor]--
			//if inDgree 0 then add to queue
			if g.inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	for _, name := range orderedNodes {
		if g.inDegree[name] > 0 {
			unprocessed = append(unprocessed, name)
		}
	}

	return order, unprocessed

}
