package gograph

import (
	"testing"
)

func expectEqualInts(value int, expectation int, t *testing.T) {
	if value != expectation {
		t.Errorf("Failed: expected %d, but found %d", expectation, value)
	}
}

func expectEqualStrings(value string, expectation string, t *testing.T) {
	if value != expectation {
		t.Errorf("Failed: expected %s, but found %s", expectation, value)
	}
}

func describe(functionName string, t *testing.T) {
	t.Logf("%s:", functionName)
}

func it(description string, t *testing.T) {
	t.Logf("    - it %s", description)
}

func context(description string, t *testing.T) {
	t.Logf(" - when %s", description)
}

func CreateDirectedGraphHelper() DirectedGraph {
	var node *DirectedNode
	var graphDirectedNodes []*DirectedNode
	return DirectedGraph{
		DirectedNodes:    graphDirectedNodes,
		RootDirectedNode: node,
		NodeMap:          make(map[string]*DirectedNode),
	}
}
