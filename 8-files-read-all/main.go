package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	filePaths := os.Args[1:]

	if len(filePaths) <= 0 {
		log.Fatal("please pass the files to compare")
	}

	duplicateFiles := make([]string, 0)

	for _, filePath := range filePaths {
		counter := make(map[string]int, 0)
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if counter[line] > 0 {
				duplicateFiles = append(duplicateFiles, filePath)
				break
			}
			counter[line]++
		}
	}

	if len(duplicateFiles) <= 0 {
		fmt.Println("no files with duplicates")
	} else {
		fmt.Println("files with duplicates are:")
		for _, filePath := range duplicateFiles {
			fmt.Println(filePath)
		}
	}
}
