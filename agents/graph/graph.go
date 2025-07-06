package graph

// Node structure represents a node in the graph.
// It contains an ID and a function pointer that takes a StateSchema and returns a StateSchema
type Node struct {
	ID         string                        `json:"id"`
	fnPtr      func(StateSchema) StateSchema `json:"-"`
	Entrypoint bool                          `json:"entrypoint"` // Indicates if this node is an entry point
	Endpoint   bool                          `json:"endpoint"`   // Indicates if this node is an endpoint
}

type Edge struct {
	ID       string `json:"id"`
	SourceID string `json:"source"`
	TargetID string `json:"target"`
}
