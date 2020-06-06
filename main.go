package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/lucasjones/reggen"
)

func main() {
	patterns := getPatterns()
	rowsLength := 100000
	threads := 1000
	thread := 1
	var rows []string
	rowDataChan := make(chan string)
	for thread <= threads {
		go generateRows(rowsLength/threads, patterns, rowDataChan)
		thread++
	}
	for {
		rows = append(rows, <-rowDataChan)
		fmt.Printf("%d \n", len(rows))
		if len(rows) >= rowsLength {
			f, err := os.OpenFile("data.csv", os.O_APPEND|os.O_WRONLY, 0644)
			content := strings.Join(rows[:], "\n")
			n, err := f.WriteString(content)
			fmt.Printf("%v %v \n", err, n)
			f.Close()
			break
		}
	}
	close(rowDataChan)
}

func generateRows(rowsLength int, patterns []string, rowDataChan chan string) {
	currentRow := 1
	for currentRow <= rowsLength {
		var row []string
		for _, pattern := range patterns {
			data, _ := generateString(pattern)
			row = append(row, data)
		}
		rowDataChan <- strings.Join(row[:], ",")
		currentRow++
	}
}
func generateString(pattern string) (data string, err error) {
	//example pattern "BT[£]{0,1}\\d{1,3}[A-Z]{1,3}"
	data, err = reggen.Generate(pattern, math.MaxInt32)
	if err != nil {
		panic(err)
	}
	return data, err
}

func getPatterns() (patterns []string) {
	patterns = append(patterns, "BT[£]{0,1}\\d{1,3}[A-Z]{1,3}")
	patterns = append(patterns, "\\d{1,8}\\.\\d{2}")
	return patterns
}
