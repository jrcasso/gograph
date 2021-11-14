package main

import (
	"fmt"

	"github.com/jrcasso/gograph"
)

type MyValue struct {
	value map[string]string
}

func main() {
	var graphA = gograph.Graph{}
	nodeA := graphA.AddNode(
		[]*gograph.Edge{},
		MyValue{
			value: map[string]string{"name": "A"},
		},
	)
	nodeB := graphA.AddNode(
		[]*gograph.Edge{},
		MyValue{
			value: map[string]string{"name": "B"},
		},
	)
	graphA.AddEdge(
		MyValue{
			value: map[string]string{"value": "1"},
		},
		nodeA,
		nodeB,
	)

	fmt.Printf("%+v", graphA)

	var graphB = gograph.DirectedGraph{}
	nodeA = graphB.AddNode(
		nil,
		nil,
		[]*gograph.Node{},
		[]*gograph.Node{},
	)
	nodeB = graphB.AddNode(
		nil,
		nil,
		[]*gograph.Node{nodeA},
		[]*gograph.Node{},
	)

	fmt.Printf("%+v", graphB)
}

func (v MyValue) GetAsMap() map[string]string {
	return v.value
}
