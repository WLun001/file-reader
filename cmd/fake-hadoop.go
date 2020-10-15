package main

import (
	"bufio"
	"encoding/json"
	"file-reader/pkg/util"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	lineLimit = 5
	filename  = "lorem.txt"
	tmpDir    = "tmp"
)

func main() {
	if err := runFn(); err != nil {
		log.Fatal(err)
	}
}

// read line by line
// when read until line 100
// split to worker func on another goroutine
// wait all mapper goroutines done
// run reducer
// print result
func runFn() error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()

	wg := sync.WaitGroup{}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	dirPath := fmt.Sprintf("%s/%s", wd, tmpDir)

	cmd := exec.Command("rm", "-rf", dirPath)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	_ = cmd.Wait()

	if err := os.Mkdir(dirPath, os.ModePerm); err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	accLines := ""
	counter := 0
	for scanner.Scan() {
		accLines += fmt.Sprintf("\n%s", scanner.Text())
		counter++
		if counter > lineLimit {
			wg.Add(1)
			go mapper(&wg, accLines, dirPath)
			// reset
			accLines = ""
			counter = 0
		}
	}

	wg.Wait()

	res, uniqueWord, reducerFile := reducer(dirPath)
	rel, err := filepath.Rel(wd, reducerFile)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Result written to %s\n", rel)
	log.Printf("First unique word: %s", uniqueWord)
	util.PrintTop5(res)

	return nil
}

// worker func
// find unique
// write result as json at tmp dir
func mapper(wg *sync.WaitGroup, text, path string) {
	dict := make(map[string]int)
	for _, w := range strings.Fields(text) {
		dict[w]++
	}

	filename := fmt.Sprintf("%s/mapper-%d.json", path, time.Now().UnixNano())

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = f.Close()
	}()

	data, err := json.MarshalIndent(dict, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	_ = ioutil.WriteFile(filename, data, 0644)

	wg.Done()
}

// read file by file
// process the result from file
// return key value pairs
// write result as json at tmp dir

// returns dict, uniqueWord, and result filename
func reducer(dir string) (map[string]int, string, string) {
	var files []string
	uniqueWord := ""
	dict := make(map[string]int)
	if err := filepath.Walk(dir, walkFn(&files)); err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		currentDict := make(map[string]int)

		f, err := os.Open(file)
		if err != nil {
			log.Fatal(f)
		}

		byteValue, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		if err := json.Unmarshal(byteValue, &currentDict); err != nil {
			log.Fatal(err)
		}

		for word, count := range currentDict {
			dict[word] += count
			if dict[word] == 1 {
				uniqueWord = word
			}
		}
		_ = f.Close()

	}

	filename := fmt.Sprintf("%s/reducer-%d.json", dir, time.Now().UnixNano())

	data, err := json.MarshalIndent(dict, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	_ = ioutil.WriteFile(filename, data, 0644)

	return dict, uniqueWord, filename
}

func walkFn(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if filepath.Ext(path) == ".json" {
			*files = append(*files, path)
		}
		return nil
	}
}
