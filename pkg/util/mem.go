package util

import (
	"fmt"
	"runtime"
)

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
// from https://golangcode.com/print-the-current-memory-usage/
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", BToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", BToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", BToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func BToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func MemUsage() (alloc, totalAlloc, sys, numGC string) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	alloc = fmt.Sprintf("%v MiB", BToMb(m.Alloc))
	totalAlloc = fmt.Sprintf("%v MiB", BToMb(m.TotalAlloc))
	sys = fmt.Sprintf("%v MiB", BToMb(m.Sys))
	numGC = fmt.Sprintf("%v", m.NumGC)
	return
}
