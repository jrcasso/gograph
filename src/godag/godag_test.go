package godag

import (
	"encoding/hex"
	"fmt"
	"testing"
)

// TestCreateGraph tests godag operations
func TestCreateGraph(t *testing.T) {

	graph := CreateGraph()
	rootNodeID, err := hex.DecodeString(graph.RootNode.ID)
	if err != nil {
		t.Error("Creating a graph should create a single root node.")
	}
	if len(graph.Nodes) != 1 {
		t.Error("Creating a graph should create a single root node.")
		fmt.Printf("Found %d\n", len(graph.Nodes))
		return
	}
	if len(rootNodeID) != 40 {
		t.Error("Creating a graph should create a single root node with a valid ID.")
		fmt.Printf("Found %d characters for %s\n", len(rootNodeID), rootNodeID)
		return
	}

}
