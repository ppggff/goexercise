package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type VisitedUrl struct {
	urls map[string]int
	mux sync.Mutex
}

func (v VisitedUrl) checkAndSave(u string) bool {
	v.mux.Lock()

	_, ok := v.urls[u];
	v.urls[u] ++

	v.mux.Unlock()
	return !ok
}

var visitedUrl = VisitedUrl{urls: make(map[string]int)}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, status chan int) {
	// This implementation doesn't do either:
	if depth <= 0 {
		status <- 0
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		status <- 0
		return
	}
	fmt.Printf("found: %s %q\n", url, body)

	childStatus := make(chan int)
	childCount := 0
	for _, u := range urls {
		if visitedUrl.checkAndSave(u) {
			childCount++
			go Crawl(u, depth-1, fetcher, childStatus)
		}
	}

	// wait for child
	for ; childCount > 0; childCount-- {
		<- childStatus
	}

	status <- 1
	return
}

func main() {
	u := "http://golang.org/"
	childStatus := make(chan int)
	if visitedUrl.checkAndSave(u) {
		go Crawl(u, 4, fetcher, childStatus)
	}
	<- childStatus
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
