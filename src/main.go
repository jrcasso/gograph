package main

import (
	"fmt"

	"./godag"
)

func main() {
	var parentNode, node, newNode *godag.Node

	// graph = godag.CreateGraph(nil)
	var graph = godag.Graph{
		Nodes:    nil,
		RootNode: nil,
	}
	graph, parentNode = godag.CreateNode(graph, nil, nil)

	for i := 0; i < 10; i++ {
		graph, node = godag.CreateNode(graph, []*godag.Node{parentNode}, nil)
		parentNode = node
	}

	graph, newNode = godag.CreateNode(graph, []*godag.Node{graph.Nodes[0]}, nil)
	if newNode == nil {
		panic("wtf")
	}

	// graph, newNode = godag.CreateNode(graph, nil, []*godag.Node{graph.Nodes[0]})
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
