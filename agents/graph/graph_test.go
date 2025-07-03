package graph

import (
	"fmt"
	"testing"
)

type BaseState struct {
	A int
	B int
}

func Node1Fn(s StateSchema) StateSchema {
	state, ok := s.(*BaseState)
	if !ok {
		panic("Node1Fn: not BaseState")
	}
	state.A += 1
	return state
}
func Node2Fn(s StateSchema) StateSchema {
	state, ok := s.(*BaseState)
	if !ok {
		panic("Node2Fn: not BaseState")
	}
	state.B += 1
	return state
}

func TestBuild(t *testing.T) {

	GraphBuilder := NewGraphBuilder(
		"test-graph",
		WithNode("first", Node1Fn),
		WithNode("second", Node2Fn),
		WithEdge("first", "second"),
	)
	fmt.Println(GraphBuilder.ID)
	compiledGraph := GraphBuilder.Compile(&BaseState{A: 0, B: 0})

	tick := make(chan StateSchema)
	go compiledGraph.Stream(&BaseState{A: 0, B: 0}, tick)
	for data := range tick {
		fmt.Println("Tick data:", data)
	}
}
