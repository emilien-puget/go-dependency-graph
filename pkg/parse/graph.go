package parse

// Node represents a node of the dependency graph
type Node struct {
	Name        string // The name of the struct
	PackageName string
	StructName  string
	Methods     []string
	Doc         string
	External    bool
}

// Graph represents the dependency graph
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

// AddNode adds a new Node to the graph
func (g *Graph) AddNode(node *Node) {
	g.Nodes = append(g.Nodes, node)
	g.Adj[node] = make([]*Adj, 0)  // Initialize the adjacency list for the new node
	g.NodeByName[node.Name] = node // Add the node to the hash map using its name as the key

	// Add the node to the NodesByPackage hash map
	g.NodesByPackage[node.PackageName] = append(g.NodesByPackage[node.PackageName], node)
}

func (g *Graph) GetAdjacency(node *Node) []*Adj {
	return g.Adj[node]
}

// GetNodeByName finds a node by its name in the graph
func (g *Graph) GetNodeByName(name string) *Node {
	return g.NodeByName[name]
}

// AddEdge adds a directed edge between two nodes in the graph
func (g *Graph) AddEdge(from *Node, to *Adj) {
	g.Adj[from] = append(g.Adj[from], to)
}

// TopologicalSort performs a topological sort on the graph
func (g *Graph) TopologicalSort() []*Node {
	visited := make(map[*Node]bool)
	sortedNodes := make([]*Node, 0)

	var dfs func(node *Node)
	dfs = func(node *Node) {
		visited[node] = true

		for _, neighbor := range g.Adj[node] {
			if !visited[neighbor.Node] {
				dfs(neighbor.Node)
			}
		}

		sortedNodes = append(sortedNodes, node)
	}

	for _, node := range g.Nodes {
		if !visited[node] {
			dfs(node)
		}
	}

	// Reverse the sortedNodes to get the topological order
	for i, j := 0, len(sortedNodes)-1; i < j; i, j = i+1, j-1 {
		sortedNodes[i], sortedNodes[j] = sortedNodes[j], sortedNodes[i]
	}

	return sortedNodes
}
