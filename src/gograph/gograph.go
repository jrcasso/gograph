package gograph

import (
	"crypto/sha1"
	"fmt"
	"strconv"
	"time"
)

// Node has a single parent node and a single child node
// TODO: Extend node to accept more than just Parents or Children. A programmer
//       should be able to specify an arbitrary number of edge relationships thats
//       can describe any abstract quality. Edges = [][]*Node, where the i-th edge
//       represents a certain qualitative relationship - e.g. like Edges[0] becomes
//       a Parent descriptor for the relationships between the nodes, and Edges[1]
//       becomes a Child descriptor for the relationships between the nodes. It's
//       really just a change in label that allows us to express more complex edges
//       in our graph than the mutual exclusivity the Parent/Child relationship suggests.
//       In truth, a mutual exclusivity would instead be enforced as a constraint on
//       updates to this more generic graph structure, that implementation must be
//       able to safely determine if nodes meet the constraints before insertion
//       into the graph.
type Node struct {
	Parents  []*Node
	Children []*Node
	ID       string
}

// Graph has many nodes
type Graph struct {
	Nodes    []*Node
	RootNode *Node
}

// CreateGraph returns a null graph object with a single root node. Does not create edges.
func CreateGraph() Graph {
	var node *Node
	var graphNodes []*Node
	return Graph{Nodes: graphNodes, RootNode: node}
}

// CreateNode returns a node with a random ID. Does not create edges.
func CreateNode(graph Graph, parents []*Node, children []*Node) (Graph, *Node) {
	var nodeID = CreateNodeID()
	var node = &Node{
		Parents:  parents,
		Children: children,
		ID:       nodeID,
	}

	// TODO: Once a more generic notion of edges is implemented, DRY the next two blocks
	//       into a single loop.
	// Create edge between children nodes by updating their Parents reference to include this new node
	if len(children) > 0 {
		for _, child := range children {
			var index = FindNode(graph, child.ID)
			if index == 0 {
				panic("Attempted to become a parent of the root node! Forbidden node creation.")
			}
			graph.Nodes[index].Parents = append(child.Parents, node)
		}
	}

	// Create edge between parent nodes by updating their Children reference to include this new node
	if len(parents) > 0 {
		for _, parent := range parents {
			var index = FindNode(graph, parent.ID)
			graph.Nodes[index].Children = append(parent.Children, node)
		}
	}

	// Should a graph own the properties that make a node unique?
	// Or should the node itself own that same property. Perhaps both
	// for more extensible operations?
	graph.Nodes = append(graph.Nodes, node)
	if len(graph.Nodes) == 1 {
		if len(graph.Nodes[0].Parents) > 0 {
			panic("Root node specified non-existent parent node!")
		}
		graph.RootNode = node
	}

	return graph, node
}

// CreateNodeID generates a random SHA-1 hash to be used as a node ID
func CreateNodeID() string {
	s := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := sha1.New()
	h.Write([]byte(s))
	sha1Hash := h.Sum(nil)
	return fmt.Sprintf("%x", sha1Hash)
}

// FindNode traverses the array of nodes in the graph and returns the index of the node with the specified ID
func FindNode(graph Graph, ID string) int {
	for index, node := range graph.Nodes {
		if node.ID == ID {
			return index
		}
	}
	return -1
}
