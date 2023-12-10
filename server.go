package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"
)

// SortRequest represents the JSON structure for incoming sorting requests.
type SortRequest struct {
	ToSort [][]int `json:"to_sort"`
}

// SortResponse represents the JSON structure for sorting responses.
type SortResponse struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNs       int64   `json:"time_ns"`
}

// sortSequential sorts each sub-array sequentially and measures time taken in nanoseconds.
func sortSequential(toSort [][]int) ([][]int, int64) {
	start := time.Now()
	sortedArrays := make([][]int, len(toSort))

	for i, arr := range toSort {
		sorted := make([]int, len(arr))
		copy(sorted, arr)
		sort.Ints(sorted)
		sortedArrays[i] = sorted
	}

	elapsed := time.Since(start)
	return sortedArrays, elapsed.Nanoseconds()
}

// sortConcurrent sorts each sub-array concurrently using goroutines and measures time taken in nanoseconds.
func sortConcurrent(toSort [][]int) ([][]int, int64) {
	start := time.Now()
	var wg sync.WaitGroup
	sortedArrays := make([][]int, len(toSort))

	for i, arr := range toSort {
		wg.Add(1)
		go func(index int, input []int) {
			defer wg.Done()
			sorted := make([]int, len(input))
			copy(sorted, input)
			sort.Ints(sorted)
			sortedArrays[index] = sorted
		}(i, arr)
	}
	wg.Wait()

	elapsed := time.Since(start)
	return sortedArrays, elapsed.Nanoseconds()
}

// processSingleHandler handles requests to sort arrays sequentially.
func processSingleHandler(w http.ResponseWriter, r *http.Request) {
	var req SortRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sortedArrays, timeTaken := sortSequential(req.ToSort)
	response := SortResponse{
		SortedArrays: sortedArrays,
		TimeNs:       timeTaken,
	}
	json.NewEncoder(w).Encode(response)
}

// processConcurrentHandler handles requests to sort arrays concurrently.
func processConcurrentHandler(w http.ResponseWriter, r *http.Request) {
	var req SortRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sortedArrays, timeTaken := sortConcurrent(req.ToSort)
	response := SortResponse{
		SortedArrays: sortedArrays,
		TimeNs:       timeTaken,
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Register handlers for /process-single and /process-concurrent endpoints
	http.HandleFunc("/process-single", processSingleHandler)
	http.HandleFunc("/process-concurrent", processConcurrentHandler)

	// Start the HTTP server on port 8000
	http.ListenAndServe(":8000", nil)
}
