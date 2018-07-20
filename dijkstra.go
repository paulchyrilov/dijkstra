package dijkstra

import "math"

//Shortest calculates the shortest path from src to dest
func (g *Graph) Shortest(src, dest int) (BestPath, error) {
	return g.evaluate(src, dest, true)
}

//Longest calculates the longest path from src to dest
func (g *Graph) Longest(src, dest int) (BestPath, error) {
	return g.evaluate(src, dest, false)
}

func (g *Graph) finally(src, dest int) (BestPath, error) {
	if !g.visitedDest {
		return BestPath{}, ErrNoPath
	}
	return g.bestPath(src, dest), nil
}

func (g *Graph) setup(shortest bool, src int, list int) {
	//-1 auto list
	//Get a new list regardless
	if list >= 0 {
		g.forceList(list)
	} else if shortest {
		g.forceList(-1)
	} else {
		g.forceList(-2)
	}
	//Reset state
	g.visitedDest = false
	//Reset the best current value (worst so it gets overwritten)
	// and set the defaults *almost* as bad
	// set all best verticies to -1 (unused)
	if shortest {
		g.setDefaults(int64(math.MaxInt64)-2)
		g.best = int64(math.MaxInt64)
	} else {
		g.setDefaults(int64(math.MinInt64)+2)
		g.best = int64(math.MinInt64)
	}
	//Set the Distance of initial Vertex 0
	g.Verticies[src].Distance = 0
	//Add the source Vertex to the list
	g.visiting.PushOrdered(&g.Verticies[src])
}

func (g *Graph) forceList(i int) {
	//-2 long auto
	//-1 short auto
	//0 short pq
	//1 long pq
	//2 short ll
	//3 long ll
	switch i {
	case -2:
		if len(g.Verticies) < 800 {
			g.forceList(2)
		} else {
			g.forceList(0)
		}
		break
	case -1:
		if len(g.Verticies) < 800 {
			g.forceList(3)
		} else {
			g.forceList(1)
		}
		break
	case 0:
		g.visiting = priorityQueueNewShort()
		break
	case 1:
		g.visiting = priorityQueueNewLong()
		break
	case 2:
		g.visiting = linkedListNewShort()
		break
	case 3:
		g.visiting = linkedListNewLong()
		break
	default:
		panic(i)
	}
}

func (g *Graph) bestPath(src, dest int) BestPath {
	var path []Vertex
	for c := g.Verticies[dest]; c.ID != src; c = g.Verticies[c.bestVertex.ID] {
		path = append(path, c)
	}
	path = append(path, g.Verticies[src])
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return BestPath{g.Verticies[dest].Distance, path}
}

func (g *Graph) evaluate(src, dest int, shortest bool) (BestPath, error) {
	//Setup graph
	g.setup(shortest, src, -1)
	return g.postSetupEvaluate(src, dest, shortest)
}

func (g *Graph) postSetupEvaluate(src, dest int, shortest bool) (BestPath, error) {
	var current *Vertex
	oldCurrent := -1
	for g.visiting.Len() > 0 {
		//Visit the current lowest distanced Vertex
		//TODO WTF
		current = g.visiting.PopOrdered()
		if oldCurrent == current.ID {
			continue
		}
		oldCurrent = current.ID
		/*
			if shortest {
				current = heap.Pop(g.visiting).(*Vertex)
			} else {
				current = heap.Pop(g.visiting).(*Vertex)
			}*/
		//If we have hit the destination set the flag, cheaper than checking it's
		// Distance change at the end
		if current.ID == dest {
			g.visitedDest = true
			continue
		}
		//If the current Distance is already worse than the best try another Vertex
		if shortest && current.Distance >= g.best { //} || (!shortest && current.Distance <= g.best) {
			continue
		}
		for v, destination := range current.destinations {
			for _, arc := range destination.Arcs {
				//If the Arcs has better access, than the current best, update the Vertex being touched
				if (shortest && current.Distance+arc.Distance < g.Verticies[v].Distance) ||
					(!shortest && current.Distance+arc.Distance > g.Verticies[v].Distance) {
					//if g.Verticies[v].bestVertex == current.ID && g.Verticies[v].ID != dest {
					if current.bestVertex.ID == v && g.Verticies[v].ID != dest {
						//also only do this if we aren't checkout out the best Distance again
						//This seems familiar 8^)
						return BestPath{}, newErrLoop(current.ID, v)
					}
					g.Verticies[v].Distance = current.Distance + arc.Distance
					g.Verticies[v].bestVertex = current
					if v == dest {
						//If this is the destination update best, so we caInfinite loop detectedn stop looking at
						// useless Verticies
						g.best = current.Distance + arc.Distance
					}
					//Push this updated Vertex into the list to be evaluated, pushes in
					// sorted form
					g.visiting.PushOrdered(&g.Verticies[v])
				}
			}
		}
	}
	return g.finally(src, dest)
}

//BestPath contains the solution of the most optimal path
type BestPath struct {
	Distance int64
	Path     []Vertex
}
