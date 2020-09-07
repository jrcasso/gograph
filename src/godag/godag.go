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

// MergeNode has many parent nodes and a single child node
type MergeNode struct {
	ParentIDs []string
	Child     string
	ID        string
}

// Node has a single parent node and a single child node
type Node struct {
	ParentID string
	ChildID  string
	ID       string
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
		ParentID: parentID,
		ChildID:  childID,
		ID:       string(sha1Hash),
	}
}

// AddNode add a node to the specified graph. Does not create edges.
func AddNode(graph Graph, node Node) Graph {
	graph.Nodes = append(graph.Nodes, node)
	return graph
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
