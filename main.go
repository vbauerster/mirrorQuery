package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
)

type result struct {
	url  string
	resp *http.Response
}

func main() {
	flag.Parse()
	lines, err := readLines(flag.Arg(0))
	if err != nil {
		fmt.Printf("cannot open: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	first := make(chan result, len(lines))
	for _, url := range lines {
		go fetch(ctx, url, first)
	}

	r := <-first
	fmt.Printf("Winner is: %s\n", r.url)
}

func fetch(ctx context.Context, url string, first chan<- result) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return
	}
	first <- result{url, resp}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
