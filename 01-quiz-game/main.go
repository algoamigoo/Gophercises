package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
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
	shuffle := flag.Bool("shuffle", false, "shuffle the quiz questions order")
	flag.Parse()

	// Open the CSV File
	file, err := os.Open(*csvFilename)
	if err != nil {
		log.Fatalf("Failed to open file '%s': %v", *csvFilename, err)
	}
	defer file.Close()

	// 1. Parse upfront into a slice with string cleaning
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

		// Clean up CSV values on import
		problems = append(problems, problem{
			q: strings.TrimSpace(record[0]),
			a: cleanString(record[1]),
		})
	}

	// Shuffle questions if flag is enabled
	if *shuffle {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		r.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	fmt.Printf("--- Quiz Starts (Time Limit: %ds) ---\n", *timeLimit)

	score := 0

	// Deferred closure ensures final score prints whether time expires or quiz finishes
	defer func() {
		fmt.Printf("\nYour score is %d out of %d!\n", score, len(problems))
	}()

	// Initialize timer
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	// 2. Run the quiz loop
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
			// Normalize user input and compare against cleaned target answer
			if cleanString(response) == p.a {
				score++
			}
		}
	}
}

// Helper function to strip leading/trailing whitespace and normalize to lowercase
func cleanString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
