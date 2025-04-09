package graph

import "github.com/replay/platform/apps/reconstruction/internal/fetch"

// Node represents a service or topic in the workflow graph.
type Node struct {
	ID      string `json:"id"`
	Kind    string `json:"kind"`
	Topic   string `json:"topic,omitempty"`
	Failed  bool   `json:"failed,omitempty"`
}

// Edge links producer and consumer topics/services.
type Edge struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Graph is the workflow artifact.
type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// BuildGraph constructs nodes and edges from event metadata.
func BuildGraph(rows []fetch.EventRow) Graph {
	g := Graph{}
	AddEdges(&g, rows)
	return g
}

// AddEdges appends service-topic edges derived from events.
func AddEdges(g *Graph, rows []fetch.EventRow) {
	seenNode := map[string]bool{}
	seenEdge := map[string]bool{}
	for _, row := range rows {
		topicID := "topic:" + row.Topic
		if !seenNode[topicID] {
			g.Nodes = append(g.Nodes, Node{ID: topicID, Kind: "topic", Topic: row.Topic})
			seenNode[topicID] = true
		}
		if row.ServiceName != "" {
			svcID := "svc:" + row.ServiceName
			if !seenNode[svcID] {
				g.Nodes = append(g.Nodes, Node{ID: svcID, Kind: "service"})
				seenNode[svcID] = true
			}
			ek := svcID + "->" + topicID
			if !seenEdge[ek] {
				g.Edges = append(g.Edges, Edge{From: svcID, To: topicID})
				seenEdge[ek] = true
			}
		}
	}
}

// MarkFailed flags nodes with error outcomes for failure highlighting.
func MarkFailed(g *Graph, rows []fetch.EventRow) {
	failed := map[string]bool{}
	for _, row := range rows {
		if row.Outcome == "error" {
			failed["topic:"+row.Topic] = true
			if row.ServiceName != "" {
				failed["svc:"+row.ServiceName] = true
			}
		}
	}
	for i := range g.Nodes {
		if failed[g.Nodes[i].ID] {
			g.Nodes[i].Failed = true
		}
	}
}

// CascadePaths returns edges downstream of failed nodes.
func CascadePaths(g Graph) []Edge {
	failed := map[string]bool{}
	for _, n := range g.Nodes {
		if n.Failed {
			failed[n.ID] = true
		}
	}
	var out []Edge
	for _, e := range g.Edges {
		if failed[e.From] {
			out = append(out, e)
		}
	}
	return out
}
