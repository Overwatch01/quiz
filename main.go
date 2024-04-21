package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	fileNamePtr := flag.String("csv", "problems.csv", "a csv file in the format of 'question, answer' (default \"problems.csv\")")
	timerPtr := flag.Int("t", 10, "configured time for the game to run")
	flag.Parse()

	quizzes := readQuizFile(*fileNamePtr)
	setQuiz(quizzes, *timerPtr)
}

func readQuizFile(fileName string) map[string]int {
	quizzes := make(map[string]int)

	// os.open opens a specific file in
	// read only mode and this return
	// a pointer of type os.File
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Error while reading file:", err)
	}

	// close the file
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalln("Error occurred while closing file:", err)
		}
	}(file)

	// The csv.NewReader() function is called in
	// which the object os.File passed as its parameters
	// and this creates a new csv reader that reads
	//from the file
	reader := csv.NewReader(file)

	// ReadAll reads all the records from the CSV file
	// and Returns them as slice of slices of string
	// and an error if any
	records, err := reader.ReadAll()

	// Checks for the error
	if err != nil {
		log.Fatalln("Error while reading records", err)
	}

	// Loop to iterate through
	// and print each of the string slice
	for _, record := range records {
		quizResult, err := strconv.Atoi(record[1])
		if err != nil {
			log.Fatalf("Cannot get result for %d", record[0], err)
		}
		quizzes[record[0]] = quizResult
	}
	return quizzes
}

func setQuiz(quizzes map[string]int, timer int) {
	score := 0
	quizCount := 0
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple Quiz")
	fmt.Println("============")

	timerChan := time.NewTimer(time.Duration(timer) * time.Second)
	go func() {
		<-timerChan.C
		setScore(len(quizzes), score)
		os.Exit(1)
	}()

	for k, v := range quizzes {
		quizCount++
		fmt.Printf("Problem #%d: %v = ", quizCount, k)
		text, _ := inputReader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		answer, err := strconv.Atoi(text)
		if err != nil {
			continue
		}
		if answer == v {
			score++
		}
	}
	setScore(len(quizzes), score)
}

func setScore(quizCount, score int) {
	fmt.Printf("\n You scored %d out of %d.", score, quizCount)
}
