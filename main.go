package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	csvFile := flag.String("csv", "questions.csv", "a csv file in the format of questions and answers")
	timeLimit := flag.Int("limit", 30, "time limit for quiz in seconds")
	flag.Parse()

	file, err := os.Open((*csvFile))
	if err != nil {
		exit(fmt.Sprintf("failed to open CSV File: %s \n", *csvFile))
	}
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()

	if err != nil {
		exit(fmt.Sprintf("failed to parse CSV File: %s \n", *csvFile))
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
problemLoop:
	for i, prob := range problems {
		fmt.Printf("Problem #%d: %s = \n", i+1, prob.q)
		answerCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\n you scored %d out of %d.\n", correct, len(problems))
			break problemLoop
		case answer := <-answerCh:
			if answer == prob.a {
				fmt.Println("Correct!")
				correct++
			} else {
				fmt.Println("Incorrect!")
			}
		}
	}
}

func parseLines(lines [][]string) []problem {

	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

type problem struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
