package godag

import (
	"testing"
)

// TestCreateGraph tests godag operations
func TestCreateGraph(t *testing.T) {
	describe("CreateGraph", t)
	var graph = CreateGraph()

	it("creates a nil root node.", t)
	if graph.RootNode != nil {
		t.Errorf("Failed: expected %s, but found %+v", "nil", graph.RootNode)
		return
	}

	it("doesn't contain any nodes.", t)
	if len(graph.Nodes) != 0 {
		t.Errorf("Failed: expected %d, but found %+v", 0, len(graph.Nodes))
		return
	}

}

// TestCreateNode tests godag operations
func TestCreateNode(t *testing.T) {
	describe("CreateNode", t)
	var graph = CreateGraph()
	var newRootNode, newNode *Node

	context("child and parent aren't specified", t)
	graph, newRootNode = CreateNode(graph, []*Node{}, []*Node{})

	it("creates an edgeless node", t)
	expectEqualInts(0, len(newRootNode.Parents), t)
	expectEqualInts(0, len(newRootNode.Children), t)

	it("assigns root node to the created node", t)
	expectEqualInts(1, len(graph.Nodes), t)
	expectEqualStrings(graph.RootNode.ID, newRootNode.ID, t)

	context("parent nodes are specified", t)
	graph, newNode = CreateNode(graph, []*Node{newRootNode}, []*Node{})

	it("creates a new node with the same parents as specified", t)
	expectEqualInts(1, len(newNode.Parents), t)
	expectEqualStrings(newNode.Parents[0].ID, graph.RootNode.ID, t)

	it("updates existing nodes' child values", t)
	expectEqualInts(2, len(graph.Nodes), t)
	expectEqualStrings(graph.RootNode.Children[0].ID, newNode.ID, t)

	context("child nodes are specified", t)
	graph, newNode = CreateNode(graph, []*Node{}, []*Node{newRootNode})

	it("creates a new node with the same parents as specified", t)
	expectEqualInts(1, len(newNode.Parents), t)
	expectEqualStrings(newNode.Parents[0].ID, graph.RootNode.ID, t)

	it("updates existing nodes' child values", t)
	expectEqualInts(2, len(graph.Nodes), t)
	expectEqualStrings(graph.RootNode.Children[0].ID, newNode.ID, t)

	// it("assigns root node to the created node", t)
	// expectEqualInts(1, len(graph.Nodes), t)
	// expectEqualStrings(graph.RootNode.ID, node.ID, t)

	// graph, node = CreateNode(graph, []*Node{}, []*Node{})
	// expectEqualInts(0, len(node.Parents), t)
	// expectEqualInts(0, len(node.Children), t)
	// expectEqualInts(2, len(graph.Nodes), t)
}
