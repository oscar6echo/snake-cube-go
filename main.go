package main

import (
	solver "snakecube/solver"
	stats "snakecube/stats"
)

func main() {

	choose := 3

	if choose == 1 {

		sequenceIn := []int{0, 0, 1, 1, 0, 1, 1, 0, 1, 0, 1, 1, 0, 1, 0, 1, 1, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0}

		sequenceIn = sequenceIn[:0]

		verbose := false
		debug := false

		state := solver.SolverState{}
		state.Init(3, sequenceIn, verbose, debug)

		StartPos := [3]int{0, 0, 0}
		state.SearchFromPos(StartPos)

		state.ShowSolutions(10)

	} else if choose == 2 {

		solver.RunSequential(true)

	} else if choose == 3 {

		solver.RunParallel(true)

	} else if choose == 5 {

		s := stats.Solutions{}

		s.Load("solutions.json")

	}

}
