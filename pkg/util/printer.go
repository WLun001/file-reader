package util

import (
	"fmt"
	"sort"
	"time"
)

// logic adopted from https://stackoverflow.com/questions/18695346/how-to-sort-a-mapstringint-by-its-values
func PrintTop5(m map[string]int) {
	n := map[int][]string{}
	var a []int
	for k, v := range m {
		n[v] = append(n[v], k)
	}
	for k := range n {
		a = append(a, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	count := 0
	fmt.Printf("top 5 words\n\n")
	for _, k := range a {
		for _, s := range n[k] {
			if count > 5 {
				break
			}
			fmt.Printf("%s, %d\n", s, k)
			count++
		}
	}
}

func PrintResult(filename string,size int64, dict map[string]int, words int, start time.Time)  {
	fmt.Printf("%s is %d bytes\n", filename, size)
	fmt.Printf("uniqueWords: %d, wordCount: %d\n", len(dict), words)
	fmt.Printf("time taken: %s\n", time.Since(start))
}
