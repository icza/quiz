// Package main is a simple command line application which takes a source input file
// that contains a list of words (one word in each line), and prints the
// longest compound-word in the list, which is also a concatenation of other sub-words that exist in the list.
// 
// The program allows the user to input different data. The input file can be
// specified/overridden using the -src flag.
// By default input file should be sorted alphabetically. If it is not sorted, the -unsorted flag must be provided.
//
// An optional -parallel flag may be provided in which case the program will perform a parallel search
// utilizing all CPU cores available on the machine.
// 
// Note #1: The task doesn't mention anything about how many words may combine the compound-word,
// so this solution finds the longest word that may be assembled from any number of sub-words (>= 2).
//
// Note #2: The task also doesn't mention if words contain only the letters of the English alphabet,
// so this solution properly handles unicode characters (does not assume English alphabet letters only).
//
// Note #3: In case there are multiple equally longest compound-words, the result may be any of these.
// If -parallel flag is not provided, the first in alphabetical order is returned.
//
// Note #4: This solution reads all words into memory and stores them in a slice.
// This is the only memory requirement.
//
// Speed Analyis
// Tested on Intel Core i5-4440 3.1 GHz (4 cores), Go 1.5. To verify, run
//     go test -bench .
//
// Running the tool on the list found at https://github.com/NodePrime/quiz/blob/master/word.list
// which contains 263533 words, the result ("antidisestablishmentarianisms")
// is found in ~0.7 ms! Less than a 1 ms out of ~263k words, that is wicked fast!
// Using the -parallel flag, it is slower (~5ms). This is partly due to the fact that the above list
// contains the longest compound-word in the first 3% of the list.
// When using a randomly generated list (10 million words, with length 5..44, having exactly one compound-word
// relatively at the end of the list), result is found in ~47 sec; when using -parallel it is found in ~15 sec (with 4 CPU cores).
//
// Original task can be found at https://github.com/NodePrime/quiz
package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"
)

var (
	src      = flag.String("src", `word.list`, "source file name containing words, 1 in each line")
	unsorted = flag.Bool("unsorted", false, "tells that input is not sorted")
	parallel = flag.Bool("parallel", false, "tells that search should utilize all CPU cores")
)

func main() {
	flag.Parse()

	// We have to read and keep all words in memory because when checking if compound,
	// sub-words may be at any position in the list.
	words, err := readLines(*src)
	if err != nil {
		log.Fatal(err)
	}

	// We rely on words being sorted, so:
	if *unsorted {
		sort.Strings(words)
	}

	longest, start := "", time.Now()
	if *parallel {
		longest = findLongestParal(words)
	} else {
		longest = findLongest(words)
	}

	log.Printf("Longest compound-word found: %s\n\t in %v\n", longest, time.Since(start))
}

// readLines reads all lines form the specified file.
func readLines(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	words := []string{}
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

// findLongest returns the longest compound-word in the sorted list,
// which is also a concatenation of other sub-words that exist in the list.
//
// Implemented as a separate function so it can be tested/benchmarked.
func findLongest(words []string) (longest string) {
	for _, word := range words {
		if len(longest) >= len(word) {
			continue // We already have a longer (or equally long) candidate
		}
		if compound(words, word) {
			longest = word
		}
	}
	return
}

// compound tells if the specified word can be assembled from others in the sorted list.
func compound(words []string, word string) bool {
	// Recursive solution:
	// Word is compound if there is another word that is a prefix of it
	// and the rest after the prefix is also a word or a compound word.
	// Max recursion depth: (rune) length of the word
	for i := range word { // range on a string steps by runes (not bytes)
		if i == 0 {
			continue
		}
		prefix := word[:i]
		if !contains(words, prefix) {
			continue
		}
		rest := word[i:]
		if contains(words, rest) || compound(words, rest) {
			return true
		}
	}

	return false
}

// contains tells if a string is contained in a sorted slice of strings.
func contains(ss []string, s string) bool {
	i := sort.SearchStrings(ss, s)
	return i < len(ss) && ss[i] == s
}

// findLongestParal returns the longest compound-word in the sorted list,
// which is also a concatenation of other sub-words that exist in the list.
//
// Implementation utilizes all CPU cores.
// Implemented as a separate function so it can be tested/benchmarked.
func findLongestParal(words []string) (longest string) {
	n := runtime.NumCPU()
	runtime.GOMAXPROCS(n)

	// This implementation delegates the expensive "is compound" checks to worker goroutines.

	wordchan := make(chan string, 1000)
	reschan := make(chan string, 1000)
	wg := sync.WaitGroup{}

	// Start "workers":
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for word := range wordchan {
				if compound(words, word) {
					reschan <- word
				}
			}
		}()
	}

	// Generate "jobs" and collect completed job results
	for _, word := range words {
		if len(longest) >= len(word) {
			continue // We already have a longer (or equally long) candidate
		}
		// Collect
		for j := len(reschan); j > 0; j-- {
			if candidate := <-reschan; len(candidate) > len(longest) {
				longest = candidate
			}
		}
		if len(longest) < len(word) {
			wordchan <- word // generate
		}
	}
	close(wordchan)
	wg.Wait() // Wait all workers to finish
	close(reschan)

	// Collect remaining completed job results
	for candidate := range reschan {
		if len(candidate) > len(longest) {
			longest = candidate
		}
	}

	return
}
