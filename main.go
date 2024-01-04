package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("lesmiserables.txt")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	const chunkSize = 10000
	buffer := make([]byte, chunkSize)

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		if err == io.EOF {
			break
		}

		fmt.Println(string(buffer[:n]))
	}
}
