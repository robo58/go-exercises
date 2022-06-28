package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type Problem struct {
	question string
	answer   string
}

type ScoreSheet struct {
	total int
	correct int
	wrong int
}

func main() {
	// read inputs from cmd line
	timeLimit := flag.Int("limit",30,"the time limit for the quiz in seconds")
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	flag.Parse()

	// initialize empty sheet
	scoreSheet := ScoreSheet{
		total: 0,
		correct: 0,
		wrong: 0,
	}

	// display result at end
	defer displayResult(&scoreSheet)

	// read file and get lines
	lines, readError := readFileByName(csvFilename)
	if readError != nil {
		exit("Failed to read the CSV file.")
	}

	// set total amount of questions on sheet
	scoreSheet.total = len(lines)

	// parse lines into problems
	problems := parseLines(lines)

	// start game loop
	playGame(problems, &scoreSheet, *timeLimit)
}

func playGame(problems []Problem, scoreSheet *ScoreSheet, timeLimit int){
	// initialize timer with duration
	timer := time.NewTimer(time.Duration(timeLimit)*time.Second)

	// create a channel that receives answers
	answerChannel := make(chan bool)

	// label for problem loop
	problemLoop:
		for index, problem := range problems {
			fmt.Printf("Problem #%d: %s = ", index+1, problem.question)
			// goroutine for checking answers
			go checkAnswer(problem, answerChannel)

			// if time runs out it stops the loop, otherwise it waits for answer from answerChannel
			select {
				case <-timer.C:
					break problemLoop
				case answerValue := <-answerChannel:
					if answerValue {
						scoreSheet.correct++
					}else {
						scoreSheet.wrong++
					}
			}
		}
}

func checkAnswer(problem Problem, channel chan bool){
	answer, answerError := readInput()
	if answerError != nil {
		exit(answerError.Error())
	}
	channel<-problem.answer==answer
}

func readInput() (string,error) {
	var input string
	_, scanError := fmt.Scanf("%s\n", &input)
	if scanError != nil {
		return "", errors.New("cannot read from input")
	}
	return input, nil
}

func readFileByName(filename *string) ([][]string, error){
	file, openError := os.Open(*filename)
	if openError != nil {
		return nil, errors.New(fmt.Sprintf("Failed to open the CSV file: %s\n", *filename))
	}
	r := csv.NewReader(file)
	return r.ReadAll()
}

func parseLines(lines [][]string) []Problem {
	ret := make([]Problem, len(lines))
	for i, line := range lines {
		ret[i] = Problem{question: line[0], answer: strings.TrimSpace(line[1])}
	}
	return ret
}

func displayResult(sheet *ScoreSheet){
	var average float32
	average = (float32(sheet.correct)/float32(sheet.total))*100
	fmt.Printf("\nQuiz is over you scored %d correct answers and %d wrong answers. \n", sheet.correct,sheet.wrong)
	fmt.Printf("That is total of %d/%d, percentage is %f%%. \n", sheet.correct,sheet.total, average)
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
