package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	filePath := flag.String("f", "", "Path to raw payloads file")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Usage: go run escaper.go -f payloads.txt")
		os.Exit(1)
	}

	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Step 1: JSON-safe escaping
		escaped, err := json.Marshal(line)
		if err != nil {
			continue
		}

		// Remove surrounding quotes added by json.Marshal
		result := string(escaped[1 : len(escaped)-1])

		// Step 2: VAPT-friendly escaping
		result = strings.ReplaceAll(result, "'", "\\'")

		fmt.Println(result)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Read error:", err)
	}
}
