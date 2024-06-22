package main

import (
	"encoding/csv"
	"log"
	"os"
	"fmt"
)

type Question struct {
	title string
	answer string
}

func main() {
	// first read a csv file (quiz question)
	file, err := os.Open("quiz_1.csv")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var questions []Question
	for _, record := range records {
		questions = append(questions, Question{
			title: record[0],
			answer: record[1],
		})
    }

	fmt.Println(questions)
	// second put in the question struct
	// third ask question in terminal and get answers
	// forth report quiz
}