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
	frequencyTable := make(map[byte]int)

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}

		if err == io.EOF {
			break
		}

		for _, v := range buffer[:n] {
			frequencyTable[v]++
		}
	}

	fmt.Println(frequencyTable)
}
