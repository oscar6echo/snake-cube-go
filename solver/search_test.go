package snakecube

import (
	"testing"
)

func TestOneSnakeFullOneStartPos(t *testing.T) {

	sequenceIn := []int{0, 0, 1, 1, 0, 1, 1, 0, 1, 0, 1, 1, 0, 1, 0, 1, 1, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0}
	StartPos := [3]int{0, 0, 0}
	verbose := true
	debug := false

	solution_path := [][][]int{
		{{1, 24, 11}, {6, 25, 10}, {7, 8, 9}},
		{{2, 23, 12}, {5, 26, 13}, {16, 15, 14}},
		{{3, 22, 21}, {4, 27, 20}, {17, 18, 19}},
	}

	state := SolverState{}
	state.Init(3, sequenceIn, verbose, debug)

	state.SearchFromPos(StartPos)

	solutions := state.SolutionStore

	nSol := len(solutions)
	if nSol != 1 {
		t.Errorf("len(SolutionStore) = %d, want %d", nSol, 1)

		sol_found := solutions[0].Path
		sol_want := solution_path
		same := state.are_equal_path_object(sol_found, sol_want)
		if !same {
			t.Errorf("solution found = %v, want %v", sol_found, sol_want)
		}
	}
}

func TestOneSnakeFullAllStartPos(t *testing.T) {

	var i int

	sequenceIn := []int{0, 0, 1, 1, 0, 1, 1, 0, 1, 0, 1, 1, 0, 1, 0, 1, 1, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0}
	verbose := false
	debug := false

	solution_paths := [][][][]int{
		{
			{{1, 24, 11}, {6, 25, 10}, {7, 8, 9}},
			{{2, 23, 12}, {5, 26, 13}, {16, 15, 14}},
			{{3, 22, 21}, {4, 27, 20}, {17, 18, 19}},
		},
	}

	state := SolverState{}
	state.Init(3, sequenceIn, verbose, debug)
	state.SearchFromAllPos()

	solutions := state.SolutionStore

	nSol_found := len(solutions)
	nSol_want := len(solution_paths)

	if nSol_found != nSol_want {
		t.Errorf("len(SolutionStore) = %d, want %d", nSol_found, nSol_want)
	}

	for i = 0; i > min(nSol_found, nSol_want); i++ {

		t.Errorf("SolutionStore = %v", solutions)
		// t.Errorf("SolutionStore = %v", state.SolutionStore)

		sol_found := solutions[i].Path
		sol_want := solution_paths[i]
		same := state.are_equal_path_object(sol_found, sol_want)
		if !same {
			t.Errorf("solution %d found = %v, want %v", i, sol_found, sol_want)
		}
	}
}

func TestOneSnakePartialAllStartPos(t *testing.T) {

	var i int

	sequenceIn := []int{0, 0, 1, 1, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	verbose := false
	debug := false

	solution_paths := [][][][]int{
		{
			{{13, 12, 9}, {16, 1, 8}, {17, 18, 7}},
			{{14, 11, 10}, {15, 2, 27}, {20, 19, 6}},
			{{23, 24, 25}, {22, 3, 26}, {21, 4, 5}},
		},

		{
			{{25, 24, 9}, {26, 1, 10}, {27, 14, 13}},
			{{20, 23, 8}, {19, 2, 11}, {16, 15, 12}},
			{{21, 22, 7}, {18, 3, 6}, {17, 4, 5}},
		},
		{
			{{21, 22, 9}, {18, 1, 10}, {17, 14, 13}},
			{{20, 23, 8}, {19, 2, 11}, {16, 15, 12}},
			{{25, 24, 7}, {26, 3, 6}, {27, 4, 5}},
		},
		{
			{{25, 24, 21}, {26, 1, 20}, {27, 16, 17}},
			{{10, 23, 22}, {11, 2, 19}, {14, 15, 18}},
			{{9, 8, 7}, {12, 3, 6}, {13, 4, 5}},
		},
		{
			{{27, 12, 13}, {26, 1, 16}, {21, 20, 17}},
			{{10, 11, 14}, {25, 2, 15}, {22, 19, 18}},
			{{9, 8, 7}, {24, 3, 6}, {23, 4, 5}},
		},
	}

	state := SolverState{}
	state.Init(3, sequenceIn, verbose, debug)
	state.SearchFromAllPos()

	solutions := state.SolutionStore

	nSol_found := len(solutions)
	nSol_want := len(solution_paths)

	if nSol_found != nSol_want {
		t.Errorf("len(SolutionStore) = %d, want %d", nSol_found, nSol_want)
	}

	for i = 0; i > min(nSol_found, nSol_want); i++ {

		t.Errorf("SolutionStore = %v", solutions)
		// t.Errorf("SolutionStore = %v", state.SolutionStore)

		sol_found := solutions[i].Path
		sol_want := solution_paths[i]
		same := state.are_equal_path_object(sol_found, sol_want)
		if !same {
			t.Errorf("solution %d found = %v, want %v", i, sol_found, sol_want)
		}
	}
}

func TestSearchAllSequential(t *testing.T) {

	solutions := RunSequential(false)

	nSeq_found := len(solutions)
	nSol_found := 0
	for _, v := range solutions {
		nSol_found += len(v)
	}

	nSeq_want := 11487
	nSol_want := 51704

	if nSeq_found != nSeq_want {
		t.Errorf("len(SolutionStore) = %d, want %d", nSeq_found, nSeq_want)
	}

	if nSol_found != nSol_want {
		t.Errorf("len(SolutionStore) = %d, want %d", nSol_found, nSol_want)
	}

}

func TestSearchAllParallel(t *testing.T) {

	solutions := RunParallel(false)

	nSeq_found := len(solutions)
	nSol_found := 0
	for _, v := range solutions {
		nSol_found += len(v)
	}

	nSeq_want := 11487
	nSol_want := 51704

	if nSeq_found != nSeq_want {
		t.Errorf("len(SolutionStore) = %d, want %d", nSeq_found, nSeq_want)
	}

	if nSol_found != nSol_want {
		t.Errorf("len(SolutionStore) = %d, want %d", nSol_found, nSol_want)
	}

}
