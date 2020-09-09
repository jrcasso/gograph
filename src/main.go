package main

import (
	"fmt"

	"./gograph"
)

func main() {
	var parentDirectedNode, node *gograph.DirectedNode

	// graph = gograph.CreateGraph(nil)
	var linkedList = gograph.DirectedGraph{
		DirectedNodes:    nil,
		RootDirectedNode: nil,
	}

	linkedList, parentDirectedNode = gograph.CreateDirectedNode(linkedList, nil, nil)

	for i := 0; i < 10; i++ {
		linkedList, node = gograph.CreateDirectedNode(linkedList, []*gograph.DirectedNode{parentDirectedNode}, nil)
		parentDirectedNode = node
	}

	linkedList, _ = gograph.CreateDirectedNode(linkedList, []*gograph.DirectedNode{linkedList.DirectedNodes[0]}, nil)

	for _, node := range linkedList.DirectedNodes {
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

	var adjMatrix = gograph.CreateAdjecencyMatrix(linkedList)
	fmt.Println("Adjacency Matrix:")
	for _, arr := range adjMatrix {
		fmt.Printf("%+v\n", arr)
	}

	var incMatrix = gograph.CreateIncidenceMatrix(linkedList)
	fmt.Println("Incidence Matrix:")
	for _, arr := range incMatrix {
		fmt.Printf("%+v\n", arr)
	}
}
