package gograph

import (
	"crypto/sha1"
	"fmt"
	"strconv"
	"time"
)

// Graph has many nodes
type Graph struct {
	DirectedNodes    []*DirectedNode
	RootDirectedNode *DirectedNode
}

// DirectedNode has a single parent node and a single child node
// TODO: Extend node to accept more than just Parents or Children. A programmer
//       should be able to specify an arbitrary number of edge relationships thats
//       can describe any abstract quality. Edges = [][]*DirectedNode, where the i-th edge
//       represents a certain qualitative relationship - e.g. like Edges[0] becomes
//       a Parent descriptor for the relationships between the nodes, and Edges[1]
//       becomes a Child descriptor for the relationships between the nodes. It's
//       really just a change in label that allows us to express more complex edges
//       in our graph than the mutual exclusivity the Parent/Child relationship suggests.
//       In truth, a mutual exclusivity would instead be enforced as a constraint on
//       updates to this more generic graph structure, that implementation must be
//       able to safely determine if nodes meet the constraints before insertion
//       into the graph.
type DirectedNode struct {
	Parents  []*DirectedNode
	Children []*DirectedNode
	ID       string
}

// CreateGraph returns a null graph object with a single root node. Does not create edges.
func CreateGraph() Graph {
	var node *DirectedNode
	var graphDirectedNodes []*DirectedNode
	return Graph{DirectedNodes: graphDirectedNodes, RootDirectedNode: node}
}

// CreateDirectedNode returns a node with a random ID. Does not create edges.
func CreateDirectedNode(graph Graph, parents []*DirectedNode, children []*DirectedNode) (Graph, *DirectedNode) {
	var nodeID = CreateDirectedNodeID()
	var node = &DirectedNode{
		Parents:  parents,
		Children: children,
		ID:       nodeID,
	}

	// TODO: Once a more generic notion of edges is implemented, DRY the next two blocks
	//       into a single loop.
	// Create edge between children nodes by updating their Parents reference to include this new node
	if len(children) > 0 {
		for _, child := range children {
			var index = FindDirectedNode(graph, child.ID)
			if index == 0 {
				panic("Attempted to become a parent of the root node! Forbidden node creation.")
			}
			graph.DirectedNodes[index].Parents = append(child.Parents, node)
		}
	}

	// Create edge between parent nodes by updating their Children reference to include this new node
	if len(parents) > 0 {
		for _, parent := range parents {
			var index = FindDirectedNode(graph, parent.ID)
			graph.DirectedNodes[index].Children = append(parent.Children, node)
		}
	}

	// Should a graph own the properties that make a node unique?
	// Or should the node itself own that same property. Perhaps both
	// for more extensible operations?
	graph.DirectedNodes = append(graph.DirectedNodes, node)
	if len(graph.DirectedNodes) == 1 {
		if len(graph.DirectedNodes[0].Parents) > 0 {
			panic("Root node specified non-existent parent node!")
		}
		graph.RootDirectedNode = node
	}

	return graph, node
}

// CreateDirectedNodeID generates a random SHA-1 hash to be used as a node ID
func CreateDirectedNodeID() string {
	s := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := sha1.New()
	h.Write([]byte(s))
	sha1Hash := h.Sum(nil)
	return fmt.Sprintf("%x", sha1Hash)
}

// FindDirectedNode traverses the array of nodes in the graph and returns the index of the node with the specified ID
func FindDirectedNode(graph Graph, ID string) int {
	for index, node := range graph.DirectedNodes {
		if node.ID == ID {
			return index
		}
	}
	return -1
}
