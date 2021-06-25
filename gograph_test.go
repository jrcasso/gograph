package gograph

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

// TestCreateGraph tests gograph operations
func TestCreateGraph(t *testing.T) {
	describe("CreateGraph", t)
	var graph = CreateGraph()

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
	// Empty graph
	describe("CreateDirectedNode", t)
	var graph = CreateGraph()
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
	var graph = Graph{}
	var newNode *Node
	var nodeIndex int

	// Create 10,000 nodes in the graph
	for i := 0; i < 10000; i++ {
		var nodeID = CreateDirectedNodeID()
		newNode = &Node{Edges: nil, ID: nodeID}
		graph.Nodes = append(graph.Nodes, newNode)
	}

	// Select a random one, the 1338th node, and do a lookup on its ID.
	// Expect the function to return the index in the graph for the same node.
	var nodeID = graph.Nodes[1337].ID
	nodeIndex = FindNode(graph, nodeID)
	expectEqualInts(nodeIndex, 1337, t)

	context("the desired node is NOT in the graph", t)
	it("returns an index of -1", t)
	var fakeNodeID = "not-a-true-id"
	nodeIndex = FindNode(graph, fakeNodeID)
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
	var graph = CreateGraph()
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
	var graph = CreateGraph()
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
	var graph = CreateGraph()
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
	var graph = CreateGraph()
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
	var graph = CreateGraph()
	var nodeA, nodeB, nodeC, nodeD, nodeE *DirectedNode
	var sortedNodes []*DirectedNode

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
	sortedNodes = TopologicalSort(graph)
	expectEqualStrings(sortedNodes[0].Values["name"], "nodeD", t)
	expectEqualStrings(sortedNodes[1].Values["name"], "nodeC", t)
	expectEqualStrings(sortedNodes[2].Values["name"], "nodeB", t)
	expectEqualStrings(sortedNodes[3].Values["name"], "nodeA", t)

	//    A
	//   / \
	//  B   C
	//   \ /
	//    D
	graph = CreateGraph()
	graph, nodeA = CreateDirectedNode(graph, map[string]string{"name": "nodeA"}, []*DirectedNode{}, []*DirectedNode{})
	graph, nodeB = CreateDirectedNode(graph, map[string]string{"name": "nodeB"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graph, nodeC = CreateDirectedNode(graph, map[string]string{"name": "nodeC"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graph, _ = CreateDirectedNode(graph, map[string]string{"name": "nodeD"}, []*DirectedNode{nodeB, nodeC}, []*DirectedNode{})
	sortedNodes = TopologicalSort(graph)
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
	graph = CreateGraph()
	graph, nodeA = CreateDirectedNode(graph, map[string]string{"name": "nodeA"}, []*DirectedNode{}, []*DirectedNode{})
	graph, nodeB = CreateDirectedNode(graph, map[string]string{"name": "nodeB"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graph, nodeC = CreateDirectedNode(graph, map[string]string{"name": "nodeC"}, []*DirectedNode{nodeA}, []*DirectedNode{})
	graph, nodeD = CreateDirectedNode(graph, map[string]string{"name": "nodeD"}, []*DirectedNode{nodeB, nodeC}, []*DirectedNode{})
	graph, nodeE = CreateDirectedNode(graph, map[string]string{"name": "nodeE"}, []*DirectedNode{nodeC}, []*DirectedNode{})
	graph, _ = CreateDirectedNode(graph, map[string]string{"name": "nodeF"}, []*DirectedNode{nodeD}, []*DirectedNode{})
	graph, _ = CreateDirectedNode(graph, map[string]string{"name": "nodeG"}, []*DirectedNode{nodeE}, []*DirectedNode{})
	sortedNodes = TopologicalSort(graph)
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
	graphA = CreateGraph()
	graphA = CreateGraph()
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
	graphB = CreateGraph()
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
	var graph = CreateGraph()
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
	graph = CreateGraph()
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
