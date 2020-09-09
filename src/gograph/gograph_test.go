package gograph

import (
	"testing"
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
	describe("CreateDirectedNode", t)
	var graph = CreateGraph()
	var newRootDirectedNode, newDirectedNode *DirectedNode

	context("child and parent aren't specified", t)
	graph, newRootDirectedNode = CreateDirectedNode(graph, []*DirectedNode{}, []*DirectedNode{})

	it("creates an edgeless node", t)
	expectEqualInts(0, len(newRootDirectedNode.Parents), t)
	expectEqualInts(0, len(newRootDirectedNode.Children), t)

	it("assigns root node to the created node", t)
	expectEqualInts(1, len(graph.DirectedNodes), t)
	expectEqualStrings(graph.RootDirectedNode.ID, newRootDirectedNode.ID, t)

	context("parent nodes are specified", t)
	graph, newDirectedNode = CreateDirectedNode(graph, []*DirectedNode{newRootDirectedNode}, []*DirectedNode{})

	it("creates a new node with the same parents as specified", t)
	expectEqualInts(1, len(newDirectedNode.Parents), t)
	expectEqualStrings(newDirectedNode.Parents[0].ID, graph.RootDirectedNode.ID, t)

	it("updates existing nodes' child values", t)
	expectEqualInts(2, len(graph.DirectedNodes), t)
	expectEqualStrings(graph.RootDirectedNode.Children[0].ID, newDirectedNode.ID, t)

	context("child nodes are specified", t)
	graph, newDirectedNode = CreateDirectedNode(graph, []*DirectedNode{}, []*DirectedNode{newRootDirectedNode})

	it("creates a new node with the same parents as specified", t)
	expectEqualInts(1, len(newDirectedNode.Parents), t)
	expectEqualStrings(newDirectedNode.Parents[0].ID, graph.RootDirectedNode.ID, t)

	it("updates existing nodes' child values", t)
	expectEqualInts(2, len(graph.DirectedNodes), t)
	expectEqualStrings(graph.RootDirectedNode.Children[0].ID, newDirectedNode.ID, t)

	// it("assigns root node to the created node", t)
	// expectEqualInts(1, len(graph.DirectedNodes), t)
	// expectEqualStrings(graph.RootDirectedNode.ID, node.ID, t)

	// graph, node = CreateDirectedNode(graph, []*DirectedNode{}, []*DirectedNode{})
	// expectEqualInts(0, len(node.Parents), t)
	// expectEqualInts(0, len(node.Children), t)
	// expectEqualInts(2, len(graph.DirectedNodes), t)
}
