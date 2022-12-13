package evilgo

type GraphEdge[T any] struct {
	Weight T
	From   *GraphNode[T]
	To     *GraphNode[T]
}

type GraphNode[T any] struct {
	NodeName string
}

type Graph[T any] struct {
	Nodes map[string]*GraphNode[T]
	Edges map[string][]*GraphEdge[T]
}

func BuildGraph[T any]() Graph[T] {
	return Graph[T]{
		Nodes: map[string]*GraphNode[T]{},
		Edges: map[string][]*GraphEdge[T]{},
	}
}

func (g Graph[T]) AddNode(name string) {
	if _, ok := g.Nodes[name]; ok {
		return
	}
	g.Nodes[name] = &GraphNode[T]{
		NodeName: name,
	}
}

func (g Graph[T]) HasNode(name string) bool {
	_, ok := g.Nodes[name]
	return ok
}

func (g Graph[T]) AddEdge(from string, to string, weight T) {
	e := GraphEdge[T]{
		Weight: weight,
	}

	if !g.HasNode(from) {
		g.AddNode(from)
	}
	e.From = g.Nodes[from]

	if !g.HasNode(to) {
		g.AddNode(to)
	}
	e.To = g.Nodes[to]

	g.Edges[from] = append(g.Edges[from], &e)
}
