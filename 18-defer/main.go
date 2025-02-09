package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func fileDefer(f *os.File) {
	fmt.Printf("closing file: %s\n", f.Name())
	f.Close()
}

func processFile(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer fileDefer(f)
	r := bufio.NewReader(f)
	length := r.Size()
	fmt.Printf("file: %s has length is: %d\n", fileName, length)
}

func main() {
	files := []string{"a.txt", "b.txt"}

	//dangerous because only calls close after all files have been open
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer fileDefer(f)
		r := bufio.NewReader(f)
		length := r.Size()
		fmt.Printf("file: %s has length is: %d\n", file, length)
	}

	//instead do this
	for _, file := range files {
		processFile(file)
	}
}
