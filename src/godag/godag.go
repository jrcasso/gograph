package godag

import (
	"crypto/sha1"
	"fmt"
	"strconv"
	"time"
)

const (
	// NullChild contains the default string for null node ID values
	NullChild = ""

	// NullParent contains the default string for null node ID values
	NullParent = ""
)

// Node has a single parent node and a single child node
type Node struct {
	ParentIDs []string
	ChildIDs  []string
	ID        string
}

// Graph has many nodes
type Graph struct {
	Nodes    []Node
	RootNode Node
}

// CreateGraph returns a graph with a single root node. Does not create edges.
func CreateGraph() Graph {
	var rootNode = CreateNode("", "")
	var graphNodes []Node
	var graph = Graph{Nodes: graphNodes, RootNode: rootNode}
	return AddNode(graph, rootNode)
}

// CreateNode returns a node with a random ID. Does not create edges.
func CreateNode(parentID string, childID string) Node {
	s := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := sha1.New()
	h.Write([]byte(s))
	sha1Hash := h.Sum(nil)
	return Node{
		ParentIDs: append([]string{}, parentID),
		ChildIDs:  append([]string{}, childID),
		ID:        fmt.Sprintf("%x", sha1Hash),
	}
}

// AddNode add a node to the specified graph. Does not create edges.
func AddNode(graph Graph, node Node) Graph {
	graph.Nodes = append(graph.Nodes, node)
	// if node.ParentID != nil {
	// 	parentNode = FindNode(graph, node.ParentID)
	// }
	return graph
}

func FindNode(graph Graph, node Node) {

}

// AddRootNode add a node to the specified graph. Does not create edges.
func AddRootNode(graph Graph, node Node) Graph {
	graph.Nodes = append(graph.Nodes, node)
	graph.RootNode = node
	return graph
}

// PrintGraphNodes prints the nodes belonging to the specified graph.
func PrintGraphNodes(graph Graph) {

	for _, s := range graph.Nodes {
		// nodeID, err := hex.DecodeString(s.ID)
		nodeID := s.ID

		fmt.Printf("%x\n", nodeID)
	}
}

// PrintGraph prints the nodes belonging to the specified graph.
func PrintGraph(graph Graph) {
	fmt.Printf("%v", graph.Nodes)
}
