package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraph(t *testing.T) {
	// Create a test graph
	graph := NewGraph()

	nodeA := &Node{Name: "A"}
	nodeB := &Node{Name: "B"}
	nodeC := &Node{Name: "C"}

	graph.AddNode(nodeA)
	graph.AddNode(nodeB)
	graph.AddNode(nodeC)

	graph.AddEdge(nodeA, &Adj{Node: nodeB})
	graph.AddEdge(nodeB, &Adj{Node: nodeC})

	t.Run("TestTopologicalSort", func(t *testing.T) {
		expectedOrder := []*Node{nodeA, nodeB, nodeC}
		actualOrder := graph.TopologicalSort()

		assert.Equal(t, expectedOrder, actualOrder)
	})

	t.Run("GetAdjacency", func(t *testing.T) {
		expectedAdjacency := []*Adj{{Node: nodeC}}
		actualAdjacency := graph.GetAdjacency(nodeB)

		assert.Equal(t, expectedAdjacency, actualAdjacency)
	})

	t.Run("TestGetLeafNodes", func(t *testing.T) {
		expectedLeafNodes := []*Node{nodeA}
		actualLeafNodes := graph.GetLeafNodes()

		assert.Equal(t, expectedLeafNodes, actualLeafNodes)
	})

	t.Run("TestGetNodesSortedByName", func(t *testing.T) {
		expectedSortedNodes := []*Node{nodeA, nodeB, nodeC}
		actualSortedNodes := graph.GetNodesSortedByName()

		assert.Equal(t, expectedSortedNodes, actualSortedNodes)
	})

	t.Run("TestGetAdjacenciesSortedByName", func(t *testing.T) {
		expectedAdjacencies := []*Adj{{Node: nodeB, Func: nil}}
		actualAdjacencies := graph.GetAdjacenciesSortedByName(nodeA)

		assert.Equal(t, expectedAdjacencies, actualAdjacencies)
	})
}
