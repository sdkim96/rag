package graph

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

type BaseState struct {
	A int
	B int
}

func Node1Fn(s StateSchema) StateSchema {
	fmt.Println("Node1Fn called")
	state, ok := s.(*BaseState)
	if !ok {
		panic("Node1Fn: not BaseState")
	}
	state.A += 1
	return state
}
func Node2Fn(s StateSchema) StateSchema {
	fmt.Println("Node2Fn called")
	state, ok := s.(*BaseState)
	if !ok {
		panic("Node2Fn: not BaseState")
	}
	state.B += 1
	return state
}
func Node3Fn(s StateSchema) StateSchema {
	fmt.Println("Node3Fn called")
	state, ok := s.(*BaseState)
	if !ok {
		panic("Node3Fn: not BaseState")
	}
	state.B += 1
	return state
}

func Node4Fn(s StateSchema) StateSchema {
	fmt.Println("Node4Fn called")
	state, ok := s.(*BaseState)
	if !ok {
		panic("Node4Fn: not BaseState")
	}
	state.A += 1
	return state
}

func TestBuild(t *testing.T) {

	ctx := context.Background()

	GraphBuilder := NewGraphBuilder(
		"test-graph",
		WithNode("first", Node1Fn),
		WithNode("second", Node2Fn),
		WithNode("third", Node3Fn),
		WithEdge("first", "second"),
		WithEdge("first", "third"),
		WithNode("second-first", Node4Fn),
		WithEdge("second", "second-first"),
	)
	compiledGraph := GraphBuilder.Compile()

	state := &BaseState{A: 0, B: 0}
	stream := make(chan StateSchema)

	go compiledGraph.Stream(state, stream, &Dispatcher{
		nextNodes: []string{},
		mu:        &sync.Mutex{},
	})
	for data := range stream {
		fmt.Println("Tick data:", data)
	}
}
