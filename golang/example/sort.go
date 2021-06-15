package main

import (
	"fmt"
	"sort"
)

type Reverse struct {
	sort.Interface
}

func (r Reverse) Less(i, j int) bool {
	return r.Less(j, i)
}

func main() {

	intSlice := []int{1, 2, 7, 9, 3, 4}

	sort.Ints(intSlice)
	fmt.Println("1: ", intSlice)

	sort.Sort(sort.IntSlice(intSlice))

	fmt.Println("sort: ", intSlice)

}
