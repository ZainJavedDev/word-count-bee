package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func ProcessFile(filePath string, routines int) (CountsResult, int, time.Duration) {
	start := time.Now()
	results := make(chan CountsResult, routines)

	pwd, _ := os.Getwd()
	file, err := os.Open(pwd + "/" + filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fileSize := stat.Size()
	chunkSize := fileSize / int64(routines)

	reader := bufio.NewReader(file)

	for i := 0; i < routines; i++ {
		chunk := make([]byte, chunkSize)
		_, err := reader.Read(chunk)
		if err != nil {
			log.Fatal(err)
		}

		go Counts(chunk, results)
	}

	totalCounts := CountsResult{}

	for i := 0; i < routines; i++ {
		result := <-results
		totalCounts.LineCount += result.LineCount
		totalCounts.WordsCount += result.WordsCount
		totalCounts.VowelsCount += result.VowelsCount
		totalCounts.PunctuationCount += result.PunctuationCount
	}
	fmt.Println("no of ", routines)
	fmt.Println("Number of lines:", totalCounts.LineCount)
	fmt.Println("Number of words:", totalCounts.WordsCount)
	fmt.Println("Number of vowels:", totalCounts.VowelsCount)
	fmt.Println("Number of punctuation:", totalCounts.PunctuationCount)
	fmt.Println("Run Time:", time.Since(start))

	return totalCounts, routines, time.Since(start)

}
