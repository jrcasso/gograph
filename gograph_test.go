package gograph

import (
	"math/rand"
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
