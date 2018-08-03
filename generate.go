package dijkstra

import "math/rand"

//Generate generates file with the amount of nodes specified
func Generate(nodes int) Graph {
	//	fmt.Println("Generating file "+filename+" with nodes ", nodes)
	graph := Graph{}
	var i int
	for i = 0; i < nodes; i++ {
		v := NewVertex(i)
		for j := 0; j < nodes; j++ {
			if j == i {
				continue
			}
			var vertex, err = graph.GetVertex(j)
			if nil != err {
				vertex := NewVertex(j)
				v.AddArc(vertex, int64(2*nodes-j)+rand.Int63n(int64(nodes)*int64(nodes-j+1)), nil)
			} else {
				v.AddArc(vertex, int64(2*nodes-j)+rand.Int63n(int64(nodes)*int64(nodes-j+1)), nil)
			}
		}
		graph.AddVerticies(*v)
	}
	return graph
}
