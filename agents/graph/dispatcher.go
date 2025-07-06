package graph

import "sync"

type Dispatcher struct {
	nextNodes []string
	mu        *sync.Mutex
}

func (dp *Dispatcher) NextNode(NodeID []string) {
	dp.nextNodes = append(dp.nextNodes, NodeID...)
}

func (dp *Dispatcher) Tick() []string {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	if len(dp.nextNodes) == 0 {
		return nil
	}

	tick := dp.nextNodes
	dp.nextNodes = []string{} // Reset next nodes for the next tick
	return tick
}
