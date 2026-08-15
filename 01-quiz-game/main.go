package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type problem struct {
	q string
	a string
}

func main() {
	// Define command-line flags
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 15, "the time limit for the quiz in seconds")
	flag.Parse()

	// Open the CSV File
	file, err := os.Open(*csvFilename)
	if err != nil {
		log.Fatalf("Failed to open file '%s': %v", *csvFilename, err)
	}
	defer file.Close()

	// 1. Parse upfront into a slice
	reader := csv.NewReader(file)
	var problems []problem

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading CSV record: %v", err)
		}
		problems = append(problems, problem{
			q: record[0],
			a: strings.TrimSpace(record[1]),
		})
	}

	fmt.Printf("--- Quiz Starts (Time Limit: %ds) ---\n", *timeLimit)

	score := 0
	
	// Deferred closure ensures final score prints whether time expires or quiz finishes
	defer func() {
		fmt.Printf("\nYour score is %d out of %d!\n", score, len(problems))
	}()

	// Initialize timer with the custom flag value
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	// 2. Run the quiz
	for i, p := range problems {
		fmt.Printf("Question %d: %s = ", i+1, p.q)

		ansCh := make(chan string)
		go func() {
			var res string
			fmt.Scan(&res)
			ansCh <- res
		}()

		select {
		case <-timer.C:
			fmt.Println("\nTime's up!")
			return
		case response := <-ansCh:
			if strings.TrimSpace(response) == p.a {
				score++
			}
		}
	}
}