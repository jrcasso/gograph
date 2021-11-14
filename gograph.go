package gograph

import (
	"github.com/google/uuid"
)

// Graph has many nodes
type Graph struct {
	Nodes []*Node
	Edges []*Edge
}

// Node is a generic recursive data structure that only has undirected edges
type Node struct {
	Edges []*Edge
	Value Value
	ID    string
}

// A generic edge can have any kind of value; an integer can be used to represent
// arbitrary states. For example, a directed edge could conceivably use two mutually
// exclusive values like 0 and 1 to represent arrow direction
// Conventionally, an edge has two vertices, but allowing an arbitrary number gets
// yields a hypergraph implementation, affording greater mathematical generality
type Edge struct {
	Value     Value
	Nodes     []*Node
	Direction int
}

// // A DirectedEdge is an edge with a binary direction associated with it
// type DirectedEdge struct {
// 	Edge      Edge
// 	Direction bool
// }

// // A DirectedEdge is an edge with a binary direction associated with it
// type DirectedNode struct {
// 	DirectedEdges []*DirectedEdge
// 	Value         Value
// 	ID            string
// }

// A DirectedGraph is a special case of a graph where the edges between
// nodes are bidirectional
type DirectedGraph struct {
	Nodes []*Node
	Edges []*Edge
}

type Value interface {
	GetAsMap() map[string]string
}

func (g *Graph) AddNodeWithID(edges []*Edge, value Value, id string) *Node {
	var newNode = Node{
		Edges: edges,
		Value: value,
		ID:    id,
	}
	g.Nodes = append(g.Nodes, &newNode)
	return &newNode
}

func (g *Graph) AddNode(edges []*Edge, value Value) *Node {
	var id = uuid.NewString()
	return g.AddNodeWithID(edges, value, id)
}

func (g *Graph) AddEdge(value Value, nodeA *Node, nodeB *Node) *Edge {
	var edge = Edge{
		Value: value,
		Nodes: []*Node{
			nodeA,
			nodeB,
		},
	}
	nodeA.Edges = append(nodeA.Edges, &edge)
	nodeB.Edges = append(nodeB.Edges, &edge)
	g.Edges = append(g.Edges, &edge)
	return &edge
}

func (dg *DirectedGraph) AddNodeWithID(nodeValue Value, edgeValue Value, parents []*Node, children []*Node, id string) *Node {
	var edge Edge
	var newNode = Node{
		Edges: []*Edge{},
		Value: nodeValue,
		ID:    id,
	}

	// Positive direction is defined as the direction from a
	// parent node to a child node. Direction generalizes in
	// hypergraph implementations.
	for _, parent := range parents {
		edge = Edge{
			Value: edgeValue,
			Nodes: []*Node{
				parent,
				&newNode,
			},
			Direction: 1,
		}
		parent.Edges = append(parent.Edges, &edge)
		newNode.Edges = append(newNode.Edges, &edge)
		dg.Edges = append(dg.Edges, &edge)
	}

	for _, child := range children {
		edge = Edge{
			Value: edgeValue,
			Nodes: []*Node{
				&newNode,
				child,
			},
			Direction: 1,
		}
		newNode.Edges = append(newNode.Edges, &edge)
		child.Edges = append(child.Edges, &edge)
		dg.Edges = append(dg.Edges, &edge)
	}

	dg.Nodes = append(dg.Nodes, &newNode)
	return &newNode
}

func (dg *DirectedGraph) AddNode(nodeValue Value, edgeValue Value, parents []*Node, children []*Node) *Node {
	var id = uuid.NewString()
	return dg.AddNodeWithID(nodeValue, edgeValue, parents, children, id)
}

// GetChildrenNodes looks at a node's edges and returns those
// which have a negative direction with other nodes. The direction
//
func (node *Node) GetChildrenNodes() []*Node {
	var children = []*Node{}
	for _, edge := range node.Edges {
		if edge.Direction > 0 {
			children = append(children, edge.Nodes[1])
		}
		if edge.Direction < 0 {
			children = append(children, edge.Nodes[0])
		}
	}
	return children
}

// GetChildrenNodes looks at a node's edges and returns those
// which have a negative direction with other nodes. The direction
//
func (node *Node) GetParentNodes() []*Node {
	var children = []*Node{}
	for _, edge := range node.Edges {
		if edge.Direction > 0 {
			children = append(children, edge.Nodes[1])
		}
		if edge.Direction < 0 {
			children = append(children, edge.Nodes[0])
		}
	}
	return children
}

