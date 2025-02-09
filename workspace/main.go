package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

type BreadthExtractor struct {
	seen map[string]bool
}

type UrlExtractorType func(item string) ([]string, error)

func (b *BreadthExtractor) BreadthFirst(f UrlExtractorType, workList []string) ([]string, error) {
	newWorkItems := []string{}
	for _, workItem := range workList {
		if !b.seen[workItem] {
			b.seen[workItem] = true
			newUrls, err := f(workItem)
			if err != nil {
				return nil, fmt.Errorf("failed to call lambda with url %s: %v", workItem, err)
			}
			newWorkItems = append(newWorkItems, newUrls...)
		}
	}

	if len(newWorkItems) > 0 {
		b.BreadthFirst(f, newWorkItems)
	}

	urls := make([]string, len(b.seen))
	i := 0
	for k := range b.seen {
		urls[i] = k
		i++
	}
	return urls, nil
}

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

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("please pass a url")
	}

	b := BreadthExtractor{}
	b.seen = make(map[string]bool)
	urls, err := b.BreadthFirst(ExtractUrls, args)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(urls)
}
