package gograph

import (
	"bytes"
	"io"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

// TestCreateGraph tests gograph operations
func TestCreateGraph(t *testing.T) {
	describe("CreateGraph", t)
	var graph = CreateDirectedGraphHelper()

	it("creates a nil root node.", t)
	if graph.RootDirectedNode != nil {
		t.Errorf("Failed: expected %s, but found %+v", "nil", graph.RootDirectedNode)
		return
	}

	it("doesn't contain any nodes.", t)
	if len(graph.DirectedNodes) != 0 {
		t.Errorf("Failed: expected %d, but found %+v", 0, len(graph.DirectedNodes))
		return
	}
}

// TestCreateDirectedNode tests gograph operations
func TestCreateDirectedNode(t *testing.T) {
	describe("CreateDirectedNode", t)
	var graph = CreateDirectedGraphHelper()
	var newRootDirectedNode, newDirectedNode1, newDirectedNode2, newDirectedNode3 *DirectedNode

	context("child and parent aren't specified", t)
	// Update the empty graph to contain a single, root node:
	//         *r
	graph, newRootDirectedNode = CreateDirectedNode(graph, map[string]string{"foo": "bar", "baz": "biz"}, []*DirectedNode{}, []*DirectedNode{})

	it("updates the graph node references", t)
	expectEqualInts(1, len(graph.DirectedNodes), t)

	it("creates an edgeless node", t)
	expectEqualInts(0, len(newRootDirectedNode.Parents), t)
	expectEqualInts(0, len(newRootDirectedNode.Children), t)

	it("has accessible values", t)
	expectEqualInts(len(graph.RootDirectedNode.Values), 2, t)
	expectEqualStrings(graph.RootDirectedNode.Values["foo"], "bar", t)
	expectEqualStrings(graph.RootDirectedNode.Values["baz"], "biz", t)

	it("has mutable values", t)
	graph.RootDirectedNode.Values["foo"] = "test"
	expectEqualStrings(graph.RootDirectedNode.Values["foo"], "test", t)

	it("assigns root node to the created node", t)
	expectEqualStrings(graph.RootDirectedNode.ID, newRootDirectedNode.ID, t)

	context("parent nodes are specified", t)
	// Update the graph to have a node with the root node specified as a parent:
	//         *r
	//        /
	//       *1
	graph, newDirectedNode1 = CreateDirectedNode(graph, nil, []*DirectedNode{newRootDirectedNode}, []*DirectedNode{})

	it("updates the graph node references", t)
	expectEqualInts(2, len(graph.DirectedNodes), t)

	it("creates a new node with the same parents as specified", t)
	expectEqualInts(1, len(newDirectedNode1.Parents), t)
	expectEqualStrings(newDirectedNode1.Parents[0].ID, graph.RootDirectedNode.ID, t)

	it("updates existing nodes' child values", t)
	expectEqualStrings(graph.RootDirectedNode.Children[0].ID, newDirectedNode1.ID, t)

	it("creates a recursive reference (i.e. a node is its child nodes parent", t)
	expectEqualStrings(graph.RootDirectedNode.Children[0].Parents[0].ID, graph.RootDirectedNode.ID, t)

	context("both parent and child nodes are specified", t)
	// Update the graph to have a node with the previous node specified as a parent:
	//         *r                   *r
	//        /                    / \
	//       *1         ---->     *1  *2
	//                                 \
	//             *3                   *3
	graph, newDirectedNode3 = CreateDirectedNode(graph, nil, []*DirectedNode{}, []*DirectedNode{})
	graph, newDirectedNode2 = CreateDirectedNode(graph, nil, []*DirectedNode{newRootDirectedNode}, []*DirectedNode{newDirectedNode3})

	it("updates the graph node references", t)
	expectEqualInts(4, len(graph.DirectedNodes), t)

	it("creates a new node with the same parents as specified", t)
	expectEqualInts(1, len(newDirectedNode2.Parents), t)
	expectEqualStrings(newDirectedNode2.Parents[0].ID, graph.RootDirectedNode.ID, t)
	expectEqualStrings(newDirectedNode3.Parents[0].ID, newDirectedNode2.ID, t)

	it("updates existing nodes' child values", t)
	expectEqualStrings(graph.RootDirectedNode.Children[1].ID, newDirectedNode2.ID, t)
	expectEqualStrings(newDirectedNode2.Children[0].ID, newDirectedNode3.ID, t)
}

// TestFindNode ensures that FindNode will find a node with a passed ID, or return -1 otherwise
func TestFindNode(t *testing.T) {
	describe("FindNode", t)

	context("the desired node is in the graph", t)
	it("should find the correct node", t)
	var graph = CreateGraph()
	var newNode *Node
	var nodeIndex int

	// Create 10,000 nodes in the graph
	for i := 0; i < 10000; i++ {
		var nodeID = CreateNodeID()
		newNode = &Node{Edges: nil, Values: make(map[string]string), ID: nodeID}
		graph.Nodes = append(graph.Nodes, newNode)
	}

	// Select a random one, the 1338th node, and do a lookup on its ID.
	// Expect the function to return the index in the graph for the same node.
	var nodeID = graph.Nodes[1337].ID
	nodeIndex, _ = FindNode(graph, nodeID)
	expectEqualInts(nodeIndex, 1337, t)

	context("the desired node is NOT in the graph", t)
	it("returns an index of -1", t)
	var fakeNodeID = "not-a-true-id"
	nodeIndex, _ = FindNode(graph, fakeNodeID)
	expectEqualInts(nodeIndex, -1, t)
}

func TestCreateDirectedNodeID(t *testing.T) {
	describe("CreateDirectedNodeID", t)
	it("can create unique IDs", t)

	// It's something like 1 in 2^80 chance that a collision becomes a concern. 10000 will do fine.
	var IDs = [10000]string{}
	var newID string
	for i := 1; i < 10000; i++ {
		newID = CreateDirectedNodeID()
		for _, ID := range IDs {
			if newID == ID {
				t.Errorf("Failed: attempted to use redundant node ID %s", ID)
				return

			}
			IDs[i] = newID
		}
	}
}

func TestFindDirectedNode(t *testing.T) {
	describe("FindDirectedNode", t)

	context("the desired node is in the graph", t)
	rand.Seed(time.Now().UnixNano())
	var graph = CreateDirectedGraphHelper()
	var parentDirectedNode, node *DirectedNode
	graph, parentDirectedNode = CreateDirectedNode(graph, nil, []*DirectedNode{}, []*DirectedNode{})

	var nodeCount = rand.Intn(1000)
	var randIndex = rand.Intn(nodeCount)
	var searchID = ""
	for i := 0; i < nodeCount; i++ {
		graph, node = CreateDirectedNode(graph, nil, []*DirectedNode{parentDirectedNode}, []*DirectedNode{})
		parentDirectedNode = node
		if randIndex-1 == i {
			searchID = node.ID
		}
	}

	it("returns the correct index", t)
	var index int
	var foundNode DirectedNode

	index, foundNode = FindDirectedNode(graph, searchID)
	expectEqualInts(index, randIndex, t)
	expectEqualStrings(graph.DirectedNodes[index].ID, searchID, t)

	it("returns the same node as specified by the index", t)
	expectEqualStrings(foundNode.ID, graph.DirectedNodes[index].ID, t)

	context("the desired node is NOT in the graph", t)
	it("returns an index of -1", t)
	var fakeNodeID = "not-a-true-id"
	index, foundNode = FindDirectedNode(graph, fakeNodeID)
	expectEqualInts(index, -1, t)
}

func TestFindNodesByValues(t *testing.T) {
	describe("FindNodesByValues", t)

	context("the desired node is in the graph", t)
	rand.Seed(time.Now().UnixNano())
	var graph = CreateDirectedGraphHelper()
	var parentDirectedNode, node, targetNode *DirectedNode
	var foundNodes []*DirectedNode

	graph, parentDirectedNode = CreateDirectedNode(graph, nil, []*DirectedNode{}, []*DirectedNode{})

	// Create a chain of 10,000 nodes
	for i := 0; i < 10000; i++ {
		s := strconv.Itoa(i)
		graph, node = CreateDirectedNode(graph, map[string]string{"index": s}, []*DirectedNode{parentDirectedNode}, []*DirectedNode{})
		parentDirectedNode = node
		if i == 1337 {
			targetNode = node
		}
	}

	it("returns the correct node", t)
	foundNodes = FindNodesByValues(graph, map[string]string{"index": "1337"})
	expectEqualInts(len(foundNodes), 1, t)
	expectEqualStrings(foundNodes[0].ID, targetNode.ID, t)
}

func TestCreateDirectedEdge(t *testing.T) {
	describe("FindDirectedNode", t)
	rand.Seed(time.Now().UnixNano())
	var graph = CreateDirectedGraphHelper()
	var parentNode, childNode, cousinNode *DirectedNode
	graph, parentNode = CreateDirectedNode(graph, nil, []*DirectedNode{}, []*DirectedNode{})
	graph, childNode = CreateDirectedNode(graph, nil, []*DirectedNode{}, []*DirectedNode{})

	graph, parentNode, childNode = CreateDirectedEdge(graph, parentNode, childNode)

	it("should correctly assign a new child to the parent node, and a new parent node to the child", t)
	expectEqualStrings(graph.DirectedNodes[0].ID, parentNode.ID, t)
	expectEqualStrings(graph.DirectedNodes[1].ID, childNode.ID, t)
	expectEqualStrings(graph.DirectedNodes[1].Parents[0].ID, parentNode.ID, t)
	expectEqualStrings(graph.DirectedNodes[0].Children[0].ID, childNode.ID, t)

	graph, cousinNode = CreateDirectedNode(graph, nil, []*DirectedNode{}, []*DirectedNode{})
	graph, _, cousinNode = CreateDirectedEdge(graph, parentNode, cousinNode)
	graph, _, cousinNode = CreateDirectedEdge(graph, childNode, cousinNode)

	it("should correctly create edges between arbitrary nodes", t)
	expectEqualStrings(graph.DirectedNodes[2].ID, cousinNode.ID, t)
	expectEqualStrings(graph.DirectedNodes[0].Children[1].ID, cousinNode.ID, t)
	expectEqualStrings(graph.DirectedNodes[1].Children[0].ID, cousinNode.ID, t)
}

func TestDeleteDirectedEdge(t *testing.T) {
	describe("DeleteDirectedEdge", t)
	rand.Seed(time.Now().UnixNano())
	var graph = CreateDirectedGraphHelper()
	var parentNode, childNode *DirectedNode
	graph, parentNode = CreateDirectedNode(graph, nil, []*DirectedNode{}, []*DirectedNode{})
	graph, childNode = CreateDirectedNode(graph, nil, []*DirectedNode{}, []*DirectedNode{})
	graph, parentNode, childNode = CreateDirectedEdge(graph, parentNode, childNode)

	it("should delete the edge relationship between a parent and child node", t)
	DeleteDirectedEdge(graph, parentNode, childNode)

	expectEqualInts(0, len(parentNode.Children), t)
	expectEqualInts(0, len(childNode.Parents), t)
}

// TestTopologicalSort tests that the function returns one possible topological ordering of the provided dag
func TestTopologicalSort(t *testing.T) {
	describe("TopologicalSort", t)
	it("should correctly sort the graph", t)
	var graph = CreateDirectedGraphHelper()
	var nodeA, nodeB, nodeC, nodeD, nodeE *DirectedNode
	var sortedNodes []*DirectedNode
	var err error

	//    A
	//   / \
	//  B   |
	//    \ |
	//      C
	//      |
	//      D
	graph, nodeA = CreateDirectedNode(graph, map[string]string{"name": "nodeA"}, []*DirectedNode{}, []*DirectedNode{})
	graph, nodeB = CreateDirectedNode(graph, map[string]string{"name": "nodeB"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graph, nodeC = CreateDirectedNode(graph, map[string]string{"name": "nodeC"}, []*DirectedNode{nodeA, nodeB}, []*DirectedNode{})
	graph, _ = CreateDirectedNode(graph, map[string]string{"name": "nodeD"}, []*DirectedNode{nodeC}, []*DirectedNode{})
	sortedNodes, err = TopologicalSort(graph)
	if err != nil {
		t.Errorf("Failed: unexpected error: %v", err)
		return
	}
	expectEqualStrings(sortedNodes[0].Values["name"], "nodeD", t)
	expectEqualStrings(sortedNodes[1].Values["name"], "nodeC", t)
	expectEqualStrings(sortedNodes[2].Values["name"], "nodeB", t)
	expectEqualStrings(sortedNodes[3].Values["name"], "nodeA", t)

	//    A
	//   / \
	//  B   C
	//   \ /
	//    D
	graph = CreateDirectedGraphHelper()
	graph, nodeA = CreateDirectedNode(graph, map[string]string{"name": "nodeA"}, []*DirectedNode{}, []*DirectedNode{})
	graph, nodeB = CreateDirectedNode(graph, map[string]string{"name": "nodeB"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graph, nodeC = CreateDirectedNode(graph, map[string]string{"name": "nodeC"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graph, _ = CreateDirectedNode(graph, map[string]string{"name": "nodeD"}, []*DirectedNode{nodeB, nodeC}, []*DirectedNode{})
	sortedNodes, err = TopologicalSort(graph)
	if err != nil {
		t.Errorf("Failed: unexpected error: %v", err)
		return
	}
	expectEqualStrings(sortedNodes[0].Values["name"], "nodeD", t)
	expectEqualStrings(sortedNodes[1].Values["name"], "nodeB", t)
	expectEqualStrings(sortedNodes[2].Values["name"], "nodeC", t)
	expectEqualStrings(sortedNodes[3].Values["name"], "nodeA", t)

	//    A
	//   / \
	//  B   C
	//   \ / \
	//    D   E
	//    |   |
	//    F   G
	graph = CreateDirectedGraphHelper()
	graph, nodeA = CreateDirectedNode(graph, map[string]string{"name": "nodeA"}, []*DirectedNode{}, []*DirectedNode{})
	graph, nodeB = CreateDirectedNode(graph, map[string]string{"name": "nodeB"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graph, nodeC = CreateDirectedNode(graph, map[string]string{"name": "nodeC"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graph, nodeD = CreateDirectedNode(graph, map[string]string{"name": "nodeD"}, []*DirectedNode{nodeB, nodeC}, []*DirectedNode{})
	graph, nodeE = CreateDirectedNode(graph, map[string]string{"name": "nodeE"}, []*DirectedNode{nodeC}, []*DirectedNode{})
	graph, _ = CreateDirectedNode(graph, map[string]string{"name": "nodeF"}, []*DirectedNode{nodeD}, []*DirectedNode{})
	graph, _ = CreateDirectedNode(graph, map[string]string{"name": "nodeG"}, []*DirectedNode{nodeE}, []*DirectedNode{})
	sortedNodes, err = TopologicalSort(graph)
	if err != nil {
		t.Errorf("Failed: unexpected error: %v", err)
		return
	}
	expectEqualStrings(sortedNodes[0].Values["name"], "nodeF", t)
	expectEqualStrings(sortedNodes[1].Values["name"], "nodeG", t)
	expectEqualStrings(sortedNodes[2].Values["name"], "nodeD", t)
	expectEqualStrings(sortedNodes[3].Values["name"], "nodeE", t)
	expectEqualStrings(sortedNodes[4].Values["name"], "nodeB", t)
	expectEqualStrings(sortedNodes[5].Values["name"], "nodeC", t)
	expectEqualStrings(sortedNodes[6].Values["name"], "nodeA", t)
}

func TestCreateAdjecencyMatrix(t *testing.T) {
	describe("CreateAdjecencyMatrix", t)
	it("should correctly describe the edge relationships between nodes", t)
	var graphA, graphB DirectedGraph
	var adjecencyMatrix [][]int
	var nodeA, nodeB, nodeC, nodeD, nodeE *DirectedNode
	var expectedMatrix = [][]int{
		{0, 1, 1, 0},
		{0, 0, 0, 1},
		{0, 0, 0, 1},
		{0, 0, 0, 0},
	}

	//    A
	//   / \
	//  B   C
	//   \ /
	//    D

	graphA = CreateDirectedGraphHelper()
	graphA, nodeA = CreateDirectedNode(graphA, map[string]string{"name": "nodeA"}, []*DirectedNode{}, []*DirectedNode{})
	graphA, nodeB = CreateDirectedNode(graphA, map[string]string{"name": "nodeB"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graphA, nodeC = CreateDirectedNode(graphA, map[string]string{"name": "nodeC"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graphA, _ = CreateDirectedNode(graphA, map[string]string{"name": "nodeD"}, []*DirectedNode{nodeB, nodeC}, []*DirectedNode{})

	//    A
	//   / \
	//  B   C
	//   \ / \
	//    D   E
	//    |   |
	//    F   G
	expectedMatrix = [][]int{
		{0, 1, 1, 0, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 1, 0, 0},
		{0, 0, 0, 0, 0, 1, 0},
		{0, 0, 0, 0, 0, 0, 1},
		{0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0},
	}
	graphB = CreateDirectedGraphHelper()
	graphB, nodeA = CreateDirectedNode(graphB, map[string]string{"name": "nodeA"}, []*DirectedNode{}, []*DirectedNode{})
	graphB, nodeB = CreateDirectedNode(graphB, map[string]string{"name": "nodeB"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graphB, nodeC = CreateDirectedNode(graphB, map[string]string{"name": "nodeC"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graphB, nodeD = CreateDirectedNode(graphB, map[string]string{"name": "nodeD"}, []*DirectedNode{nodeB, nodeC}, []*DirectedNode{})
	graphB, nodeE = CreateDirectedNode(graphB, map[string]string{"name": "nodeE"}, []*DirectedNode{nodeC}, []*DirectedNode{})
	graphB, _ = CreateDirectedNode(graphB, map[string]string{"name": "nodeF"}, []*DirectedNode{nodeD}, []*DirectedNode{})
	graphB, _ = CreateDirectedNode(graphB, map[string]string{"name": "nodeG"}, []*DirectedNode{nodeE}, []*DirectedNode{})

	for _, graph := range []DirectedGraph{graphA, graphB} {
		adjecencyMatrix = CreateAdjecencyMatrix(graph)
		for i := range adjecencyMatrix {
			for j := range adjecencyMatrix {
				if adjecencyMatrix[i][j] != expectedMatrix[i][j] {
					t.Errorf("Failed: expected %+v, but found %+v", expectedMatrix, adjecencyMatrix)
					return
				}
			}
		}
	}
}

func TestCreateIncidenceMatrix(t *testing.T) {
	describe("CreateIncidenceMatrix", t)
	it("should correctly describe the edge relationships between nodes", t)
	var graph = CreateDirectedGraphHelper()
	var incidenceMatrix [][]int
	var nodeA, nodeB, nodeC *DirectedNode
	var expectedMatrix = [][]int{
		{0, 1, 1, 0},
		{-1, 0, 0, 1},
		{-1, 0, 0, 1},
		{0, -1, -1, 0},
	}

	//    A
	//   / \
	//  B   C
	//   \ /
	//    D
	graph, nodeA = CreateDirectedNode(graph, map[string]string{"name": "nodeA"}, []*DirectedNode{}, []*DirectedNode{})
	graph, nodeB = CreateDirectedNode(graph, map[string]string{"name": "nodeB"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graph, nodeC = CreateDirectedNode(graph, map[string]string{"name": "nodeC"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graph, _ = CreateDirectedNode(graph, map[string]string{"name": "nodeD"}, []*DirectedNode{nodeB, nodeC}, []*DirectedNode{})
	incidenceMatrix = CreateIncidenceMatrix(graph)

	for i := range incidenceMatrix {
		for j := range incidenceMatrix {
			if incidenceMatrix[i][j] != expectedMatrix[i][j] {
				t.Errorf("Failed: expected %+v, but found %+v", expectedMatrix, incidenceMatrix)
				return
			}
		}
	}
}

func TestHasCycle(t *testing.T) {
	describe("HasCycle", t)

	context("graph has no cycles", t)
	it("should return false", t)
	var graph = CreateDirectedGraphHelper()
	var nodeA, nodeB, nodeC, nodeD *DirectedNode

	// Create a DAG:
	//    A
	//   / \
	//  B   C
	//   \ /
	//    D
	graph, nodeA = CreateDirectedNode(graph, map[string]string{"name": "nodeA"}, []*DirectedNode{}, []*DirectedNode{})
	graph, nodeB = CreateDirectedNode(graph, map[string]string{"name": "nodeB"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graph, nodeC = CreateDirectedNode(graph, map[string]string{"name": "nodeC"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graph, _ = CreateDirectedNode(graph, map[string]string{"name": "nodeD"}, []*DirectedNode{nodeB, nodeC}, []*DirectedNode{})

	if HasCycle(graph) {
		t.Errorf("Failed: expected no cycle in acyclic graph")
	}

	context("graph has a cycle", t)
	it("should return true", t)
	// Create a graph with a cycle:
	//    A
	//   / \
	//  B-->C
	//  ^   |
	//  |   v
	//  D<--E
	graph = CreateDirectedGraphHelper()
	graph, nodeA = CreateDirectedNode(graph, map[string]string{"name": "nodeA"}, []*DirectedNode{}, []*DirectedNode{})
	graph, nodeB = CreateDirectedNode(graph, map[string]string{"name": "nodeB"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graph, nodeC = CreateDirectedNode(graph, map[string]string{"name": "nodeC"}, []*DirectedNode{nodeA, nodeB}, []*DirectedNode{})
	graph, nodeD = CreateDirectedNode(graph, map[string]string{"name": "nodeD"}, []*DirectedNode{}, []*DirectedNode{})
	graph, _ = CreateDirectedNode(graph, map[string]string{"name": "nodeE"}, []*DirectedNode{nodeC}, []*DirectedNode{})

	// Add edges to create a cycle: B <- D <- E <- C <- B
	graph, _, _ = CreateDirectedEdge(graph, nodeD, nodeB)
	graph, _, _ = CreateDirectedEdge(graph, nodeC, nodeD)

	if !HasCycle(graph) {
		t.Errorf("Failed: expected to detect cycle in cyclic graph")
	}
}

func TestFindPath(t *testing.T) {
	describe("FindPath", t)

	context("path exists", t)
	it("should return the correct path", t)
	var graph = CreateDirectedGraphHelper()
	var nodeA, nodeB, nodeC, nodeD *DirectedNode

	// Create a graph:
	//    A
	//   / \
	//  B   C
	//  |   |
	//  D   E
	graph, nodeA = CreateDirectedNode(graph, map[string]string{"name": "nodeA"}, []*DirectedNode{}, []*DirectedNode{})
	graph, nodeB = CreateDirectedNode(graph, map[string]string{"name": "nodeB"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graph, nodeC = CreateDirectedNode(graph, map[string]string{"name": "nodeC"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graph, nodeD = CreateDirectedNode(graph, map[string]string{"name": "nodeD"}, []*DirectedNode{nodeB}, []*DirectedNode{})
	graph, _ = CreateDirectedNode(graph, map[string]string{"name": "nodeE"}, []*DirectedNode{nodeC}, []*DirectedNode{})

	path, err := FindPath(graph, nodeA.ID, nodeD.ID)
	if err != nil {
		t.Errorf("Failed: unexpected error: %v", err)
		return
	}

	it("should return a path of length 3", t)
	expectEqualInts(len(path), 3, t)

	it("should return the correct path nodes", t)
	expectEqualStrings(path[0].Values["name"], "nodeA", t)
	expectEqualStrings(path[1].Values["name"], "nodeB", t)
	expectEqualStrings(path[2].Values["name"], "nodeD", t)

	// A path from D to C is impossible, so we should get an error
	context("path does not exist", t)
	it("should return an error", t)
	_, err = FindPath(graph, nodeD.ID, nodeC.ID)
	if err == nil {
		t.Errorf("Failed: expected error for non-existent path")
	}
}

func TestFindStronglyConnectedComponents(t *testing.T) {
	describe("FindStronglyConnectedComponents", t)

	context("a graph has multiple components", t)
	var graph = CreateDirectedGraphHelper()
	var nodeA, nodeB, nodeC, nodeD, nodeE, nodeF *DirectedNode

	graph, nodeA = CreateDirectedNode(graph, map[string]string{"name": "nodeA"}, []*DirectedNode{}, []*DirectedNode{})
	graph, nodeB = CreateDirectedNode(graph, map[string]string{"name": "nodeB"}, []*DirectedNode{}, []*DirectedNode{})
	graph, nodeC = CreateDirectedNode(graph, map[string]string{"name": "nodeC"}, []*DirectedNode{}, []*DirectedNode{})
	graph, nodeD = CreateDirectedNode(graph, map[string]string{"name": "nodeD"}, []*DirectedNode{}, []*DirectedNode{})
	graph, nodeE = CreateDirectedNode(graph, map[string]string{"name": "nodeE"}, []*DirectedNode{}, []*DirectedNode{})
	graph, nodeF = CreateDirectedNode(graph, map[string]string{"name": "nodeF"}, []*DirectedNode{}, []*DirectedNode{})

	// Create two components with cycles
	graph, _, _ = CreateDirectedEdge(graph, nodeA, nodeB)
	graph, _, _ = CreateDirectedEdge(graph, nodeB, nodeC)
	graph, _, _ = CreateDirectedEdge(graph, nodeC, nodeA)

	graph, _, _ = CreateDirectedEdge(graph, nodeD, nodeE)
	graph, _, _ = CreateDirectedEdge(graph, nodeE, nodeF)
	graph, _, _ = CreateDirectedEdge(graph, nodeF, nodeD)

	// Connect components with a one-way edge
	graph, _, _ = CreateDirectedEdge(graph, nodeC, nodeD)
	components := FindStronglyConnectedComponents(graph)

	it("should have 2 components", t)
	if len(components) != 2 {
		t.Errorf("Expected 2 components, got %d", len(components))
		return
	}

	compSizes := make([]int, 0, len(components))
	for _, comp := range components {
		compSizes = append(compSizes, len(comp))
	}

	it("should return components that each have 3 nodes", t)
	if compSizes[0] != 3 || compSizes[1] != 3 {
		t.Errorf("Expected components of size 3, got %v", compSizes)
	}
}

// TestBasicGraph tests the undirected graph operations
func TestBasicGraph(t *testing.T) {
	describe("Basic Undirected Graph", t)

	it("creates an empty graph", t)
	graph := CreateGraph()
	expectEqualInts(len(graph.Nodes), 0, t)

	context("when adding nodes", t)
	it("adds nodes correctly", t)
	graph, node1 := CreateNode(graph, map[string]string{"name": "node1"})
	graph, node2 := CreateNode(graph, map[string]string{"name": "node2"})
	graph, node3 := CreateNode(graph, map[string]string{"name": "node3"})

	expectEqualInts(len(graph.Nodes), 3, t)
	expectEqualStrings(node1.Values["name"], "node1", t)
	expectEqualStrings(node2.Values["name"], "node2", t)
	expectEqualStrings(node3.Values["name"], "node3", t)

	context("when adding edges", t)
	it("creates undirected edges correctly", t)
	graph = AddUndirectedEdge(graph, node1, node2)
	graph = AddUndirectedEdge(graph, node2, node3)

	expectEqualInts(len(node1.Edges), 1, t)
	expectEqualInts(len(node2.Edges), 2, t)
	expectEqualInts(len(node3.Edges), 1, t)

	expectEqualStrings(node1.Edges[0].ID, node2.ID, t)
	expectEqualStrings(node2.Edges[0].ID, node1.ID, t)
	expectEqualStrings(node2.Edges[1].ID, node3.ID, t)
	expectEqualStrings(node3.Edges[0].ID, node2.ID, t)

	context("when removing edges", t)
	it("removes undirected edges correctly", t)
	graph = RemoveUndirectedEdge(graph, node1, node2)

	expectEqualInts(len(node1.Edges), 0, t)
	expectEqualInts(len(node2.Edges), 1, t)
	expectEqualStrings(node2.Edges[0].ID, node3.ID, t)

	context("when finding nodes", t)
	it("finds nodes by ID correctly", t)
	index, foundNode := FindNode(graph, node2.ID)
	expectEqualInts(index, 1, t)
	expectEqualStrings(foundNode.ID, node2.ID, t)

	index, foundNode = FindNode(graph, "non-existent-id")
	expectEqualInts(index, -1, t)
	if foundNode != nil {
		t.Errorf("Expected nil node for non-existent ID, got %+v", foundNode)
	}
}

// TestGenericGraph tests the generic flexible graph operations
func TestGenericGraph(t *testing.T) {
	describe("Generic Flexible Graph", t)

	it("creates a graph with specified edge types", t)
	edgeTypes := []string{"Friends", "Enemies", "Colleagues"}
	graph := CreateGenericGraph(edgeTypes)

	expectEqualInts(len(graph.Nodes), 0, t)

	context("when adding nodes", t)
	it("adds nodes with correct edge types", t)

	// Create the first node manually since we need to set up the edge types
	nodeID := CreateNodeID()
	edgeTypeMap := make(map[string]int)
	for i, edgeType := range edgeTypes {
		edgeTypeMap[edgeType] = i
	}

	edges := make([][]*GenericNode, len(edgeTypes))
	firstNode := &GenericNode{
		Edges:     edges,
		EdgeTypes: edgeTypeMap,
		Values:    map[string]string{"name": "node1"},
		ID:        nodeID,
	}

	graph.Nodes = append(graph.Nodes, firstNode)
	graph.NodeMap[firstNode.ID] = firstNode

	graph, node2 := CreateGenericNode(graph, map[string]string{"name": "node2"})
	graph, node3 := CreateGenericNode(graph, map[string]string{"name": "node3"})

	expectEqualInts(len(graph.Nodes), 3, t)
	expectEqualStrings(graph.Nodes[0].Values["name"], "node1", t)
	expectEqualStrings(node2.Values["name"], "node2", t)
	expectEqualStrings(node3.Values["name"], "node3", t)

	context("when adding edges", t)
	it("creates edges of different types correctly", t)

	var err error
	graph, err = AddEdge(graph, firstNode, node2, "Friends")
	if err != nil {
		t.Errorf("Failed to add edge: %v", err)
	}

	graph, err = AddEdge(graph, firstNode, node3, "Enemies")
	if err != nil {
		t.Errorf("Failed to add edge: %v", err)
	}

	graph, err = AddEdge(graph, node2, node3, "Colleagues")
	if err != nil {
		t.Errorf("Failed to add edge: %v", err)
	}

	_, err = AddEdge(graph, node2, firstNode, "InvalidType")
	if err == nil {
		t.Errorf("Expected error for invalid edge type, got nil")
	}

	context("when retrieving related nodes", t)
	it("gets nodes related by a specific edge type", t)

	friends, err := GetRelatedNodes(firstNode, "Friends")
	if err != nil {
		t.Errorf("Failed to get related nodes: %v", err)
	}
	expectEqualInts(len(friends), 1, t)
	expectEqualStrings(friends[0].ID, node2.ID, t)

	enemies, err := GetRelatedNodes(firstNode, "Enemies")
	if err != nil {
		t.Errorf("Failed to get related nodes: %v", err)
	}
	expectEqualInts(len(enemies), 1, t)
	expectEqualStrings(enemies[0].ID, node3.ID, t)

	it("handles invalid relationship types", t)
	_, err = GetRelatedNodes(firstNode, "InvalidType")
	if err == nil {
		t.Errorf("Expected error for invalid relationship type, got nil")
	}

	context("when removing edges", t)
	it("removes edges correctly", t)

	graph, err = RemoveEdge(graph, firstNode, node2, "Friends")
	if err != nil {
		t.Errorf("Failed to remove edge: %v", err)
	}

	friends, _ = GetRelatedNodes(firstNode, "Friends")
	expectEqualInts(len(friends), 0, t)

	context("when converting between graph types", t)
	it("converts DirectedGraph to GenericGraph correctly", t)

	dirGraph := CreateDirectedGraphHelper()
	dirGraph, dirNode1 := CreateDirectedNode(dirGraph, map[string]string{"name": "dirNode1"}, []*DirectedNode{}, []*DirectedNode{})
	dirGraph, dirNode2 := CreateDirectedNode(dirGraph, map[string]string{"name": "dirNode2"}, []*DirectedNode{}, []*DirectedNode{})
	dirGraph, dirNode3 := CreateDirectedNode(dirGraph, map[string]string{"name": "dirNode3"}, []*DirectedNode{}, []*DirectedNode{})

	// Add edges for parent-child relationships
	dirGraph, _, _ = CreateDirectedEdge(dirGraph, dirNode1, dirNode2)
	dirGraph, _, _ = CreateDirectedEdge(dirGraph, dirNode2, dirNode3)
	genGraph := ConvertDirectedToGeneric(dirGraph)

	expectEqualInts(len(genGraph.Nodes), 3, t)

	var genNode1, genNode2, genNode3 *GenericNode
	for _, node := range genGraph.Nodes {
		if node.Values["name"] == "dirNode1" {
			genNode1 = node
		} else if node.Values["name"] == "dirNode2" {
			genNode2 = node
		} else if node.Values["name"] == "dirNode3" {
			genNode3 = node
		}
	}

	if genNode1 == nil || genNode2 == nil || genNode3 == nil {
		t.Errorf("Failed to find all nodes in converted generic graph")
		return
	}

	children, err := GetRelatedNodes(genNode1, "Children")
	if err != nil {
		t.Errorf("Failed to get children: %v", err)
		return
	}
	expectEqualInts(len(children), 1, t)
	expectEqualStrings(children[0].Values["name"], "dirNode2", t)

	parents, err := GetRelatedNodes(genNode2, "Parents")
	if err != nil {
		t.Errorf("Failed to get parents: %v", err)
		return
	}
	expectEqualInts(len(parents), 1, t)
	expectEqualStrings(parents[0].Values["name"], "dirNode1", t)

	parents, err = GetRelatedNodes(genNode3, "Parents")
	if err != nil {
		t.Errorf("Failed to get parents: %v", err)
		return
	}
	expectEqualInts(len(parents), 1, t)

	it("converts GenericGraph back to DirectedGraph correctly", t)
	convertedDirGraph, err := ConvertGenericToDirected(genGraph)
	if err != nil {
		t.Errorf("Failed to convert back to directed graph: %v", err)
		return
	}

	expectEqualInts(len(convertedDirGraph.DirectedNodes), 3, t)
	var convDirNode1, convDirNode2, convDirNode3 *DirectedNode
	for _, node := range convertedDirGraph.DirectedNodes {
		if node.Values["name"] == "dirNode1" {
			convDirNode1 = node
		} else if node.Values["name"] == "dirNode2" {
			convDirNode2 = node
		} else if node.Values["name"] == "dirNode3" {
			convDirNode3 = node
		}
	}

	if convDirNode1 == nil || convDirNode2 == nil || convDirNode3 == nil {
		t.Errorf("Failed to find all nodes in converted directed graph")
		return
	}

	expectEqualInts(len(convDirNode1.Children), 1, t)
	expectEqualStrings(convDirNode1.Children[0].Values["name"], "dirNode2", t)

	expectEqualInts(len(convDirNode2.Parents), 1, t)
	expectEqualStrings(convDirNode2.Parents[0].Values["name"], "dirNode1", t)

	expectEqualInts(len(convDirNode2.Children), 1, t)
	expectEqualStrings(convDirNode2.Children[0].Values["name"], "dirNode3", t)

	expectEqualInts(len(convDirNode3.Parents), 1, t)
	expectEqualStrings(convDirNode3.Parents[0].Values["name"], "dirNode2", t)
}

// TestIsAntisymmetricMatrix tests the IsAntisymmetricMatrix function
func TestIsAntisymmetricMatrix(t *testing.T) {
	describe("IsAntisymmetricMatrix", t)

	context("when matrix is antisymmetric", t)
	it("should return true", t)
	antisymmetricMatrix := [][]int{
		{0, 1, -2},
		{-1, 0, 3},
		{2, -3, 0},
	}

	result := IsAntisymmetricMatrix(antisymmetricMatrix)
	if !result {
		t.Errorf("Expected antisymmetric matrix to return true, got false")
	}

	context("when matrix is not antisymmetric", t)
	it("should return false for non-zero diagonal", t)

	nonZeroDiagonal := [][]int{
		{1, 1, -2},
		{-1, 0, 3},
		{2, -3, 0},
	}

	result = IsAntisymmetricMatrix(nonZeroDiagonal)
	if result {
		t.Errorf("Expected matrix with non-zero diagonal to return false, got true")
	}

	it("should return false for non-antisymmetric elements", t)

	nonAntisymmetric := [][]int{
		{0, 1, -2},
		{-1, 0, 3},
		{2, 4, 0},
	}

	result = IsAntisymmetricMatrix(nonAntisymmetric)
	if result {
		t.Errorf("Expected non-antisymmetric matrix to return false, got true")
	}

	context("when matrix is not square", t)
	it("should return false", t)

	nonSquare := [][]int{
		{0, 1, -2},
		{-1, 0, 3},
	}

	result = IsAntisymmetricMatrix(nonSquare)
	if result {
		t.Errorf("Expected non-square matrix to return false, got true")
	}

	context("when matrix is empty", t)
	it("should handle empty matrix", t)

	empty := [][]int{}

	result = IsAntisymmetricMatrix(empty)
	if result {
		t.Errorf("Expected empty matrix to return false, got true")
	}
}

// TestPrintMatrix tests the PrintMatrix function
func TestPrintMatrix(t *testing.T) {
	describe("PrintMatrix", t)
	it("should print the matrix correctly", t)
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	PrintMatrix(matrix)
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()
	expectedOutput := "[1 2 3]\n[4 5 6]\n[7 8 9]\n"
	if output != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, output)
	}

	it("should handle empty matrix", t)
	emptyMatrix := [][]int{}
	r, w, _ = os.Pipe()
	os.Stdout = w
	PrintMatrix(emptyMatrix)
	w.Close()
	os.Stdout = oldStdout
	buf.Reset()
	io.Copy(&buf, r)
	output = buf.String()
	if output != "" {
		t.Errorf("Expected empty output for empty matrix, got: %s", output)
	}
}
