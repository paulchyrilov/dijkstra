package dijkstra

import "sort"

//Vertex is a single node in the network, contains it's ID, best Distance (to
// itself from the src) and the weight to go to each other connected node (Vertex)
type Vertex struct {
	//ID of the Vertex
	ID int
	Name string
	//Best Distance to the Vertex
	Distance   int64
	bestVertex *Vertex
	bestArc *Arc
	//A set of all weights to the nodes in the map
	destinations map[int]Destination
}

type Destination struct {
	Vertex *Vertex
	Arcs   map[int]Arc
}

type Arc struct {
	Distance   int64
	Attributes interface{}
}

//NewVertex creates a new Vertex
func NewVertex(ID int) *Vertex {
	return &Vertex{ID: ID, destinations: map[int]Destination{}}
}

//AddVerticies adds the listed verticies to the graph, overwrites any existing
// Vertex with the same ID.
func (g *Graph) AddVerticies(verticies ...Vertex) {
	for _, v := range verticies {
		if v.ID >= len(g.Verticies) {
			newV := make([]Vertex, v.ID+1-len(g.Verticies))
			g.Verticies = append(g.Verticies, newV...)
		}
		g.Verticies[v.ID] = v
	}
}

//AddArc adds an arc to the Vertex, it's up to the user to make sure this is used
// correctly, firstly ensuring to use before adding to graph, or to use referenced
// of the Vertex instead of a copy. Secondly, to ensure the destination is a valid
// Vertex in the graph. Note that AddArc will overwrite any existing Distance set
// if there is already an arc set to Destination.
func (v *Vertex) AddArc(destinationVertex *Vertex, Distance int64, attributes interface{}) {
	if v.destinations == nil {
		v.destinations = map[int]Destination{}
	}

	if destination, ok := v.destinations[destinationVertex.ID]; ok {
		destination.Arcs[len(destination.Arcs)] = Arc{Distance:Distance, Attributes:attributes}
		sort.Slice(destination.Arcs, func(i, j int) bool {
			return destination.Arcs[i].Distance > destination.Arcs[j].Distance
		})
	} else {
		newDestination := Destination{Vertex: destinationVertex, Arcs: map[int]Arc{}}
		newDestination.Arcs[len(newDestination.Arcs)] = Arc{Distance:Distance, Attributes:attributes}
		v.destinations[destinationVertex.ID] = newDestination
	}
}

/*
I decided you don't get that kind of privelage
#checkyourprivelage
//RemoveArc completely removes the arc to Destination (if it exists)
func (v *Vertex) RemoveArc(Destination int) {
	delete(v.Arcs, Destination)
}*/

//GetDestination gets the specified arc to Destination, bool is false if no arc found
func (v *Vertex) GetDestination(DestinationId int) (destinations Destination, ok bool) {
	if v.destinations == nil {
		return Destination{}, false
	}
	//idk why but doesn't work on one line?
	destinations, ok = v.destinations[DestinationId]
	return
}

//GetDestination gets the specified arc to Destination, bool is false if no arc found
func (v *Vertex) GetDestinationByVertex(destinationVertex *Vertex) (destinations Destination, ok bool) {
	if v.destinations == nil {
		return Destination{}, false
	}
	//idk why but doesn't work on one line?
	destinations, ok = v.destinations[destinationVertex.ID]
	return
}
