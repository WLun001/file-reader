package main

import (
	"bufio"
	"encoding/json"
	"file-reader/pkg/util"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	log.Println("starting http server at :3000")
	http.HandleFunc("/word", readFileHandler)
	return http.ListenAndServe(":3000", nil)
}

func readFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := fmt.Sprintf("file-%d.txt", time.Now().Unix())

	fileUrl, ok := r.URL.Query()["file"]
	log.Println(fileUrl)

	if r.Method != "GET" {
		http.Error(w, "only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	if !ok || len(fileUrl[0]) < 1 {
		http.Error(w, "Url Param 'file' is missing", http.StatusBadRequest)
		return
	}

	err := util.DownloadFile(filename, fileUrl[0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	start := time.Now()
	f, err := os.Stat(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	dict, words, uniqueWord := fileProcessing(file)

	alloc, totalAlloc, sys, numGC := util.MemUsage()

	res := make(map[string]interface{})
	res["mem"] = map[string]string{
		"alloc":      alloc,
		"totalAlloc": totalAlloc,
		"sys":        sys,
		"numGC":      numGC,
	}
	res["timeTaken"] = time.Since(start).String()
	res["file"] = map[string]interface{}{
		"size": f.Size(),
		"MiB":  util.BToMb(uint64(f.Size())),
		"url":  fileUrl[0],
		"name": f.Name(),
	}
	res["words"] = map[string]interface{}{
		"uniqueWords":     len(dict),
		"wordCount":       words,
		"firstUniqueWord": uniqueWord,
		"top5":            util.Top5(dict),
	}

	log.Printf("deleting %s", filename)
	err = os.Remove(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonRes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonRes); err != nil {
		panic(err)
	}
}

// read file
// return unique word dict, word count, and first unique word
func fileProcessing(f *os.File) (map[string]int, int, string) {
	uniqueWord := ""
	dict := make(map[string]int)
	words := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		for _, w := range strings.Fields(scanner.Text()) {
			dict[w]++
			words++
			if dict[w] == 1 {
				uniqueWord = w
			}
		}
	}
	return dict, words, uniqueWord
}
