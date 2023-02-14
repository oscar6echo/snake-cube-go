package snakecube

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func (state *SolverState) buildSequenceInit(reset bool) []int {
	var i int
	sequence := make([]int, state.L)

	if reset {
		for i = 0; i < state.L; i++ {
			sequence[i] = 0
		}
	}
	return sequence
}

func copySlice(src []int) []int {
	dst := make([]int, len(src))
	copy(dst, src)
	return dst
}

func (state *SolverState) buildDirectionInit(reset bool) []int {
	var i int
	direction := make([]int, state.L-1)

	if reset {
		for i = 0; i < state.L-1; i++ {
			direction[i] = 0
		}
	}
	return direction
}

func (state *SolverState) buildPathInit(reset bool) [][][]int {
	var i, j, k int
	path := make([][][]int, state.N)

	for i = 0; i < state.N; i++ {
		path[i] = make([][]int, state.N)
		for j = 0; j < state.N; j++ {
			path[i][j] = make([]int, state.N)
			for k = 0; k < state.N; k++ {
				if reset {
					path[i][j][k] = 0
				}
			}
		}
	}
	return path
}

func (state *SolverState) copyPath(src [][][]int) [][][]int {
	var i, j, k int
	dst := make([][][]int, state.N)

	for i = 0; i < state.N; i++ {
		dst[i] = make([][]int, state.N)
		for j = 0; j < state.N; j++ {
			dst[i][j] = make([]int, state.N)
			for k = 0; k < state.N; k++ {
				dst[i][j][k] = src[i][j][k]
			}
		}
	}
	return dst
}

func (state *SolverState) are_equal_path_object(a [][][]int, b [][][]int) bool {
	var i, j, k int
	same := true

	for i = 0; i < state.N; i++ {

		for j = 0; j < state.N; j++ {
			if !same {
				break
			}
			for k = 0; k < state.N; k++ {
				if !same {
					break
				}
				if a[i][j][k] != b[i][j][k] {
					same = false
					break
				}
			}
		}
	}
	return same
}

func compareLexicographicOrder(a []int, b []int) int {
	var i int

	if len(a) != len(b) {
		log.Fatal("UNEXPECTED")
	}

	for i = 0; i < len(a); i++ {
		if a[i] < b[i] {
			return -1
		} else if a[i] > b[i] {
			return +1
		}
	}
	return 0

}

func buildLexicographicSmallerSequence(sequence []int) []int {

	var out []int
	rev_sequence := reverseSlice(sequence)

	comp := compareLexicographicOrder(sequence, rev_sequence)
	if comp <= 0 {
		out = copySlice(sequence)
	} else {
		out = rev_sequence
	}
	return out

}

func SaveSolutions(solutions map[string][]SnakeSolution, filename string) {
	// func SaveSolutions(solutions []Solution, filename string) {

	t0 := time.Now()

	obj, _ := json.Marshal(solutions)
	_ = ioutil.WriteFile(filename, obj, 0644)

	time := time.Since(t0)
	fmt.Printf("solutions saved to: %v - done in %s\n", filename, time)

}
