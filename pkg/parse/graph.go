package parse

import (
	"go/types"
	"sort"

	"github.com/emilien-puget/go-dependency-graph/pkg/parse/struct_decl"
	"golang.org/x/tools/go/packages"
)

// Node represents a node of the dependency graph.
type Node struct {
	Name            string // The fully qualified name of the struct, PackageName.StructName
	PackageName     string
	StructName      string
	Methods         []struct_decl.Method
	Doc             string
	External        bool
	InboundEdges    []*Node
	ActualNamedType *types.Named
	P               *packages.Package
	FilePath        string
}

func (n *Node) MergeAdditionalFields(other *Node) {
	if len(other.Methods) > 0 {
		n.Methods = other.Methods
	}
	if n.ActualNamedType == nil && other.ActualNamedType != nil {
		n.ActualNamedType = other.ActualNamedType
	}
	if n.P == nil && other.P != nil {
		n.P = other.P
	}
	if n.FilePath == "" && other.FilePath != "" {
		n.FilePath = other.FilePath
	}
}

// Graph represents the dependency graph.
type Graph struct {
	Nodes          []*Node            // A slice of pointers to all the nodes in the graph
	Adj            map[*Node][]*Adj   // An adjacency list mapping each node to its adjacent nodes
	NodeByName     map[string]*Node   // A hash map to store nodes by their names
	NodesByPackage map[string][]*Node // A hash map to store nodes by their packages
}

type Adj struct {
	Node *Node
	Func []string
}

func NewGraph() *Graph {
	return &Graph{
		Nodes:          make([]*Node, 0),
		Adj:            make(map[*Node][]*Adj),
		NodeByName:     make(map[string]*Node),
		NodesByPackage: make(map[string][]*Node),
	}
}

// AddNode adds a new Node to the graph.
func (g *Graph) AddNode(node *Node) {
	existingNode, exists := g.NodeByName[node.Name]
	if exists {
		existingNode.MergeAdditionalFields(node)
		return
	}

	g.Nodes = append(g.Nodes, node)
	g.Adj[node] = make([]*Adj, 0)  // Initialize the adjacency list for the new node
	g.NodeByName[node.Name] = node // Add the node to the hash map using its name as the key

	// Update the NodesByPackage hash map
	g.NodesByPackage[node.PackageName] = append(g.NodesByPackage[node.PackageName], node)

	// Initialize the inbound edges slice for the new node
	node.InboundEdges = make([]*Node, 0)
}

func (g *Graph) GetAdjacency(node *Node) []*Adj {
	return g.Adj[node]
}

// GetNodeByName finds a node by its name in the graph.
func (g *Graph) GetNodeByName(name string) *Node {
	return g.NodeByName[name]
}

// AddEdge adds a directed edge between two nodes in the graph.
func (g *Graph) AddEdge(from *Node, to *Adj) {
	existingFromNode, exists := g.NodeByName[from.Name]
	if exists {
		from = existingFromNode
	}

	existingToNode, exists := g.NodeByName[to.Node.Name]
	if exists {
		existingToNode.MergeAdditionalFields(to.Node)
		to.Node = existingToNode
	}

	// Update the adjacency list
	g.Adj[from] = append(g.Adj[from], to)

	// Update the inbound edges for the 'to' node
	to.Node.InboundEdges = append(to.Node.InboundEdges, from)
}

// TopologicalSort performs a topological sort on the graph.
func (g *Graph) TopologicalSort() []*Node {
	visited := make(map[*Node]bool)
	stack := make([]*Node, 0)

	var dfs func(node *Node)
	dfs = func(node *Node) {
		visited[node] = true

		for _, neighbor := range g.Adj[node] {
			if !visited[neighbor.Node] {
				dfs(neighbor.Node)
			}
		}

		stack = append(stack, node) // Push the node onto the stack
	}

	for _, node := range g.Nodes {
		if !visited[node] {
			dfs(node)
		}
	}

	// Pop nodes from the stack to get the topological order
	sortedNodes := make([]*Node, 0, len(stack))
	for i := len(stack) - 1; i >= 0; i-- {
		sortedNodes = append(sortedNodes, stack[i])
	}

	return sortedNodes
}

// GetLeafNodes returns all the leaf nodes.
// A leaf node is a node without any inbound nodes.
func (g *Graph) GetLeafNodes() []*Node {
	leafNodes := make([]*Node, 0)

	for _, node := range g.Nodes {
		if len(node.InboundEdges) == 0 && len(g.Adj[node]) > 0 {
			leafNodes = append(leafNodes, node)
		}
	}

	return leafNodes
}

// GetNodesSortedByName returns all nodes sorted by node name.
func (g *Graph) GetNodesSortedByName() []*Node {
	sortedNodes := make([]*Node, len(g.Nodes))
	copy(sortedNodes, g.Nodes)
	sort.Slice(sortedNodes, func(i, j int) bool {
		return sortedNodes[i].Name < sortedNodes[j].Name
	})
	return sortedNodes
}

// GetAdjacenciesSortedByName returns the adjacencies sorted by the adjacent node names.
func (g *Graph) GetAdjacenciesSortedByName(node *Node) []*Adj {
	adjacencies := g.Adj[node]
	sortedAdjacencies := make([]*Adj, len(adjacencies))
	copy(sortedAdjacencies, adjacencies)
	sort.Slice(sortedAdjacencies, func(i, j int) bool {
		return sortedAdjacencies[i].Node.Name < sortedAdjacencies[j].Node.Name
	})
	return sortedAdjacencies
}
