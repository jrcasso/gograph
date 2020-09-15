package gograph

import (
	"crypto/sha1"
	"fmt"
	"strconv"
	"time"
)

// Graph has many nodes
type Graph struct {
	Nodes []*Node
}

// Node is a generic recursive data structure that only has undirected edges.
type Node struct {
	Edges []*Node
	ID    string
}

// DirectedGraph has many nodes with directed edges
type DirectedGraph struct {
	DirectedNodes    []*DirectedNode
	RootDirectedNode *DirectedNode
}

// DirectedNode has a single parent node edges and a single child node edges
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
func CreateGraph() DirectedGraph {
	var node *DirectedNode
	var graphDirectedNodes []*DirectedNode
	return DirectedGraph{DirectedNodes: graphDirectedNodes, RootDirectedNode: node}
}

// CreateDirectedNode returns a node with a random ID. Does not create edges.
func CreateDirectedNode(graph DirectedGraph, parents []*DirectedNode, children []*DirectedNode) (DirectedGraph, *DirectedNode) {
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

// FindNode traverses the array of nodes in the graph and returns the index of the node with the specified ID
func FindNode(graph Graph, ID string) int {
	for index, node := range graph.Nodes {
		if node.ID == ID {
			return index
		}
	}
	return -1
}

// FindDirectedNode traverses the array of nodes in the graph and returns the index of the node with the specified ID
func FindDirectedNode(graph DirectedGraph, ID string) int {
	for index, node := range graph.DirectedNodes {
		if node.ID == ID {
			return index
		}
	}
	return -1
}

// CreateAdjecencyMatrix initial implementation, whether a directed edge exists from
// j-th element to the i-th element. We define the parent to child direction as the
// j-th to i-th direction.
func CreateAdjecencyMatrix(graph DirectedGraph) [][]int {
	var numNodes = len(graph.DirectedNodes)
	var adjMatrix = make([][]int, numNodes)
	for a := range graph.DirectedNodes {
		adjMatrix[a] = make([]int, numNodes)
	}

	for i, inode := range graph.DirectedNodes {
		for j, jnode := range graph.DirectedNodes {
			for _, jNodeParent := range jnode.Parents {
				if inode.ID == jNodeParent.ID {
					adjMatrix[i][j] = 1
				}
			}
		}
	}

	return adjMatrix
}

// CreateIncidenceMatrix initial implementation, whether a directed edge exists from
// j-th element to the i-th element. We define the parent to child direction as the
// j-th to i-th direction.
func CreateIncidenceMatrix(graph DirectedGraph) [][]int {
	var numNodes = len(graph.DirectedNodes)
	var adjMatrix = make([][]int, numNodes)
	for a := range graph.DirectedNodes {
		adjMatrix[a] = make([]int, numNodes)
	}

	for i, inode := range graph.DirectedNodes {
		for j, jnode := range graph.DirectedNodes {
			for _, jNodeParent := range jnode.Parents {
				if inode.ID == jNodeParent.ID {
					adjMatrix[i][j] = 1
				}
			}
			for _, jNodeChild := range jnode.Children {
				if inode.ID == jNodeChild.ID {
					adjMatrix[i][j] = -1
				}
			}
		}
	}

	return adjMatrix
}

// PrintMatrix prints a matrix of integers
func PrintMatrix(matrix [][]int) {
	for _, parent := range matrix {
		fmt.Printf("%+v\n", parent)
	}
}

// IsAntisymmetricMatrix checks that
func IsAntisymmetricMatrix(matrix [][]int) bool {
	var jsize = len(matrix)
	for j := 0; j < jsize; j++ {
		var isize = len(matrix[j])
		for i := j; i < isize; i++ {
			if len(matrix[i]) != len(matrix[j]) {
				panic("Provided matrix is not a square matrix!")
			}
			if matrix[j][i] != -matrix[i][j] {
				return false
			}
		}
	}
	return true
}
