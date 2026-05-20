package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	filePath := flag.String("f", "", "File containing URLs (one per line)")
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Usage: go run gocode.go -f urls.txt")
		os.Exit(1)
	}

	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Burp proxy
	proxyURL, _ := url.Parse("http://127.0.0.1:8080")

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // required for Burp
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   20 * time.Second,
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		target := strings.TrimSpace(scanner.Text())
		if target == "" {
			continue
		}

		req, err := http.NewRequest("GET", target, nil)
		if err != nil {
			fmt.Println("Invalid URL:", target)
			continue
		}

		req.Header.Set("User-Agent", "Go-Burp-Client")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Request failed:", target, err)
			continue
		}

		fmt.Println("[OK]", target, "->", resp.Status)
		resp.Body.Close()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("File read error:", err)
	}
}
