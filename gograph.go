package gograph

import (
	"crypto/sha1"
	"fmt"
	"strconv"
	"time"
)

// Graph has many nodes with undirected edges
type Graph struct {
	Nodes   []*Node
	NodeMap map[string]*Node
}

// Node is a generic data structure with undirected edges
type Node struct {
	Edges  []*Node
	Values map[string]string
	ID     string
}

// GenericGraph represents a graph with flexible edge relationships
type GenericGraph struct {
	Nodes   []*GenericNode
	NodeMap map[string]*GenericNode
}

// GenericNode represents a node with flexible edge relationships
type GenericNode struct {
	// Edges is a slice of slices of nodes, where each index represents a different relationship type
	// For example, Edges[0] might be "Parents", Edges[1] might be "Children", Edges[2] might be "Friends", etc.
	// We thus track relationship types in EdgeTypes, which maps relationship names to their index in the Edges slice
	Edges     [][]*GenericNode
	EdgeTypes map[string]int
	Values    map[string]string
	ID        string
}

// DirectedGraph has many nodes with directed edges
type DirectedGraph struct {
	DirectedNodes    []*DirectedNode
	RootDirectedNode *DirectedNode
	NodeMap          map[string]*DirectedNode
}

// DirectedNode represents a node with directed edges
type DirectedNode struct {
	Parents  []*DirectedNode
	Children []*DirectedNode
	Values   map[string]string
	ID       string
}

// CreateNodeID generates a random SHA-1 hash to be used as a node ID
func CreateNodeID() string {
	s := strconv.FormatInt(time.Now().UnixNano(), 10)
	h := sha1.New()
	h.Write([]byte(s))
	sha1Hash := h.Sum(nil)
	return fmt.Sprintf("%x", sha1Hash)
}

// PrintMatrix prints a matrix of integers
func PrintMatrix(matrix [][]int) {
	for _, row := range matrix {
		fmt.Printf("%+v\n", row)
	}
}

// CreateGraph returns a new undirected graph
func CreateGraph() Graph {
	return Graph{
		Nodes:   []*Node{},
		NodeMap: make(map[string]*Node),
	}
}

// CreateNode creates a new node in the undirected graph
func CreateNode(graph Graph, values map[string]string) (Graph, *Node) {
	nodeID := CreateNodeID()
	node := &Node{
		Edges:  []*Node{},
		Values: values,
		ID:     nodeID,
	}

	graph.Nodes = append(graph.Nodes, node)
	graph.NodeMap[node.ID] = node

	return graph, node
}

// AddUndirectedEdge creates an undirected edge between two nodes
func AddUndirectedEdge(graph Graph, node1 *Node, node2 *Node) Graph {
	node1.Edges = append(node1.Edges, node2)
	node2.Edges = append(node2.Edges, node1)
	return graph
}

// RemoveUndirectedEdge removes an undirected edge between two nodes
func RemoveUndirectedEdge(graph Graph, node1 *Node, node2 *Node) Graph {
	for i, node := range node1.Edges {
		if node.ID == node2.ID {
			node1.Edges = append(node1.Edges[:i], node1.Edges[i+1:]...)
			break
		}
	}

	for i, node := range node2.Edges {
		if node.ID == node1.ID {
			node2.Edges = append(node2.Edges[:i], node2.Edges[i+1:]...)
			break
		}
	}

	return graph
}

// FindNode finds a node by ID in an undirected graph
func FindNode(graph Graph, ID string) (int, *Node) {
	// First check the map for O(1) lookup
	if node, ok := graph.NodeMap[ID]; ok {
		for index, n := range graph.Nodes {
			if n.ID == ID {
				return index, node
			}
		}
	}

	// Fallback to linear search
	for index, node := range graph.Nodes {
		if node.ID == ID {
			return index, node
		}
	}
	return -1, nil
}

// CreateDirectedEdge creates a parent-child relationship between two specified nodes.
func CreateDirectedEdge(graph DirectedGraph, parent *DirectedNode, child *DirectedNode) (DirectedGraph, *DirectedNode, *DirectedNode) {
	child.Parents = append(child.Parents, parent)
	parent.Children = append(parent.Children, child)
	return graph, parent, child
}

