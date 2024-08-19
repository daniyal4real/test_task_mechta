package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Record struct {
	A int `json:"a"`
	B int `json:"b"`
}

func read(fileName string) []byte {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

func numberOfBlocks(records []Record) int {
	numBlocks := len(records) / 100
	if len(records)%100 != 0 {
		numBlocks++
	}
	return numBlocks
}

func sumOfNumbers(records []Record, numWorkers int) (int, error) {
	sumCh := make(chan int)
	blockCh := make(chan int)
	numBlocks := numberOfBlocks(records)
	go func() {
		for i := 0; i < numBlocks; i++ {
			blockCh <- i
		}
		close(blockCh)
	}()

	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				idx, ok := <-blockCh
				if !ok {
					return
				}
				startIdx := idx * 100
				endIdx := (idx + 1) * 100
				if endIdx > len(records) {
					endIdx = len(records)
				}
				sum := 0
				for j := startIdx; j < endIdx; j++ {
					sum += records[j].A + records[j].B
				}
				sumCh <- sum
			}
		}()
	}
	totalSum := 0
	for i := 0; i < numBlocks; i++ {
		sum := <-sumCh
		totalSum += sum
	}
	close(sumCh)
	return totalSum, nil
}

func parseValues(bytes []byte) ([]Record, error) {
	var records []Record
	if err := json.Unmarshal(bytes, &records); err != nil {
		log.Fatal(err)
	}
	return records, nil
}

func main() {
	numWorkers := flag.Int("workers", 10, "number of workers")
	flag.Parse()
	bytes := read("input.json")

	records, err := parseValues(bytes)
	if err != nil {
		log.Fatal(err)
	}

	result, err := sumOfNumbers(records, *numWorkers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Общая сумма:", result)
}