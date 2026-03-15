package graph

import "github.com/example/teleport-cluster-digital-twin/internal/model"

// Traverse performs breadth-first traversal from a start node.
func (g *Graph) Traverse(startID string, maxDepth int, allowed map[model.RelationshipType]bool) []string {
	if maxDepth < 0 {
		return nil
	}
	type item struct {
		id    string
		depth int
	}
	queue := []item{{id: startID, depth: 0}}
	seen := map[string]bool{startID: true}
	order := []string{startID}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur.depth >= maxDepth {
			continue
		}
		for _, e := range g.OutgoingEdges(cur.id) {
			if len(allowed) > 0 && !allowed[e.Type] {
				continue
			}
			if seen[e.ToID] {
				continue
			}
			seen[e.ToID] = true
			order = append(order, e.ToID)
			queue = append(queue, item{id: e.ToID, depth: cur.depth + 1})
		}
	}
	return order
}