// CreateDirectedNode returns a node with a random ID. Does not create edges.
func CreateDirectedNode(graph DirectedGraph, values map[string]string, parents []*DirectedNode, children []*DirectedNode) (DirectedGraph, *DirectedNode) {
	if graph.NodeMap == nil {
		graph.NodeMap = make(map[string]*DirectedNode)
	}

	var nodeID = CreateDirectedNodeID()
	var node = &DirectedNode{
		Parents:  parents,
		Children: children,
		Values:   values,
		ID:       nodeID,
	}

	updateRelationships := func(relatedNodes []*DirectedNode, isParentRelation bool) {
		for _, relatedNode := range relatedNodes {
			var index, _ = FindDirectedNode(graph, relatedNode.ID)

			if !isParentRelation && index == 0 {
				panic("Attempted to become a parent of the root node! Forbidden node creation.")
			}

			if isParentRelation {
				graph.DirectedNodes[index].Children = append(graph.DirectedNodes[index].Children, node)
			} else {
				graph.DirectedNodes[index].Parents = append(graph.DirectedNodes[index].Parents, node)
			}
		}
	}

	if len(children) > 0 {
		updateRelationships(children, false)
	}

	if len(parents) > 0 {
		updateRelationships(parents, true)
	}

	graph.DirectedNodes = append(graph.DirectedNodes, node)
	if len(graph.DirectedNodes) == 1 {
		if len(graph.DirectedNodes[0].Parents) > 0 {
			panic("Root node specified non-existent parent node!")
		}
		graph.RootDirectedNode = node
	}

	graph.NodeMap[node.ID] = node
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

// DeleteDirectedEdge is an in-place function that deletes the edge connection between a child and parent node
func DeleteDirectedEdge(graph DirectedGraph, parent *DirectedNode, child *DirectedNode) (DirectedGraph, *DirectedNode, *DirectedNode) {
	for index, childNode := range parent.Children {
		if childNode.ID == child.ID {
			parent.Children = append(parent.Children[:index], parent.Children[index+1:]...)
			break
		}
	}
	for index, parentNode := range child.Parents {
		if parentNode.ID == parent.ID {
			child.Parents = append(child.Parents[:index], child.Parents[index+1:]...)
			break
		}
	}

	return graph, parent, child
}

// FindDirectedNode finds a node by ID with O(1) complexity using the NodeMap
func FindDirectedNode(graph DirectedGraph, ID string) (int, DirectedNode) {
	// First check the map for O(1) lookup if it exists
	if graph.NodeMap != nil {
		if node, ok := graph.NodeMap[ID]; ok {
			for index, n := range graph.DirectedNodes {
				if n.ID == ID {
					return index, *node
				}
			}
		}
	}

	// Fallback to linear search (as done generically in FindNode)
	for index, node := range graph.DirectedNodes {
		if node.ID == ID {
			return index, *node
		}
	}
	return -1, DirectedNode{}
}

// FindNodesByValues traverses the array of nodes in the graph and returns the nodes that match the passed values
func FindNodesByValues(graph DirectedGraph, values map[string]string) []*DirectedNode {
	var isMatch bool
	var results []*DirectedNode

	for _, node := range graph.DirectedNodes {
		isMatch = true
		for key, value := range values {
			if node.Values[key] != value {
				isMatch = false
			}
		}
		if isMatch {
			results = append(results, node)
		}
	}
	return results
}

// TopologicalSort implements Kahn's algorithm to sort a directed acyclic graph
// Returns an error if the graph contains a cycle
func TopologicalSort(graph DirectedGraph) ([]*DirectedNode, error) {
	var sorted []*DirectedNode
	var childlessNodes []*DirectedNode
	for _, node := range graph.DirectedNodes {
		if len(node.Children) == 0 {
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
		nextNode.Parents = []*DirectedNode{}
	}

	if len(sorted) != len(graph.DirectedNodes) {
		return nil, fmt.Errorf("provided graph is not a directed, acyclic graph, and does not have a valid topological ordering")
	}

	return sorted, nil
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

// IsAntisymmetricMatrix checks if a matrix is antisymmetric
func IsAntisymmetricMatrix(matrix [][]int) bool {
	if len(matrix) == 0 {
		return false
	}

	// Check if matrix is square
	n := len(matrix)
	for _, row := range matrix {
		if len(row) != n {
			return false
		}
	}

	for i := 0; i < n; i++ {
		if matrix[i][i] != 0 {
			return false
		}
		for j := 0; j < n; j++ {
			if matrix[i][j] != -matrix[j][i] {
				return false
			}
		}
	}
	return true
}

// HasCycle checks if the directed graph contains any cycles
func HasCycle(graph DirectedGraph) bool {
	visited := make(map[string]bool)
	recursionStack := make(map[string]bool)

	for _, node := range graph.DirectedNodes {
		if !visited[node.ID] {
			if hasCycleUtil(node, visited, recursionStack) {
				return true
			}
		}
	}
	return false
}

// FindPath finds a path between source and target nodes with a BFS, if one exists
func FindPath(graph DirectedGraph, sourceID, targetID string) ([]*DirectedNode, error) {
	sourceIndex, _ := FindDirectedNode(graph, sourceID)
	targetIndex, _ := FindDirectedNode(graph, targetID)

	if sourceIndex == -1 || targetIndex == -1 {
		return nil, fmt.Errorf("source or target node not found in graph")
	}

	visited := make(map[string]bool)
	queue := [][]*DirectedNode{{graph.DirectedNodes[sourceIndex]}}

	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		node := path[len(path)-1]

		if node.ID == targetID {
			return path, nil
		}

		if visited[node.ID] {
			continue
		}

		visited[node.ID] = true

		for _, child := range node.Children {
			if !visited[child.ID] {
				newPath := make([]*DirectedNode, len(path)+1)
				copy(newPath, path)
				newPath[len(path)] = child
				queue = append(queue, newPath)
			}
		}
	}

	return nil, fmt.Errorf("no path found from source to target")
}

// FindStronglyConnectedComponents finds all strongly connected components in the graph (Tarjan's algorithm)
func FindStronglyConnectedComponents(graph DirectedGraph) [][]*DirectedNode {
	if len(graph.DirectedNodes) == 0 {
		return [][]*DirectedNode{}
	}

	var components [][]*DirectedNode
	var index int = 0
	var stack []*DirectedNode

	indices := make(map[string]int)
	lowLinks := make(map[string]int)
	onStack := make(map[string]bool)

	for _, node := range graph.DirectedNodes {
		indices[node.ID] = -1
		lowLinks[node.ID] = -1
		onStack[node.ID] = false
	}

	var strongConnect func(*DirectedNode)
	strongConnect = func(node *DirectedNode) {
		indices[node.ID] = index
		lowLinks[node.ID] = index
		index++
		stack = append(stack, node)
		onStack[node.ID] = true

		for _, child := range node.Children {
			if indices[child.ID] == -1 {
				strongConnect(child)
				lowLinks[node.ID] = min(lowLinks[node.ID], lowLinks[child.ID])
			} else if onStack[child.ID] {
				lowLinks[node.ID] = min(lowLinks[node.ID], indices[child.ID])
			}
		}

		if lowLinks[node.ID] == indices[node.ID] {
			var component []*DirectedNode
			var componentNode *DirectedNode
			for {
				componentNode, stack = stack[len(stack)-1], stack[:len(stack)-1]
				onStack[componentNode.ID] = false
				component = append(component, componentNode)
				if componentNode.ID == node.ID {
					break
				}
			}
			components = append(components, component)
		}
	}

	for _, node := range graph.DirectedNodes {
		if indices[node.ID] == -1 {
			strongConnect(node)
		}
	}

	return components
}

// CreateGenericGraph returns a new generic graph
func CreateGenericGraph(edgeTypes []string) GenericGraph {
	edgeTypeMap := make(map[string]int)
	for i, edgeType := range edgeTypes {
		edgeTypeMap[edgeType] = i
	}

	return GenericGraph{
		Nodes:   []*GenericNode{},
		NodeMap: make(map[string]*GenericNode),
	}
}

// CreateGenericNode creates a new node with the specified edge types
func CreateGenericNode(graph GenericGraph, values map[string]string) (GenericGraph, *GenericNode) {
	nodeID := CreateDirectedNodeID() // Reuse the existing ID generation function
	edges := make([][]*GenericNode, len(graph.Nodes[0].EdgeTypes))

	node := &GenericNode{
		Edges:     edges,
		EdgeTypes: graph.Nodes[0].EdgeTypes,
		Values:    values,
		ID:        nodeID,
	}

	graph.Nodes = append(graph.Nodes, node)
	graph.NodeMap[node.ID] = node

	return graph, node
}

// AddEdge adds a directional edge of the specified type between two nodes
func AddEdge(graph GenericGraph, fromNode *GenericNode, toNode *GenericNode, edgeType string) (GenericGraph, error) {
	edgeIndex, exists := fromNode.EdgeTypes[edgeType]
	if !exists {
		return graph, fmt.Errorf("edge type '%s' not defined in graph", edgeType)
	}
	fromNode.Edges[edgeIndex] = append(fromNode.Edges[edgeIndex], toNode)
	return graph, nil
}

// RemoveEdge removes a directional edge of the specified type between two nodes
func RemoveEdge(graph GenericGraph, fromNode *GenericNode, toNode *GenericNode, edgeType string) (GenericGraph, error) {
	edgeIndex, exists := fromNode.EdgeTypes[edgeType]
	if !exists {
		return graph, fmt.Errorf("edge type '%s' not defined in graph", edgeType)
	}

	for i, node := range fromNode.Edges[edgeIndex] {
		if node.ID == toNode.ID {
			fromNode.Edges[edgeIndex] = append(fromNode.Edges[edgeIndex][:i], fromNode.Edges[edgeIndex][i+1:]...)
			break
		}
	}

	return graph, nil
}

// GetRelatedNodes returns all nodes related to the given node by the specified edge type
func GetRelatedNodes(node *GenericNode, edgeType string) ([]*GenericNode, error) {
	edgeIndex, exists := node.EdgeTypes[edgeType]
	if !exists {
		return nil, fmt.Errorf("edge type '%s' not defined in graph", edgeType)
	}

	return node.Edges[edgeIndex], nil
}

// ConvertDirectedToGeneric converts a DirectedGraph to a GenericGraph
func ConvertDirectedToGeneric(directedGraph DirectedGraph) GenericGraph {
	edgeTypes := []string{"Parents", "Children"}
	genericGraph := CreateGenericGraph(edgeTypes)

	edgeTypeMap := make(map[string]int)
	for i, edgeType := range edgeTypes {
		edgeTypeMap[edgeType] = i
	}

	nodeMap := make(map[string]*GenericNode)
	for _, dirNode := range directedGraph.DirectedNodes {
		nodeID := dirNode.ID
		edges := make([][]*GenericNode, len(edgeTypes))

		genericNode := &GenericNode{
			Edges:     edges,
			EdgeTypes: edgeTypeMap,
			Values:    dirNode.Values,
			ID:        nodeID,
		}

		genericGraph.Nodes = append(genericGraph.Nodes, genericNode)
		genericGraph.NodeMap[nodeID] = genericNode
		nodeMap[nodeID] = genericNode
	}

	for _, dirNode := range directedGraph.DirectedNodes {
		genericNode := nodeMap[dirNode.ID]

		// Add parent edges
		parentsIndex := genericNode.EdgeTypes["Parents"]
		for _, parent := range dirNode.Parents {
			genericParent := nodeMap[parent.ID]
			genericNode.Edges[parentsIndex] = append(genericNode.Edges[parentsIndex], genericParent)
		}

		// Add child edges
		childrenIndex := genericNode.EdgeTypes["Children"]
		for _, child := range dirNode.Children {
			genericChild := nodeMap[child.ID]
			genericNode.Edges[childrenIndex] = append(genericNode.Edges[childrenIndex], genericChild)
		}
	}

	return genericGraph
}

// ConvertGenericToDirected converts a GenericGraph to a DirectedGraph
// Note: This assumes the GenericGraph has "Parents" and "Children" edge types
func ConvertGenericToDirected(genericGraph GenericGraph) (DirectedGraph, error) {
	directedGraph := DirectedGraph{
		DirectedNodes: []*DirectedNode{},
		NodeMap:       make(map[string]*DirectedNode),
	}

	// Check if the generic graph has the required edge types
	if len(genericGraph.Nodes) > 0 {
		_, hasParents := genericGraph.Nodes[0].EdgeTypes["Parents"]
		_, hasChildren := genericGraph.Nodes[0].EdgeTypes["Children"]
		if !hasParents || !hasChildren {
			return directedGraph, fmt.Errorf("generic graph does not have required 'Parents' and 'Children' edge types")
		}
	}

	nodeMap := make(map[string]*DirectedNode)
	for _, genNode := range genericGraph.Nodes {
		dirNode := &DirectedNode{
			Parents:  []*DirectedNode{},
			Children: []*DirectedNode{},
			Values:   genNode.Values,
			ID:       genNode.ID,
		}

		directedGraph.DirectedNodes = append(directedGraph.DirectedNodes, dirNode)
		directedGraph.NodeMap[dirNode.ID] = dirNode
		nodeMap[genNode.ID] = dirNode
	}

	for _, genNode := range genericGraph.Nodes {
		dirNode := nodeMap[genNode.ID]

		// Add parent edges...
		parentsIndex := genNode.EdgeTypes["Parents"]
		for _, parent := range genNode.Edges[parentsIndex] {
			dirParent := nodeMap[parent.ID]
			dirNode.Parents = append(dirNode.Parents, dirParent)
			dirParent.Children = append(dirParent.Children, dirNode)
		}
	}

	if len(directedGraph.DirectedNodes) > 0 {
		for _, node := range directedGraph.DirectedNodes {
			if len(node.Parents) == 0 {
				directedGraph.RootDirectedNode = node
				break
			}
		}

		// If no node without parents is found, just use the first node
		if directedGraph.RootDirectedNode == nil {
			directedGraph.RootDirectedNode = directedGraph.DirectedNodes[0]
		}
	}

	return directedGraph, nil
}

// hasCycleUtil is a helper function for HasCycle using DFS
func hasCycleUtil(node *DirectedNode, visited, recursionStack map[string]bool) bool {
	visited[node.ID] = true
	recursionStack[node.ID] = true

	for _, child := range node.Children {
		if !visited[child.ID] {
			if hasCycleUtil(child, visited, recursionStack) {
				return true
			}
		} else if recursionStack[child.ID] {
			return true
		}
	}

	recursionStack[node.ID] = false
	return false
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
