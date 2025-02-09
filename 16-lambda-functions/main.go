package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func forEachNode(node *html.Node, pre, post func(node *html.Node)) {
	// pre and post are function values with type func(node *html.Node)
	//forEachNode is considered a first class function because it takes in a function as an argument
	if node.Type == html.ElementNode {
		pre(node)
	}
	for c := range node.ChildNodes() {
		forEachNode(c, pre, post)
	}
	if node.Type == html.ElementNode {
		post(node)
	}
}

func hasChildenElements(node *html.Node) bool {
	for c := range node.ChildNodes() {
		if c.Type == html.ElementNode {
			return true
		}
	}
	return false
}

func printHtml(node *html.Node) {
	depth := 0
	pre := func(node *html.Node) {
		if hasChildenElements(node) {
			fmt.Printf("%*s<%s>\n", depth*2, " ", node.Data)
		} else {
			fmt.Printf("%*s<%s/>\n", depth*2, " ", node.Data)
		}
		depth++
	}
	post := func(node *html.Node) {
		depth--
		if hasChildenElements(node) {
			fmt.Printf("%*s</%s>\n", depth*2, " ", node.Data)
		}
	}
	forEachNode(node, pre, post)
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
		log.Fatal("Please enter a URL")
	}

	for _, url := range args {
		doc, err := getHtmlFromUrl(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to process url %s: %v", url, err)
			continue
		}

		printHtml(doc)
	}
}
