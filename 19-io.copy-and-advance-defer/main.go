package main

import (
	"io"
	"net/http"
	"os"
	"path"
)

func fetch(url string) (n int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}

	f, err := os.Create(local)
	if err != nil {
		return 0, err
	}
	defer func() {
		// overwrite the error message returned from the function
		err = f.Close()
	}()

	numCopied, err := io.Copy(f, resp.Body)
	if err != nil {
		return 0, err
	}

	return int(numCopied), nil
}

func main() {
	args := os.Args[1:]
	for _, arg := range args {
		fetch(arg)
	}
}
