package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

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

func ExtractParallel(worklist chan []string, link string, n *sync.WaitGroup) {
	defer n.Done()
	futureWork, err := ExtractUrls(link)
	if err != nil {
		fmt.Printf("failed to extract url: %v\n", err)
		worklist <- []string{}
	} else {
		worklist <- futureWork
	}
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("please pass a url")
	}

	GroupsOfUrls := [][]string{args}

	// only allow up to 5 extractors to run at a time
	sem := make(chan struct{}, 5)

	seen := make(map[string]bool)

	depth := 0
	for len(GroupsOfUrls) > 0 {
		worklist := make(chan []string)

		fmt.Printf("search depth is: %v\n", depth)
		depth++
		var n sync.WaitGroup
		// extract out urls from lists of urls
		for _, list := range GroupsOfUrls {
			// create goroutines
			for _, link := range list {
				if !seen[link] {
					sem <- struct{}{}
					defer func() { <-sem }()

					fmt.Printf("exploring url: %v\n", link)
					seen[link] = true
					n.Add(1)
					go ExtractParallel(worklist, link, &n)
				}
			}
		}
		go func() {
			n.Wait()
			close(worklist)
		}()

		GroupsOfUrls = [][]string{}
		// Wait for all gorountines for this depth to finish
		for futureWork := range worklist {
			if len(futureWork) > 0 {
				GroupsOfUrls = append(GroupsOfUrls, futureWork)
			}
		}
	}

	fmt.Println("URLs seen:")
	for k := range seen {
		fmt.Println(k)
	}
}
