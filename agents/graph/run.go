package graph

type Spec struct {
	Trigger bool `json:"trigger"` // If true, this node will be executed as a trigger node
	Done    bool `json:"done"`    // If true, this node will be executed as a done node
	Error   bool `json:"error"`   // If true, this node will be executed as an error node
}
type CompiledNode struct {
	ID    string                        `json:"id"`
	fnPtr func(StateSchema) StateSchema `json:"-"`
	Spec  Spec                          `json:"spec"` // The spec of the node
}

type CompiledGraph struct {
	ID            string         `json:"id"`
	CompiledNodes []CompiledNode `json:"compied_nodes"`
	EntryNode     string         `json:"entry_node"` // The entry node of the graph
}

func (gb *GraphBuilder) Compile(data StateSchema) *CompiledGraph {
	seenNodes := make(map[string]*Spec)
	compiledNodes := make([]CompiledNode, 0, len(gb.Nodes))

	for i, node := range gb.Nodes {
		trigger := false
		if i == 0 {
			trigger = true
		}
		seenNodes[node.ID] = &Spec{
			Trigger: trigger,
			Done:    false,
			Error:   false,
		}
		compiledNode := CompiledNode{
			ID:    node.ID,
			fnPtr: node.fnPtr,
			Spec:  *seenNodes[node.ID],
		}
		compiledNodes = append(compiledNodes, compiledNode)
	}
	return &CompiledGraph{
		ID:            gb.ID,
		CompiledNodes: compiledNodes,
		EntryNode:     gb.Nodes[0].ID,
	}

}

func NodeWrapFn(fnPtr func(StateSchema) StateSchema, data StateSchema, ch chan StateSchema) {
	result := fnPtr(data)
	ch <- result
}

func (cp *CompiledGraph) Run(data StateSchema) StateSchema {

	// for tick := range cp.Stream(data) {

	// }
	// return data
	return data
}

func (cp *CompiledGraph) Stream(data StateSchema, tickData chan StateSchema) {

	ch := make(chan StateSchema)

	start := cp.EntryNode
	for _, node := range cp.CompiledNodes {
		if node.ID == start {
			go NodeWrapFn(node.fnPtr, data, ch)
			break
		}
	}

	val := <-ch
	tickData <- val
	close(ch)

}
