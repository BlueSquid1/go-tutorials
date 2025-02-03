package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	counter := make(map[string]int)
	args := os.Args[1:]
	if len(args) == 0 {
		countLines(os.Stdin, counter)
	} else {
		for _, filePath := range args {
			file, err := os.Open(filePath)
			if err != nil {
				log.Fatal(err)
			}
			countLines(file, counter)
		}
	}

	for line, count := range counter {
		fmt.Printf("%v occured: %v\n", line, count)
	}
}

func countLines(f *os.File, counter map[string]int) {
	buffer := bufio.NewScanner(f)
	for buffer.Scan() {
		counter[buffer.Text()]++
	}
}
