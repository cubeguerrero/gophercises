package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func readCSV(filename string) ([]problem, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	r := csv.NewReader(file)
	problems := []problem{}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		problems = append(problems, problem{question: row[0], answer: row[1]})
	}
	return problems, nil
}

func askQuestion(p problem, reader bufio.Reader, c chan string) {
	fmt.Print(p.question, " ")
	answer, _ := reader.ReadString('\n')
	c <- strings.Trim(answer, "\n")
}

func main() {
	filenamePtr := flag.String("filename", "problems.csv", "filename of the CSV quiz")
	timerLengthPtr := flag.Int("time", 10, "timer length for the quiz (in seconds)")
	flag.Parse()
	questionsAndAnswers, err := readCSV(*filenamePtr)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	score := 0
	stdinReader := bufio.NewReader(os.Stdin)
	answerChannel := make(chan string)
	timerChannel := time.NewTimer(time.Duration(*timerLengthPtr) * time.Second)
	for _, p := range questionsAndAnswers {
		go askQuestion(p, *stdinReader, answerChannel)

		select {
		case <-timerChannel.C:
			fmt.Println("\nTime's up!")
			fmt.Printf("You got %d out of %d", score, len(questionsAndAnswers))
			return
		case answer := <-answerChannel:
			if answer == p.answer {
				score++
			}
		}
	}

	fmt.Printf("You got %d out of %d", score, len(questionsAndAnswers))
}
