package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

type KeywordCount struct {
	Keyword string
	Count   int
}

func ProcessLogFile(filePath string, keywords []string) ([]KeywordCount, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	keywordMap := make(map[string]string)
	for _, kw := range keywords {
		keywordMap[strings.ToLower(kw)] = kw
	}

	linesCh := make(chan string, 100)
	resultsCh := make(chan map[string]int, 10)
	var wg sync.WaitGroup

	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go CountKeywords(linesCh, resultsCh, keywordMap, &wg)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	scanner := bufio.NewScanner(file)
	go func() {
		for scanner.Scan() {
			linesCh <- scanner.Text()
		}
		close(linesCh)
	}()

	finalCounts := make(map[string]int)
	for result := range resultsCh {
		for keyword, count := range result {
			finalCounts[keyword] += count
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var counts []KeywordCount
	for keyword, count := range finalCounts {
		counts = append(counts, KeywordCount{Keyword: keyword, Count: count})
	}

	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})

	return counts, nil
}

func CountKeywords(linesCh <-chan string, resultsCh chan<- map[string]int, keywordMap map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

	localCounts := make(map[string]int)

	for line := range linesCh {
		lowerLine := strings.ToLower(line)
		
		for lowerKeyword, originalKeyword := range keywordMap {
			if strings.Contains(lowerLine, lowerKeyword) {
				localCounts[originalKeyword]++
			}
		}
	}

	resultsCh <- localCounts
}

func main() {
	keywords := []string{"INFO", "ERROR", "DEBUG"}
	
	counts, err := ProcessLogFile("log.txt", keywords)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("=== Log Processing Results ===")
	for _, count := range counts {
		fmt.Printf("%s: %d\n", count.Keyword, count.Count)
	}
}