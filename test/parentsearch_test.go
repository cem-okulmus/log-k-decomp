package tests

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/cem-okulmus/BalancedGo/lib"
	logk "github.com/cem-okulmus/log-k-decomp/lib"
)

var EDGE int

// max returns the larger of two integers a and b
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

//getRandomEdge will produce a random Edge
func getRandomEdge(size int) lib.Edge {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	arity := r.Intn(size) + 1
	var vertices []int
	name := r.Intn(size*10) + EDGE + 1
	EDGE = name

	for i := 0; i < arity; i++ {

		vertices = append(vertices, r.Intn(size*10)+i+1)
	}

	return lib.Edge{Name: name, Vertices: vertices}
}

//getRandomGraph will produce a random Graph
func getRandomGraph(size int) (lib.Graph, map[string]int) {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	card := r.Intn(size) + 1

	var edges []lib.Edge
	var SpEdges []lib.Edges

	for i := 0; i < card; i++ {
		edges = append(edges, getRandomEdge(size))
	}

	outString := lib.Graph{Edges: lib.NewEdges(edges), Special: SpEdges}.ToHyberBenchFormat()
	parsedGraph, pGraph := lib.GetGraph(outString)

	return parsedGraph, pGraph.Encoding
}

//getRandomEdges will produce a random Edges struct
func getRandomEdges(size int) lib.Edges {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	k := r.Intn(size) + 1
	var edges []lib.Edge
	for j := 0; j < k; j++ {
		edges = append(edges, getRandomEdge(size*10))
	}

	return lib.NewEdges(edges)
}

func getRandomSep(g lib.Graph, size int) lib.Edges {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	k := r.Intn(size) + 1

	var selection []int
	for j := 0; j < k; j++ {
		selection = append(selection, r.Intn(g.Edges.Len()))
	}

	return lib.GetSubset(g.Edges, selection)
}

//TestSearchPar ensures that the parallel search for good parents always returns the same results,
// no matter how many splits are generated and run in parallel
func TestSearchPar(t *testing.T) {

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	randGraph, _ := getRandomGraph(30)
	k := r.Intn(5) + 1
	prevSep := getRandomSep(randGraph, 5)

	k = max(k, prevSep.Len())

	allowedParent := lib.FilterVertices(randGraph.Edges, prevSep.Vertices())

	combinParallel := lib.SplitCombin(allowedParent.Len(), k, runtime.GOMAXPROCS(-1), false)
	combinSeq := lib.SplitCombin(allowedParent.Len(), k, 1, false)

	parallelSearch := lib.ParallelSearch{
		H:          &randGraph,
		Edges:      &allowedParent,
		BalFactor:  2,
		Generators: combinParallel,
	}
	seqSearch := lib.ParallelSearch{
		H:          &randGraph,
		Edges:      &allowedParent,
		BalFactor:  2,
		Generators: combinSeq,
	}

	conn := lib.Inter(prevSep.Vertices(), randGraph.Vertices())
	predPar := logk.ParentCheck{Conn: conn, Child: prevSep.Vertices()}

	var allSepsSeq []lib.Edges
	var allSepsPar []lib.Edges

	fmt.Println("Starting par Search:")

	parallelSearch.FindNext(predPar)

	for ; !parallelSearch.ExhaustedSearch; parallelSearch.FindNext(predPar) {

		sep := lib.GetSubset(allowedParent, parallelSearch.Result)
		allSepsPar = append(allSepsPar, sep)

		fmt.Println("From indices", parallelSearch.Result, " the sep ", sep, " created")
	}

	fmt.Println("Starting seq Search:")

	for seqSearch.FindNext(predPar); !seqSearch.ExhaustedSearch; seqSearch.FindNext(predPar) {

		sep := lib.GetSubset(allowedParent, seqSearch.Result)
		allSepsSeq = append(allSepsSeq, sep)

		fmt.Println("From indices", seqSearch.Result, " the sep ", sep, " created")
	}

OUTER:
	for i := range allSepsSeq {
		sep := allSepsSeq[i]

		for j := range allSepsPar {
			other := allSepsPar[j]
			if other.Hash() == sep.Hash() {
				continue OUTER // found matching sep
			}
		}

		if len(allSepsSeq) != len(allSepsPar) {

			fmt.Println("Graph", randGraph)
			fmt.Println("prevSep", prevSep)
			fmt.Println("k: ", k)
			fmt.Println("Conn: ", lib.PrintVertices(conn))

			combinParallel2 := lib.SplitCombin(allowedParent.Len(), k, runtime.GOMAXPROCS(-1), false)
			combinSeq2 := lib.SplitCombin(allowedParent.Len(), k, 1, false)

			fmt.Print("\n All stuff in combinPar: ")
			for _, combin := range combinParallel2 {

				fmt.Print("\n")

				for combin.HasNext() {
					j := combin.GetNext()
					fmt.Print(j)
					combin.Confirm()
				}

				fmt.Print("\n")

			}

			fmt.Print("\n\n All stuff in combinSeq2: ")
			for _, combin := range combinSeq2 {

				for combin.HasNext() {
					j := combin.GetNext()
					fmt.Print(j)
					combin.Confirm()
				}

			}

			fmt.Println("\n\n Number of splits in parallel: ", len(combinParallel))

			fmt.Println("Allowed, ", allowedParent)

			fmt.Println("Seps found by seq, ", allSepsSeq)
			fmt.Println("Seps found by par, ", allSepsPar)

		}

		t.Errorf("Mismatch in returned seps between sequential and parallel Search, numSepsSeq %v, numSepsPar %v", len(allSepsSeq), len(allSepsPar))

	}

}
