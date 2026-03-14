package graph

import "github.com/example/teleport-cluster-digital-twin/internal/model"

// ShortestPath returns the shortest path (BFS) of node IDs between start and goal.
func (g *Graph) ShortestPath(startID string, goalID string, allowed map[model.RelationshipType]bool) []string {
	if startID == goalID {
		return []string{startID}
	}
	type item struct{ id string }
	queue := []item{{id: startID}}
	prev := map[string]string{}
	seen := map[string]bool{startID: true}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, e := range g.OutgoingEdges(cur.id) {
			if len(allowed) > 0 && !allowed[e.Type] {
				continue
			}
			if seen[e.ToID] {
				continue
			}
			seen[e.ToID] = true
			prev[e.ToID] = cur.id
			if e.ToID == goalID {
				return unwind(prev, startID, goalID)
			}
			queue = append(queue, item{id: e.ToID})
		}
	}
	return nil
}

func unwind(prev map[string]string, startID, goalID string) []string {
	path := []string{goalID}
	for cur := goalID; cur != startID; {
		p, ok := prev[cur]
		if !ok {
			return nil
		}
		path = append([]string{p}, path...)
		cur = p
	}
	return path
}
