package main

import (
	"fmt"

	"./gograph"
)

func main() {
	var parentDirectedNode, node *gograph.DirectedNode

	// graph = gograph.CreateGraph(nil)
	var graph = gograph.Graph{
		DirectedNodes:    nil,
		RootDirectedNode: nil,
	}
	graph, parentDirectedNode = gograph.CreateDirectedNode(graph, nil, nil)

	for i := 0; i < 10; i++ {
		graph, node = gograph.CreateDirectedNode(graph, []*gograph.DirectedNode{parentDirectedNode}, nil)
		parentDirectedNode = node
	}

	graph, _ = gograph.CreateDirectedNode(graph, []*gograph.DirectedNode{graph.DirectedNodes[0]}, nil)

	for _, node := range graph.DirectedNodes {
		if len(node.Children) > 0 && len(node.Parents) == 0 {
			fmt.Printf("The child of node %s is %+v\n", node.ID, node.Children[0].ID)
		}
		if len(node.Parents) > 0 && len(node.Children) > 0 {
			fmt.Printf("The child of node %s is %+v and the parent is %+v\n", node.ID, node.Children[0].ID, node.Parents[0].ID)
		}
		if len(node.Parents) > 0 && len(node.Children) == 0 {
			fmt.Printf("The parent of node %s is %+v\n", node.ID, node.Parents[0].ID)
		}
	}
}
