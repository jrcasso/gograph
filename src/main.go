package main

import (
	"fmt"

	"./godag"
)

func main() {
	fmt.Println("main")

	var graph = godag.CreateGraph()
	var newNode = godag.CreateNode(godag.NullChild, godag.NullParent)
	graph = godag.AddNode(graph, newNode)
	// var graph = godag.CreateGraph()
	godag.PrintGraphNodes(graph)
	// root = godag.CreateRootNode()
}