func TopologicalSort(graph DirectedGraph) []*Node {
	var sorted []*Node
	var childlessNodes []*Node
	for _, node := range graph.Nodes {
		if len(node.GetChildrenNodes()) == 0 {
			childlessNodes = append(childlessNodes, node)
		}
	}

	for len(childlessNodes) > 0 {
		var nextNode = childlessNodes[0]
		childlessNodes = childlessNodes[1:]
		sorted = append(sorted, nextNode)
		for _, parent := range nextNode.Parents {
			// Remove edge from parent to this child node
			for index, childNode := range parent.Children {
				if childNode.ID == nextNode.ID {
					parent.Children = append(parent.Children[:index], parent.Children[index+1:]...)
					break
				}
			}
			if len(parent.Children) == 0 {
				childlessNodes = append(childlessNodes, parent)
			}
		}
		// Now that the edges have been deleted from the parents to this particular child node,
		// we can more efficiently mass-delete the edges from the child to its parents.
		// One might ask, "Why don't we call the DeleteDirectedEdge function?", to which I would respond,
		// "The DeleteDirectedEdge function is an in-place function that will mutate the nextNode.Parents slice
		// and lead to a 6-hour bug; the dummy variable would iterate through a mutated nextNode.Parents above.".
		nextNode.Parents = []*Node{}
	}

	if len(sorted) != len(graph.Nodes) {
		panic("Provided graph is not a directed, acyclic graph, and does not have a valid topological ordering")
	}

	return sorted
}

// // CreateGraph returns a null graph object with a single root node. Does not create edges.
// func CreateGraph() DirectedGraph {
// 	var node *DirectedNode
// 	var graphDirectedNodes []*DirectedNode
// 	return DirectedGraph{DirectedNodes: graphDirectedNodes, RootDirectedNode: node}
// }

// // CreateDirectedEdge creates a parent-child relationship between two specified nodes.
// func CreateDirectedEdge(graph DirectedGraph, parent *DirectedNode, child *DirectedNode) (DirectedGraph, *DirectedNode, *DirectedNode) {
// 	child.Parents = append(child.Parents, parent)
// 	parent.Children = append(parent.Children, child)
// 	return graph, parent, child
// }

// // CreateDirectedNode returns a node with a random ID. Does not create edges.
// func CreateDirectedNode(graph DirectedGraph, values map[string]string, parents []*DirectedNode, children []*DirectedNode) (DirectedGraph, *DirectedNode) {
// 	var nodeID = CreateDirectedNodeID()
// 	var node = &DirectedNode{
// 		Parents:  parents,
// 		Children: children,
// 		Values:   values,
// 		ID:       nodeID,
// 	}

// 	// TODO: Once a more generic notion of edges is implemented, DRY the next two blocks
// 	//       into a single loop.
// 	// Create edge between children nodes by updating their Parents reference to include this new node
// 	if len(children) > 0 {
// 		for _, child := range children {
// 			var index, _ = FindDirectedNode(graph, child.ID)
// 			if index == 0 {
// 				panic("Attempted to become a parent of the root node! Forbidden node creation.")
// 			}
// 			graph.DirectedNodes[index].Parents = append(child.Parents, node)
// 		}
// 	}

// 	// Create edge between parent nodes by updating their Children reference to include this new node
// 	if len(parents) > 0 {
// 		for _, parent := range parents {
// 			var index, _ = FindDirectedNode(graph, parent.ID)
// 			graph.DirectedNodes[index].Children = append(parent.Children, node)
// 		}
// 	}

// 	// Should a graph own the properties that make a node unique?
// 	// Or should the node itself own that same property. Perhaps both
// 	// for more extensible operations?
// 	graph.DirectedNodes = append(graph.DirectedNodes, node)
// 	if len(graph.DirectedNodes) == 1 {
// 		if len(graph.DirectedNodes[0].Parents) > 0 {
// 			panic("Root node specified non-existent parent node!")
// 		}
// 		graph.RootDirectedNode = node
// 	}

// 	return graph, node
// }

// // CreateDirectedNodeID generates a random SHA-1 hash to be used as a node ID
// func CreateDirectedNodeID() string {
// 	s := strconv.FormatInt(time.Now().UnixNano(), 10)
// 	h := sha1.New()
// 	h.Write([]byte(s))
// 	sha1Hash := h.Sum(nil)
// 	return fmt.Sprintf("%x", sha1Hash)
// }

// // DeleteDirectedEdge is an in-place function that deletes the edge connection between a child and parent node
// // TODO: Something about the in-placeness, which is in contrast to the CreateDirectedEdge function.
// func DeleteDirectedEdge(graph DirectedGraph, parent *DirectedNode, child *DirectedNode) (DirectedGraph, *DirectedNode, *DirectedNode) {
// 	for index, childNode := range parent.Children {
// 		if childNode.ID == child.ID {
// 			parent.Children = append(parent.Children[:index], parent.Children[index+1:]...)
// 			break
// 		}
// 	}
// 	for index, parentNode := range child.Parents {
// 		if parentNode.ID == parent.ID {
// 			child.Parents = append(child.Parents[:index], child.Parents[index+1:]...)
// 			break
// 		}
// 	}

// 	return graph, parent, child
// }

// // FindNode traverses the array of nodes in the graph and returns the index of the node with the specified ID
// func FindNode(graph Graph, ID string) int {
// 	for index, node := range graph.Nodes {
// 		if node.ID == ID {
// 			return index
// 		}
// 	}
// 	return -1
// }

// // FindDirectedNode traverses the array of nodes in the graph and returns the index of the node with the specified ID
// func FindDirectedNode(graph DirectedGraph, ID string) (int, DirectedNode) {
// 	for index, node := range graph.DirectedNodes {
// 		if node.ID == ID {
// 			return index, *node
// 		}
// 	}
// 	return -1, DirectedNode{}
// }

