package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"golang.org/x/net/html"
)

func ExtractUrls(targetUrl string) ([]string, error) {
	doc, err := getHtmlFromUrl(targetUrl)
	if err != nil {
		return nil, err
	}
	links := []string{}
	urlExtractor := func(node *html.Node) {
		if node.Data == "a" {
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					if _, err := url.Parse(attr.Val); err != nil {
						continue
					}
					links = append(links, attr.Val)
				}
			}
		}
	}
	forEachNode(doc, urlExtractor)
	return links, nil
}

func forEachNode(node *html.Node, f func(node *html.Node)) {
	if node.Type == html.ElementNode {
		f(node)
	}
	for c := range node.ChildNodes() {
		forEachNode(c, f)
	}
}

func getHtmlFromUrl(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to do GET request for %s: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("return code was: %d", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %v", err)
	}
	return doc, nil
}

func ExtractParallel(worklist chan []string, link string, n *sync.WaitGroup, done chan struct{}) {
	defer n.Done()
	// artifically slow extractor to see performance
	time.Sleep(time.Second)
	futureWork, err := ExtractUrls(link)
	if err != nil {
		fmt.Printf("failed to extract url: %v\n", err)
		futureWork = []string{}
	}
	// if done has been closed then don't write to the worklist
	select {
	case <-done:
	case worklist <- futureWork:
	}
}

func main() {
	// This does a breadth first search for all links on a website
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("please pass a url")
	}

	GroupsOfUrls := [][]string{args}

	done := make(chan struct{})
	go func() {
		// press ctrl + D to stop
		os.Stdin.Read(make([]byte, 1))
		// closing the channel means that the select statement will trigger
		// indefinately for any reads to this channel
		close(done)
	}()

	// Only allow up to 5 extractors to run at a time.
	sem := make(chan struct{}, 5)

	// Only run extractor on URLs that haven't been seen before.
	seen := make(map[string]bool)

	enableVerbose := true
	// tick defaults to null
	var tick <-chan time.Time
	if enableVerbose {
		tick = time.Tick(500 * time.Millisecond)
	}

	depth := 0
	for len(GroupsOfUrls) > 0 {
		// Results from each extractor is store in this channel.
		worklist := make(chan []string)

		depth++
		var n sync.WaitGroup
		// extract out urls from lists of urls
		for _, list := range GroupsOfUrls {
			for _, link := range list {
				if !seen[link] {
					// wait if more than 5 gorountines are running
					sem <- struct{}{}
					defer func() { <-sem }()

					seen[link] = true
					// Count each gorountine that starts
					n.Add(1)
					go ExtractParallel(worklist, link, &n, done)
				}
			}
		}
		go func() {
			// Wait for all gorountines for this depth to finish
			n.Wait()
			// when a channel is closed, all the existing values can still be read
			// after all the values have been read it will return the default value
			// for the channel type for infinite.
			close(worklist)
		}()

		// clear the URLs so the urls for the next depth level can been appended
		GroupsOfUrls = [][]string{}
		// look over channel. Will only stop looping once all values have been read from the closed channel.
	loop:
		for {
			select {
			case futureWork, ok := <-worklist:
				if !ok {
					break loop
				}
				GroupsOfUrls = append(GroupsOfUrls, futureWork)
			case <-tick: // if tick defaults to null then this case will never occur
				fmt.Printf("Currently at depth: %v, total URLs seen: %v\n", depth, len(seen))
			case <-done:
				// Cancel current goroutines
				close(worklist)
				fmt.Println("exit early")
				return
			}
		}
	}

	fmt.Println("URLs seen:")
	for k := range seen {
		fmt.Println(k)
	}
}
