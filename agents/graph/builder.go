package graph

// GraphBuilder
//
// A GraphBuilder is a structure that allows you to build a graph by adding nodes and edges.
type GraphBuilder struct {
	ID    string `json:"id"`
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// Node functions must accept one argument of type StateSchema and return a value of the same type.
type StateSchema interface{}

func NewGraphBuilder(id string, opts ...func(*GraphBuilder)) *GraphBuilder {
	gb := &GraphBuilder{
		ID:    id,
		Nodes: []Node{},
		Edges: []Edge{},
	}
	for _, opt := range opts {
		opt(gb)
	}
	return gb
}

func WithNode(
	nodename string,
	fn func(StateSchema) StateSchema,
) func(*GraphBuilder) {
	return func(gb *GraphBuilder) {
		node := Node{
			ID:    nodename,
			fnPtr: fn,
		}
		gb.Nodes = append(gb.Nodes, node)
	}
}

func WithEdge(
	leftnode string,
	rightnode string,
) func(*GraphBuilder) {
	return func(gb *GraphBuilder) {
		edge := Edge{
			ID:       leftnode + "-" + rightnode,
			SourceID: leftnode,
			TargetID: rightnode,
		}
		gb.Edges = append(gb.Edges, edge)
	}
}

// func WithConditionalEdge(
// 	leftnode string,
// 	pathFunc func(StateSchema) string,
// ) func(*GraphBuilder) {
// 	return func(gb *GraphBuilder) {
// 		edge := Edge{
// 			ID:       leftnode + "-" + rightnode,
// 			SourceID: leftnode,
// 			TargetID: rightnode,
// 			fnPtr:    fn,
// 		}
// 		gb.Edges = append(gb.Edges, edge)
// 	}
// }
