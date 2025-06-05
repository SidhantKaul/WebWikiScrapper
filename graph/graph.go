package graph

import (
	"LinkScrapper/jsonwriter"
)

var Width int
var Depth int

// Graph represents a node in a graph structure.
type Graph struct {
	URL   string
	index int
	List  []*Graph
}

// SetGraphDimensions sets the global width and depth used when creating graphs.
func SetGraphDimensions(pWidth int, pDepth int) {
	Width = pWidth
	Depth = pDepth
}

// NewGraph creates a new Graph node with the given URL and allocates space for children.
func NewGraph(pURL string) *Graph {
	return &Graph{
		URL:   pURL,
		index: 0,
		List:  make([]*Graph, Width),
	}
}

// AddChild adds a child graph node with the given URL. Returns false if the child limit is reached.
func (graph *Graph) AddChild(URL string) bool {
	if graph.index == Width {
		return false

	}

	child := NewGraph(URL)
	graph.List[graph.index] = child

	graph.index++

	return true
}

// GetWidth returns the number of children currently added to the node.
func (graph *Graph) GetWidth() int {
	return graph.index
}

// GraphToJson writes the graph's structure to the specified JSON file.
func (graph *Graph) GraphToJson(filePath string) {
	writer := jsonwriter.NewJsonWriter(filePath)
	graph.PrintNodeInfo(writer)
	writer.CloseFile()
}

// PrintNodeInfo recursively writes the graph node and its children in JSON format.
func (graph *Graph) PrintNodeInfo(writer *jsonwriter.JsonWriter) {
	if graph == nil {
		return
	}

	writer.StartObject()
	writer.AddValue("Link", graph.URL)
	writer.StartArray("Children")

	for i := 0; i < graph.GetWidth(); i++ {
		if graph.List[i] != nil {
			graph.List[i].PrintNodeInfo(writer)
		}
	}

	writer.EndArray()
	writer.EndObject()
}
