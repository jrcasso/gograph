package godag

import (
	"testing"
)

func expectEqualInts(value int, expectation int, t *testing.T) {
	if value != expectation {
		t.Errorf("Failed: expected %d, but found %d", value, expectation)
		return
	}
}

func expectEqualStrings(value string, expectation string, t *testing.T) {
	if value != expectation {
		t.Errorf("Failed: expected %s, but found %s", value, expectation)
		return
	}
}

func describe(functionName string, t *testing.T) {
	t.Logf("%s:", functionName)
}

func it(description string, t *testing.T) {
	t.Logf(" - it %s", description)
}

func when(description string, t *testing.T) {
	t.Logf("      when %s", description)
}
