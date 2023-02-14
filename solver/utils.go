package snakecube

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func runningtime(s string) (string, time.Time) {
	log.Println("Start:	", s)
	return s, time.Now()
}

func track(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("End:	", s, "took", endTime.Sub(startTime))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func uniqueStr(arr []string) []string {
	occured := map[string]bool{}
	result := []string{}
	for _, e := range arr {
		if occured[e] != true {
			occured[e] = true
			result = append(result, e)
		}
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func signInt(a int) int {
	if a < 0 {
		return -1
	} else if a > 0 {
		return +1
	}
	return 0
}

func reverseSlice(a []int) []int {
	var i int

	b := make([]int, len(a))
	for i = 0; i < len(a); i++ {
		b[i] = a[len(a)-1-i]
	}

	return b
}

func arrayToString(a []int, sep string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", sep, -1), "[]")
}
