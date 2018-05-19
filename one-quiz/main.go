package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

//This file takes a CSV of the format `problem`,`answer`
//and lets you play a timed quiz on the command line

//Open the CSV file and read it into a Q:A map
func prepareFile(file string) map[string]string {
	csvFile, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(bufio.NewReader(csvFile))
	records := make(map[string]string)
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		records[line[0]] = line[1]
	}
	return records
}

func main() {
	//Create time limit flag
	var limit int
	flag.IntVar(&limit, "limit", 30, "time limit for quiz")
	flag.Parse()

	//Assign Q:A map to var
	problems := prepareFile("problems.csv")

	totalQs := len(problems)
	correct := 0
	ch := make(chan string)

	//Run the quiz
	go func() {
		for q := range problems {
			fmt.Println("Problem: " + q + " = ")
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(answer)
			if answer == problems[q] {
				correct++
			}
		}
		ch <- "Your score was " + strconv.Itoa(correct) + " out of " + strconv.Itoa(totalQs)
	}()

	//Handle end-of-quiz, whether finished or timed out
	select {
	case score := <-ch:
		fmt.Println(score)
	case <-time.After(time.Duration(limit) * time.Second):
		fmt.Println("Timed out. Your score was " + strconv.Itoa(correct) + " out of " + strconv.Itoa(totalQs))
	}
}
