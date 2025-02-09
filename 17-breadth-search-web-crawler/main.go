package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

type UrlExtractorType func(item string) ([]string, error)

func BreadthFirst(f UrlExtractorType, workList []string) ([]string, error) {
	seen := make(map[string]bool)

	for len(workList) > 0 {
		items := workList
		workList = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				childResults, err := f(item)
				if err != nil {
					return nil, fmt.Errorf("lambda failed for url %s: %v", item, err)
				}
				workList = append(workList, childResults...)
			}
		}
	}
	urls := make([]string, len(seen))
	i := 0
	for k := range seen {
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

	urls, err := BreadthFirst(ExtractUrls, args)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(urls)
}
