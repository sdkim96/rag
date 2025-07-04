package graph
 
type Spec struct {
    Trigger bool `json:"trigger"` 
    Done    bool `json:"done"`    
    Error   bool `json:"error"`   
	Next []string `json:"next"` // Next nodes to trigger
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
}
 
func (gb *GraphBuilder) Compile(data StateSchema) *CompiledGraph {
    nodesHashmap := make(map[string]*CompiledNode)
    entryNode := ""
 
    for i, node := range gb.Nodes {
        trigger := false
        if i == 0 {
            entryNode = node.ID
            trigger = true
        }
        nodesHashmap[node.ID] = &CompiledNode{
            ID:    node.ID,
            fnPtr: node.fnPtr,
            Spec: Spec{
                Trigger: trigger,
                Done:    false,
                Error:   false,
            },
        }
    }
	for j, edge := range gb.Edges {
		from := edge.SourceID
		to := edge.TargetID

		nodesHashmap[from].Spec.Next = append(nodesHashmap[from].Spec.Next, edge.TargetID)
	}
 
    return &CompiledGraph{
        ID:            gb.ID,
        CompiledNodes: nodesHashmap,
        EntryNode:     entryNode,
    }
}
 
func InvokeNode(
    node *CompiledNode,
    data StateSchema,
    ch chan StateSchema,
	dp *Dispatcher, // Assuming you have a Dispatcher to handle next nodes
) {
    result := node.fnPtr(data)
    node.Spec.Done = true


	// TODO: Give dispatcher next nodes to trigger
	dp.Recieve(node.Spec.Next)

	ch <- result
}
 
func (cp *CompiledGraph) Run(data StateSchema) StateSchema {
 
    // for tick := range cp.Stream(data) {
 
    // }
    // return data
    return data
}
 
func (cp *CompiledGraph) Stream(data StateSchema, stream chan StateSchema) {
 
    ch := make(chan StateSchema)
    defer close(ch)
 
    node := cp.CompiledNodes[cp.EntryNode]
    go InvokeNode(node, data, ch, cp.Dispatcher)
    val := <-ch
    stream <- val
 
    for {
        tick := cp.Dispatcher.Tick()

		if len(tick) == 0 {
			break
		}
        for _, node := range tick {
            if node.Spec.Trigger == true {
                go InvokeNode(node, val, ch cp.Dispatcher)
                val := <-ch
                stream <- val
            }
           
        }
 
    }
 
}