package snakecube

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Pos [3]int

type Solution struct {
	Sequence   []int
	Direction  []int
	Path       [][][]int
	StartPos   Pos
	Palindrome bool
}

type SolverState struct {
	N             int
	L             int
	sequenceIn    []int
	sequence      []int
	direction     []int
	path          [][][]int
	StartPos      Pos
	SolutionStore []Solution
	verbose       bool
	debug         bool
}

type SolutionbyStartPos struct {
	StartPos  Pos
	solutions []Solution
}

type SnakeSolution struct {
	StartPos  Pos       `json:"startPos"`
	Direction []int     `json:"direction"`
	Path      [][][]int `json:"path"`
}

func (state *SolverState) Init(N int, sequenceIn []int, verbose bool, debug bool) {

	state.verbose = verbose
	state.debug = debug

	if N < 3 || N > 4 {
		log.Fatal("N must be between 3 and 5")
	}
	state.N = N
	state.L = N * N * N

	nSeqIn := len(sequenceIn)
	if nSeqIn > state.L {
		log.Fatal("len(sequenceIn must be <= N**3")
	}
	state.sequenceIn = sequenceIn

	if state.verbose {
		fmt.Println("init state")
		fmt.Printf("size = %dx%dx%d\n", N, N, N)
		fmt.Printf("N = %d\n", state.N)
		fmt.Printf("L = %d\n", state.L)
		fmt.Printf("input sequence = %v\n", sequenceIn)
	}

}

func (state *SolverState) SearchFromPos(StartPos Pos) {

	var i int

	if state.verbose {
		fmt.Println(" ")
		defer track(runningtime("SearchFromPos from StartPos=" + fmt.Sprintf("%v", StartPos)))
	}

	state.StartPos = StartPos

	// reset sequence
	state.sequence = state.buildSequenceInit(true)
	if i < len(state.sequenceIn) {
		state.sequence[i] = state.sequenceIn[i]
	}

	// reset direction
	state.direction = state.buildDirectionInit(true)

	// reset path
	state.path = state.buildPathInit(true)

	state.path[StartPos[0]][StartPos[1]][StartPos[2]] = 1
	state.step(0, StartPos, 1, 1)

}

func (state *SolverState) SearchFromAllPos() {

	var i, j, k int

	for i = 0; i < state.N; i++ {
		for j = 0; j < state.N; j++ {
			for k = 0; k < state.N; k++ {
				StartPos := Pos{i, j, k}
				state.SearchFromPos(StartPos)
			}
		}
	}
}

func (state *SolverState) step(n int, pos Pos, direct int, exploredDim int) {

	sign := signInt(direct)

	newPos := pos
	newPos[abs(direct)-1] += sign

	if state.isValidPosition(newPos) {

		state.path[newPos[0]][newPos[1]][newPos[2]] = n + 2
		state.direction[n] = direct

		state.branch(n+1, newPos, direct, exploredDim)

		state.path[newPos[0]][newPos[1]][newPos[2]] = 0
		state.direction[n] = 0
	}
}

func (state *SolverState) isValidPosition(pos Pos) bool {

	if pos[0] >= 0 && pos[0] < 3 && pos[1] >= 0 && pos[1] < 3 && pos[2] >= 0 && pos[2] < 3 {
		if state.path[pos[0]][pos[1]][pos[2]] == 0 {
			return true
		}
		return false
	}
	return false

}

func (state *SolverState) branch(n int, pos Pos, direct int, exploredDim int) {

	var k int

	if n == state.L-1 {
		// path is complete
		if state.debug {
			fmt.Printf("path complete = %v\n", state.path)
		}
		state.sequence[n] = 0
		isLexicographicallySmallerOrEqual, isPalindrome := state.checkSolution()

		if isLexicographicallySmallerOrEqual {

			sequence := buildLexicographicSmallerSequence(state.sequence)
			direction := copySlice(state.direction)
			path := state.copyPath(state.path)
			startPos := state.StartPos

			solution := Solution{
				Sequence:   sequence,
				Direction:  direction,
				Path:       path,
				StartPos:   startPos,
				Palindrome: isPalindrome,
			}
			state.SolutionStore = append(state.SolutionStore, solution)

			if state.debug {
				fmt.Printf("==> solution = %v\n", solution)
			}
			// nSol := len(state.SolutionStore)
			// if nSol%1000 == 0 {
			// 	fmt.Printf("%d ", nSol)
			// }
		}
	} else {
		if n >= len(state.sequenceIn) || state.sequenceIn[n] == 0 {
			// go straight
			if state.debug {
				fmt.Printf("n=%d | go straight -> %v\n", n, state.path)
			}
			state.sequence[n] = 0
			state.step(n, pos, direct, exploredDim)
		}
		if n >= len(state.sequenceIn) || state.sequenceIn[n] == 1 {
			// make turn
			if state.debug {
				fmt.Printf("n=%d | make turn -> %v\n", n, state.path)
			}
			state.sequence[n] = 1
			for k = 1; k <= min(exploredDim, 3); k++ {
				if k != abs(direct) {
					state.step(n, pos, +k, exploredDim)
					state.step(n, pos, -k, exploredDim)

				}
			}
			if exploredDim < 3 {
				// move up one dimension
				state.step(n, pos, exploredDim+1, exploredDim+1)
			}
		}
	}
}

func (state *SolverState) checkSolution() (bool, bool) {
	var i int

	// palindrome
	isPalindrome := true
	for i = 0; i < abs(state.L/2)+1; i++ {
		if state.sequence[i] != state.sequence[state.L-1-i] {
			isPalindrome = false
			break
		}
	}

	// order
	revDirection := state.buildReverseDirection()
	order := compareLexicographicOrder(state.direction, revDirection)
	var isLexicographicallySmallerOrEqual bool
	if order <= 0 {
		isLexicographicallySmallerOrEqual = true
	} else {
		isLexicographicallySmallerOrEqual = false

	}

	return isLexicographicallySmallerOrEqual, isPalindrome

}

