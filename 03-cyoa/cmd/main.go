package main

import (
	"fmt"
	cyoa "gophercises/03-cyoa"
	"log"
	"net/http"
	"os"
)

func main() {
	file, err := os.Open("../gopher.json")
	if err!=nil{
		log.Fatal("Cannot read json file")
	}
	defer file.Close()
	
	book,err := cyoa.JsonStory(file)
	if err != nil {
		log.Fatalf("Cannot parse json: %v", err)
	}

	fmt.Printf("Successfully loaded %d chapters!\n", len(book))
	
	h := cyoa.NewHandler(book)

	fmt.Println("Server is starting on http://localhost:8080...")
	
	// ListenAndServe blocks here indefinitely until the server stops
	log.Fatal(http.ListenAndServe(":8080", h))
}
