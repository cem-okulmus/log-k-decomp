package lib

import (
	"github.com/cem-okulmus/BalancedGo/lib"
	"github.com/cem-okulmus/disjoint"
)

// ParentCheck looks a separator that could function as the direct ancestor (or "parent")
// of some child node in the GHD, where the connecting vertices "Conn" are explicitly provided
type ParentCheck struct {
	Conn  []int
	Child []int
}

// Check performs the needed computation to ensure whether sep is a good parent
func (p ParentCheck) Check(H *lib.Graph, sep *lib.Edges, balFactor int, Vertices map[int]*disjoint.Element) bool {

	//balancedness condition
	comps, _, _ := H.GetComponents(*sep, Vertices)

	foundCompLow := false
	var compLow lib.Graph

	balancednessLimit := (((H.Len()) * (balFactor - 1)) / balFactor)

	for i := range comps {
		if comps[i].Len() > balancednessLimit {
			foundCompLow = true
			compLow = comps[i]
		}
	}

	if !foundCompLow {
		return false // a bad parent :(
	}

	vertCompLow := compLow.Vertices()
	childχ := lib.Inter(p.Child, vertCompLow)

	if !lib.Subset(lib.Inter(vertCompLow, p.Conn), sep.Vertices()) {
		return false // also a bad parent :(
	}

	// Connectivity check
	if !lib.Subset(lib.Inter(vertCompLow, sep.Vertices()), childχ) {
		return false // again a bad parent :( Calling child services ...
	}

	return true // found a good parent :)
}