// // FindNodesByValues traverses the array of nodes in the graph and returns the nodes that match the passed values
// func FindNodesByValues(graph DirectedGraph, values map[string]string) []*DirectedNode {
// 	var isMatch bool
// 	var results []*DirectedNode

// 	for _, node := range graph.DirectedNodes {
// 		isMatch = true
// 		for key, value := range values {
// 			if node.Values[key] != value {
// 				isMatch = false
// 			}
// 		}
// 		if isMatch {
// 			results = append(results, node)
// 		}
// 	}
// 	return results
// }

// // TopologicalSort implements Kahn's algorithm to sort a directed acyclic graph
// func TopologicalSort(graph DirectedGraph) []*DirectedNode {
// 	var sorted []*DirectedNode
// 	var childlessNodes []*DirectedNode
// 	for _, node := range graph.DirectedNodes {
// 		if len(node.Children) == 0 {
// 			childlessNodes = append(childlessNodes, node)
// 		}
// 	}

// 	for len(childlessNodes) > 0 {
// 		var nextNode = childlessNodes[0]
// 		childlessNodes = childlessNodes[1:]
// 		sorted = append(sorted, nextNode)
// 		for _, parent := range nextNode.Parents {
// 			// Remove edge from parent to this child node
// 			for index, childNode := range parent.Children {
// 				if childNode.ID == nextNode.ID {
// 					parent.Children = append(parent.Children[:index], parent.Children[index+1:]...)
// 					break
// 				}
// 			}
// 			if len(parent.Children) == 0 {
// 				childlessNodes = append(childlessNodes, parent)
// 			}
// 		}
// 		// Now that the edges have been deleted from the parents to this particular child node,
// 		// we can more efficiently mass-delete the edges from the child to its parents.
// 		// One might ask, "Why don't we call the DeleteDirectedEdge function?", to which I would respond,
// 		// "The DeleteDirectedEdge function is an in-place function that will mutate the nextNode.Parents slice
// 		// and lead to a 6-hour bug; the dummy variable would iterate through a mutated nextNode.Parents above.".
// 		nextNode.Parents = []*DirectedNode{}
// 	}

// 	if len(sorted) != len(graph.DirectedNodes) {
// 		panic("Provided graph is not a directed, acyclic graph, and does not have a valid topological ordering")
// 	}

// 	return sorted
// }

// // CreateAdjecencyMatrix initial implementation, whether a directed edge exists from
// // j-th element to the i-th element. We define the parent to child direction as the
// // j-th to i-th direction.
// func CreateAdjecencyMatrix(graph DirectedGraph) [][]int {
// 	var numNodes = len(graph.DirectedNodes)
// 	var adjMatrix = make([][]int, numNodes)
// 	for a := range graph.DirectedNodes {
// 		adjMatrix[a] = make([]int, numNodes)
// 	}

// 	for i, inode := range graph.DirectedNodes {
// 		for j, jnode := range graph.DirectedNodes {
// 			for _, jNodeParent := range jnode.Parents {
// 				if inode.ID == jNodeParent.ID {
// 					adjMatrix[i][j] = 1
// 				}
// 			}
// 		}
// 	}

// 	return adjMatrix
// }

// // CreateIncidenceMatrix initial implementation, whether a directed edge exists from
// // j-th element to the i-th element. We define the parent to child direction as the
// // j-th to i-th direction.
// func CreateIncidenceMatrix(graph DirectedGraph) [][]int {
// 	var numNodes = len(graph.DirectedNodes)
// 	var adjMatrix = make([][]int, numNodes)
// 	for a := range graph.DirectedNodes {
// 		adjMatrix[a] = make([]int, numNodes)
// 	}

// 	for i, inode := range graph.DirectedNodes {
// 		for j, jnode := range graph.DirectedNodes {
// 			for _, jNodeParent := range jnode.Parents {
// 				if inode.ID == jNodeParent.ID {
// 					adjMatrix[i][j] = 1
// 				}
// 			}
// 			for _, jNodeChild := range jnode.Children {
// 				if inode.ID == jNodeChild.ID {
// 					adjMatrix[i][j] = -1
// 				}
// 			}
// 		}
// 	}

// 	return adjMatrix
// }

// // PrintMatrix prints a matrix of integers
// func PrintMatrix(matrix [][]int) {
// 	for _, parent := range matrix {
// 		fmt.Printf("%+v\n", parent)
// 	}
// }

// // IsAntisymmetricMatrix checks that
// func IsAntisymmetricMatrix(matrix [][]int) bool {
// 	var jsize = len(matrix)
// 	for j := 0; j < jsize; j++ {
// 		var isize = len(matrix[j])
// 		for i := j; i < isize; i++ {
// 			if len(matrix[i]) != len(matrix[j]) {
// 				panic("Provided matrix is not a square matrix!")
// 			}
// 			if matrix[j][i] != -matrix[i][j] {
// 				return false
// 			}
// 		}
// 	}
// 	return true
// }
