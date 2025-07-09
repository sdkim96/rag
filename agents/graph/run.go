package graph

type Spec struct {
	Trigger bool     `json:"trigger"`
	Done    bool     `json:"done"`
	Error   bool     `json:"error"`
	Next    []string `json:"next"`
}
type CompiledNode struct {
	ID    string                        `json:"id"`
	fnPtr func(StateSchema) StateSchema `json:"-"`
	Spec  Spec                          `json:"spec"` // The spec of the node
}

type CompiledGraph struct {
	ID            string                   `json:"id"`
	CompiledNodes map[string]*CompiledNode `json:"compied_nodes"`
	EntryNode     string                   `json:"entry_node"` // The entry node of the graph
	Dispatcher    *Dispatcher              `json:"dispatcher"` // The dispatcher for the graph
}

func (gb *GraphBuilder) Compile() *CompiledGraph {
	compiled := make(map[string]*CompiledNode)
	entryNode := ""

	for i, node := range gb.Nodes {
		trigger := false
		if i == 0 {
			entryNode = node.ID
			trigger = true
		}
		compiled[node.ID] = &CompiledNode{
			ID:    node.ID,
			fnPtr: node.fnPtr,
			Spec: Spec{
				Trigger: trigger,
				Done:    false,
				Error:   false,
			},
		}
	}
	for _, edge := range gb.Edges {
		from := edge.SourceID

		compiled[from].Spec.Next = append(compiled[from].Spec.Next, edge.TargetID)
	}

	return &CompiledGraph{
		ID:            gb.ID,
		CompiledNodes: compiled,
		EntryNode:     entryNode,
	}
}

func (cp *CompiledGraph) invokeNode(
	nodeID string,
	data StateSchema,
	ch chan StateSchema,
) {
	thisNode := cp.CompiledNodes[nodeID]
	result := thisNode.fnPtr(data)
	thisNode.Spec.Done = true

	// TODO: Give dispatcher next nodes to trigger
	cp.Dispatcher.NextNode(thisNode.Spec.Next)

	ch <- result
}

func (cp *CompiledGraph) Run(data StateSchema) StateSchema {

	// for tick := range cp.Stream(data) {

	// }
	// return data
	return data
}

func (cp *CompiledGraph) Stream(
	data StateSchema,
	stream chan StateSchema,
	dp *Dispatcher,
) {
	cp.Dispatcher = dp
	internalCh := make(chan StateSchema)
	defer close(internalCh)
	defer close(stream)

	go cp.invokeNode(cp.EntryNode, data, internalCh)
	val := <-internalCh
	stream <- val

	for {
		tick := cp.Dispatcher.Tick()

		if len(tick) == 0 {
			break
		}
		for _, nodeID := range tick {
			go cp.invokeNode(nodeID, val, internalCh)
			val = <-internalCh
			stream <- val
		}

	}

}
