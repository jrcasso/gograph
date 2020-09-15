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
	// Empty graph
	describe("CreateDirectedNode", t)
	var graph = CreateGraph()
	var newRootDirectedNode, newDirectedNode1, newDirectedNode2, newDirectedNode3 *DirectedNode

	context("child and parent aren't specified", t)
	// Update the empty graph to contain a single, root node:
	//         *r
	graph, newRootDirectedNode = CreateDirectedNode(graph, []*DirectedNode{}, []*DirectedNode{})

	it("updates the graph node references", t)
	expectEqualInts(1, len(graph.DirectedNodes), t)

	it("creates an edgeless node", t)
	expectEqualInts(0, len(newRootDirectedNode.Parents), t)
	expectEqualInts(0, len(newRootDirectedNode.Children), t)

	it("assigns root node to the created node", t)
	expectEqualStrings(graph.RootDirectedNode.ID, newRootDirectedNode.ID, t)

	context("parent nodes are specified", t)
	// Update the graph to have a node with the root node specified as a parent:
	//         *r
	//        /
	//       *1
	graph, newDirectedNode1 = CreateDirectedNode(graph, []*DirectedNode{newRootDirectedNode}, []*DirectedNode{})

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
	graph, newDirectedNode3 = CreateDirectedNode(graph, []*DirectedNode{}, []*DirectedNode{})
	graph, newDirectedNode2 = CreateDirectedNode(graph, []*DirectedNode{newRootDirectedNode}, []*DirectedNode{newDirectedNode3})

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
