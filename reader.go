package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
	"unique/pkg/util"
)

const kb = 1024

var file = flag.String("f", "emoji.txt", "text file")
var limit = flag.Int("l", 500, "limit in KB")

func main() {
	flag.Parse()

	start := time.Now()

	wg := sync.WaitGroup{}
	channel := make(chan string)
	dict := make(map[string]int)
	done := make(chan bool, 1)

	words := 0

	go func() {
		for s := range channel {
			words++
			dict[s]++
		}

		done <- true
	}()

	var current int64

	f, err := os.Stat(*file)
	if err != nil {
		panic(err)
	}

	goroutines := 1
	limitInBytes := int64(*limit * kb)
	if f.Size() > limitInBytes {
		goroutines = int(f.Size() / limitInBytes)
	}

	for i := 0; i < goroutines; i++ {
		wg.Add(1)

		go func() {
			read(current, limitInBytes, *file, channel)
			fmt.Printf("%d goroutine has been completed \n", i)
			wg.Done()
		}()

		current += limitInBytes + 1
	}

	wg.Wait()
	close(channel)

	<-done
	close(done)

	util.PrintMemUsage()
	util.PrintResult(*file, f.Size(), dict, words, start)
	util.PrintTop5(dict)
}

func read(offset int64, limit int64, fileName string, channel chan string) {
	file, err := os.Open(fileName)

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	if err != nil {
		panic(err)
	}

	// Move the pointer of the file to the start of designated chunk.
	file.Seek(offset, 0)
	reader := bufio.NewReader(file)

	// This block of code ensures that the start of chunk is a new word. If
	// a character is encountered at the given position it moves a few bytes till
	// the end of the word.
	if offset != 0 {
		_, err = reader.ReadBytes(' ')
		if err == io.EOF {
			fmt.Println("EOF")
			return
		}

		if err != nil {
			panic(err)
		}
	}

	var cummulativeSize int64
	for {
		if cummulativeSize > limit {
			break
		}

		b, err := reader.ReadBytes(' ')

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		cummulativeSize += int64(len(b))
		s := strings.TrimSpace(string(b))
		if s != "" {
			channel <- s
		}
	}
}
