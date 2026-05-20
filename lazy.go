package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	urlTemplate string
	wordlist    string
	delay       time.Duration
	threads     int

	filterSize   string
	filterStatus string

	filterSizes   map[int64]bool
	filterStatusC map[int]bool
)

func main() {
	flag.StringVar(&urlTemplate, "u", "", "URL with FUZZ keyword")
	flag.StringVar(&wordlist, "w", "", "Wordlist path")
	flag.DurationVar(&delay, "delay", 0, "Delay between requests")
	flag.IntVar(&threads, "t", 10, "Number of threads")

	flag.StringVar(&filterSize, "fs", "", "Filter by response size in bytes (e.g. -fs 103,3386)")
	flag.StringVar(&filterStatus, "fc", "", "Filter by status code (e.g. -fc 404,400)")

	flag.Parse()

	if urlTemplate == "" || wordlist == "" {
		fmt.Println("Usage: lazy -u https://target/FUZZ -w wordlist.txt [-fs 103] [-fc 404]")
		os.Exit(1)
	}

	filterSizes = parseSizeFilters(filterSize)
	filterStatusC = parseStatusFilters(filterStatus)

	fmt.Println("Warming up connection...")
	time.Sleep(500 * time.Millisecond)
	fmt.Println("Warmup done. Starting fuzzing...")

	file, err := os.Open(wordlist)
	if err != nil {
		fmt.Println("Error opening wordlist:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	jobs := make(chan string)
	wg := sync.WaitGroup{}

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker(jobs, &wg)
	}

	fmt.Printf("[%s] Scanning:\n", time.Now().Format("15:04:05"))

	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word == "" {
			continue
		}
		jobs <- word
	}

	close(jobs)
	wg.Wait()
}

func worker(jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for word := range jobs {
		target := strings.Replace(urlTemplate, "FUZZ", word, 1)

		resp, err := client.Get(target)
		if err != nil {
			continue
		}

		body, _ := io.ReadAll(resp.Body)
		size := int64(len(body))
		status := resp.StatusCode
		location := resp.Header.Get("Location")
		path := "/" + word

		resp.Body.Close()

		// ---- Filters ----
		if filterSizes[size] {
			continue
		}
		if filterStatusC[status] {
			continue
		}
		// -----------------

		printResult(status, size, path, location)

		if delay > 0 {
			time.Sleep(delay)
		}
	}
}

func printResult(status int, size int64, path string, location string) {
	timestamp := time.Now().Format("15:04:05")

	if location != "" {
		fmt.Printf("[%s] %3d - %5dB - %-20s -> location: %s\n",
			timestamp,
			status,
			size,
			path,
			location,
		)
	} else {
		fmt.Printf("[%s] %3d - %5dB - %-20s\n",
			timestamp,
			status,
			size,
			path,
		)
	}
}

func parseSizeFilters(input string) map[int64]bool {
	m := make(map[int64]bool)
	if input == "" {
		return m
	}
	parts := strings.Split(input, ",")
	for _, p := range parts {
		v, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64)
		if err == nil {
			m[v] = true
		}
	}
	return m
}

func parseStatusFilters(input string) map[int]bool {
	m := make(map[int]bool)
	if input == "" {
		return m
	}
	parts := strings.Split(input, ",")
	for _, p := range parts {
		v, err := strconv.Atoi(strings.TrimSpace(p))
		if err == nil {
			m[v] = true
		}
	}
	return m
}