package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("expected 1 arguement")
	}

	url := args[0]
	containsHttps, err := regexp.MatchString(`^http[s]://.*`, url)
	if err != nil {
		log.Fatal(err)
	}

	if !containsHttps {
		url = "https://" + url
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", string(data))

	fmt.Printf("status code: %v\n", resp.StatusCode)
}
