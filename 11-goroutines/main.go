package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	channel := make(chan string)
	go download("https://www.google.com", channel)
	fmt.Println(<-channel)
}

func download(url string, ch chan<- string) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	bytesDownloaded, err := io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}

	ch <- fmt.Sprintf("Bytes downloaded: %d for %s", bytesDownloaded, url)
}
