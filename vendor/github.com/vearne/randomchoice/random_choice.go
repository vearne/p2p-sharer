package randomchoice

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomChoice(count int, choiceCount int) []int {
	if count < choiceCount {
		choiceCount = count
	}

	intSlice := make([]int, count)
	for i := 0; i < count; i++ {
		intSlice[i] = i
	}

	idx := 0
	for i := 0; i < choiceCount; i++ {
		idx = rand.Int()%count + i
		// swap
		//swap(intSlice, i, idx)
		intSlice[i], intSlice[idx] = intSlice[idx], intSlice[i]
		count--

	}
	return intSlice[0:choiceCount]
}

func swap(s []int, i int, j int) {
	s[i], s[j] = s[j], s[i]
}
