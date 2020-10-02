package main

import (
	"bufio"
	"os"
	"strings"
	"time"
	"unique/pkg/util"
)

func main() {
	fName := "emoji.txt"

	start := time.Now()
	f, err := os.Stat(fName)
	if err != nil {
		panic(err)
	}

	file, err := os.Open(fName)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	dict, words := readFile(file)

	util.PrintMemUsage()
	util.PrintResult(fName, f.Size(), dict, words, start)
	util.PrintTop5(dict)
}

func readFile(f *os.File) (map[string]int, int) {
	dict := make(map[string]int)
	words := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		for _, w := range strings.Fields(scanner.Text()) {
			dict[w]++
			words++
		}
	}
	return dict, words
}
