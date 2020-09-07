package godag

import (
	"testing"
)

// TestCreateGraph tests godag operations
func TestCreateGraph(t *testing.T) {
	go describe("CreateGraph", t)
	graph := CreateGraph()

	it("creates a single root node.", t)
	expectEqualInts(1, len(graph.Nodes), t)

	it("creates a single root node with an ID of length 40.", t)
	rootNodeID := graph.RootNode.ID
	expectEqualInts(40, len(rootNodeID), t)

	it("creates a single root node without a parent.", t)
	rootNode := graph.Nodes[0]
	// There should only be one empty string (the null parent)
	expectEqualInts(1, len(rootNode.ParentIDs), t)
	expectEqualStrings("", rootNode.ParentIDs[0], t)
}

// TestCreateNode tests godag operations
func TestCreateNode(t *testing.T) {
	go describe("CreateNode", t)

	it("creates an edgeless node", t)
	when("child and parent aren't specified", t)
	newNode := CreateNode("", "")
	expectEqualInts(1, len(newNode.ParentIDs), t)
	expectEqualInts(1, len(newNode.ChildIDs), t)

	it("creates a node with only a parent edge", t)
	when("parent is specified and child is not", t)
	newNode = CreateNode("ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMN", "")
	expectEqualInts(1, len(newNode.ParentIDs), t)
	expectEqualInts(40, len(newNode.ParentIDs[0]), t)

	it("creates a node with only a child edge", t)
	when("child is specified and parent is not", t)
	newNode = CreateNode("", "ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMN")
	expectEqualInts(1, len(newNode.ChildIDs), t)
	expectEqualInts(40, len(newNode.ChildIDs[0]), t)

	it("creates a node with both a child edge and a parent edge", t)
	when("both the child and parent are specified", t)
	newNode = CreateNode("ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMN", "ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMN")
	expectEqualInts(1, len(newNode.ChildIDs), t)
	expectEqualInts(40, len(newNode.ChildIDs[0]), t)
	expectEqualInts(1, len(newNode.ParentIDs), t)
	expectEqualInts(40, len(newNode.ParentIDs[0]), t)
}

// func TestAddNode(T *testing.T) {
// 	go describe("AddNode", t)

// 	var wg sync.WaitGroup
// 	graph := CreateGraph()
// 	sampleChan := make(chan sample)

// 	for i, line := range 10 {
// 		wg.Add(1)
// 		newNode := CreateNode("", "")

//     go AddNode(graph, newNode)
// 	}

// 	go func() {
// 			wg.Wait()
// 			close(sampleChan)
// 	}()

// 	for s := range sampleChan {
// 		..
// 	}
// }