func (state *SolverState) buildReverseDirection() []int {
	var i, d, k int

	directionMap := make([]int, 3+1+3)
	exploredDim := make([]int, 3)
	for i = 0; i < len(directionMap); i++ {
		directionMap[i] = 0
	}
	for i = 0; i < len(exploredDim); i++ {
		exploredDim[i] = 0
	}

	plainReverseDirection := make([]int, state.L-1)
	for i = 0; i < len(plainReverseDirection); i++ {
		plainReverseDirection[i] = state.direction[state.L-2-i]
	}

	i = 0
	d = 0
	for d < 3 {
		k = abs(plainReverseDirection[i])
		if exploredDim[k-1] == 0 {
			exploredDim[k-1] = 1
			d += 1
			directionMap[3+plainReverseDirection[i]] = d
			directionMap[3-plainReverseDirection[i]] = -d
		}
		i += 1
	}

	reverseDirection := make([]int, state.L-1)
	for i = 0; i < state.L-1; i++ {
		reverseDirection[i] = directionMap[3+plainReverseDirection[i]]
	}
	return reverseDirection

}

func (state *SolverState) ShowSolutions(nSolShowMax int) {
	var i int

	mps := make(map[Pos][]Solution, 0)
	for _, e := range state.SolutionStore {
		handle := mps[e.StartPos]
		handle = append(handle, e)
		mps[e.StartPos] = handle
	}

	for StartPos, arr := range mps {

		nSol := len(arr)
		if nSol > 0 {

			nSolShow := min(nSol, nSolShowMax)
			if nSolShow <= 0 {
				nSolShow = nSol
			}

			fmt.Printf("\nshow %d/%d solutions from %v\n", nSolShow, nSol, StartPos)
			for i = 0; i < nSolShow; i++ {
				fmt.Printf("solution %d: %v\n", i, mps[StartPos][i].Path)
			}
		}
	}
}

func RunSequential(save bool) map[string][]SnakeSolution {

	var i, j, k int
	N := 3

	state := SolverState{}
	state.Init(3, []int{}, false, false)

	fmt.Printf("start RunSequential\n...")
	t0 := time.Now()

	for i = 0; i < N; i++ {
		for j = 0; j < N; j++ {
			for k = 0; k < N; k++ {
				StartPos := Pos{i, j, k}
				state.SearchFromPos(StartPos)
			}
		}
	}

	time1 := time.Since(t0)
	solutions := buildSnakeSolutions(state.SolutionStore)
	t1 := time.Now()

	time2 := time.Since(t1)
	nSeq := len(solutions)
	nSol := 0
	for _, v := range solutions {
		nSol += len(v)
	}

	fmt.Printf("\nsearch time: %s\nshape time: %s\nnb sequences: %d\nnb solutions: %d\n", time1, time2, nSeq, nSol)

	if save {
		SaveSolutions(solutions, "solutions.json")
	}

	return solutions
}

func buildSnakeSolutions(solutionsIn []Solution) map[string][]SnakeSolution {
	mSeqSol := make(map[string][]SnakeSolution, 0)

	for _, e := range solutionsIn {
		seq := arrayToString(e.Sequence, "")
		if _, ok := mSeqSol[seq]; !ok {
			mSeqSol[seq] = []SnakeSolution{}
		}

		sol := SnakeSolution{
			StartPos:  e.StartPos,
			Direction: e.Direction,
			Path:      e.Path,
		}
		mSeqSol[seq] = append(mSeqSol[seq], sol)
	}
	return mSeqSol

}

func runOneStartPos(StartPos Pos, wg *sync.WaitGroup, c chan SolutionbyStartPos) {

	defer wg.Done()

	state := SolverState{}
	state.Init(3, []int{}, false, false)
	state.SearchFromPos(StartPos)

	res := SolutionbyStartPos{
		StartPos:  StartPos,
		solutions: state.SolutionStore,
	}
	c <- res

}

func RunParallel(save bool) map[string][]SnakeSolution {

	var i, j, k int

	fmt.Printf("start RunParallel\n...")
	t0 := time.Now()

	N := 3
	arrStartPos := make([]Pos, 0)

	for i = 0; i < N; i++ {
		for j = 0; j < N; j++ {
			for k = 0; k < N; k++ {
				arrStartPos = append(arrStartPos, [3]int{i, j, k})
			}
		}
	}

	wg := sync.WaitGroup{}
	c := make(chan SolutionbyStartPos)

	for _, StartPos := range arrStartPos {
		wg.Add(1)
		go runOneStartPos(StartPos, &wg, c)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	_solutions := make([]Solution, 0)

	for response := range c {

		StartPos := response.StartPos
		NewSolutions := response.solutions
		nSol := len(NewSolutions)

		if nSol > 0 {
			fmt.Printf("\n==> %d solutions for %v", nSol, StartPos)
			_solutions = append(_solutions, NewSolutions...)
		}
	}

	time1 := time.Since(t0)
	t1 := time.Now()

	solutions := buildSnakeSolutions(_solutions)

	time2 := time.Since(t1)
	nSeq := len(solutions)
	nSol := 0
	for _, v := range solutions {
		nSol += len(v)
	}

	fmt.Printf("\nsearch time: %s\nshape time: %s\nnb sequences: %d\nnb solutions: %d\n", time1, time2, nSeq, nSol)

	if save {
		SaveSolutions(solutions, "solutions.json")
	}

	return solutions
}
