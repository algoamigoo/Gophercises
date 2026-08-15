package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func main()  {
	// open the CSV File
	file,err :=os.Open("problems.csv")
	if err!=nil {
		log.Fatalf("Failed to open file : %v",err)
	}

	defer file.Close()
	
	//initialize CSV reader
	reader := csv.NewReader(file)
	fmt.Println("quiz starts")

	// read line by line
	i := 1
	score :=0
	for {
		record,err := reader.Read()
		if(err==io.EOF){
			break // End of file reached
		}
		if(err!=nil){
			log.Fatalf("Error reading CSV record: %v", err)
		}
		question :=record[0]
		answer :=record[1]
		
		fmt.Printf("question %d : %s = \n", i, question)
		i++
		var response string
		_, err = fmt.Scan(&response)
		if err != nil {
			fmt.Println("Invalid input:", err)
			return
		}
		if response ==answer {
			score++
		}
	}
	fmt.Printf("\nYour score is %d out of %d!\n", score, i-1)
}