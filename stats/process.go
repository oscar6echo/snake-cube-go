package stats

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

type Solutions []Solution

func (state *Solutions) Load(pathFile string) {

	file, err := os.Open(pathFile)
	if err != nil {
		log.Fatalf("missing file %s", pathFile)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)
	var solutions Solutions

	t0 := time.Now()

	json.Unmarshal([]byte(byteValue), &solutions)

	time := time.Since(t0)
	fmt.Printf("loaded file %s- done in %s\n", pathFile, time)

	nSol := len(solutions)
	fmt.Printf("%d solutions\n", nSol)

}
