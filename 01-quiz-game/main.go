package main

import (
	"encoding/csv"
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
	file, err := os.Open("problems.csv")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Parse all problems upfront so total count is known
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

	fmt.Println("Quiz starts!")

	score := 0
	// Deferred evaluation ensures final score prints on normal exit or timer expiration
	defer func() {
		fmt.Printf("\nYour score is %d out of %d!\n", score, len(problems))
	}()

	timer := time.NewTimer(15 * time.Second)

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