package main

import (
	"fmt"

	"./gograph"
)

func main() {
	var parentNode, node, newNode *gograph.Node

	// graph = gograph.CreateGraph(nil)
	var graph = gograph.Graph{
		Nodes:    nil,
		RootNode: nil,
	}
	graph, parentNode = gograph.CreateNode(graph, nil, nil)

	for i := 0; i < 10; i++ {
		graph, node = gograph.CreateNode(graph, []*gograph.Node{parentNode}, nil)
		parentNode = node
	}

	graph, newNode = gograph.CreateNode(graph, []*gograph.Node{graph.Nodes[0]}, nil)
	if newNode == nil {
		panic("wtf")
	}

	// graph, newNode = gograph.CreateNode(graph, nil, []*gograph.Node{graph.Nodes[0]})
	// if newNode == nil {
	// 	panic("wtf")
	// }

	for _, node := range graph.Nodes {
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
