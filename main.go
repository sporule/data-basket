package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/lucasjones/reggen"
)

func main() {
	//timer
	start := time.Now()

	//load config
	flag.Parse()
	if flag.NArg() == 0 {
		panic("Could not find config file, please check the path")
	}
	configPath := flag.Arg(0)
	fileName, rowsLimit, columnNames, columnPatterns := readConfig(configPath)
	fmt.Printf("Writing %d rows to %v with column names: %v \n", rowsLimit, fileName, columnNames)

	//write column headings
	os.Remove(fileName)
	writeToFile(strings.Join(columnNames[:], ",")+"\n", fileName)

	//run goroutine jobs to generate rows
	goroutineSize := runtime.NumCPU() * 4
	fmt.Printf("Go Routine Size: %d \n", goroutineSize)
	rowDataChan := make(chan string)
	for i := 1; i < goroutineSize; i++ {
		go generateRows(columnPatterns, rowDataChan)
	}

	//write data to files
	var rows []string
	writeThreshold := 500000
	writtenRows := 0
	for {
		rows = append(rows, <-rowDataChan)
		if len(rows) >= min(rowsLimit-writtenRows, writeThreshold) {
			writeToFile(strings.Join(rows[:], "\n"), fileName)
			writtenRows += min(rowsLimit-writtenRows, writeThreshold)
			progress := float64(writtenRows) / float64(rowsLimit) * 100
			fmt.Printf("%.2f %% : Written %d rows to file %v \n", progress, writtenRows, fileName)
			rows = rows[:0]
			if writtenRows >= rowsLimit {
				break
			}
		}
	}
	close(rowDataChan)
	fmt.Printf("Total Time used: %v", time.Since(start))
}

func readConfig(configPath string) (fileName string, rowsLimit int, columnNames, columnPatterns []string) {
	f, err := ioutil.ReadFile(configPath)
	var config = make(map[string]interface{})
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(f, &config)
	if err != nil {
		panic(err)
	}
	fileName = config["fileName"].(string)
	rowsLimit = int(config["rows"].(float64))
	columns := config["columns"].(map[string]interface{})
	if fileName == "" || rowsLimit <= 0 || columns == nil {
		panic("Lack of information to generate test data")
	}
	for k, v := range columns {
		columnNames = append(columnNames, k)
		columnPatterns = append(columnPatterns, v.(string))
	}
	return fileName, rowsLimit, columnNames, columnPatterns
}

func writeToFile(content string, fileName string) {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		panic(err)
	}
}

func generateRows(columnPatterns []string, rowDataChan chan string) {
	defer func() {
		if r := recover().(error); r != nil {
			if r.Error() != "send on closed channel" {
				fmt.Printf("Recovering from panic in generateRows, error: %v \n", r)
			}
		}
	}()

	var generators []reggen.Generator

	for _, pattern := range columnPatterns {
		generators = append(generators, createGenerator(pattern))
	}

	for {
		var row []string
		for _, generator := range generators {
			data := generator.Generate(math.MaxInt8)
			row = append(row, data)
		}
		rowDataChan <- strings.Join(row[:], ",")
	}
}

func createGenerator(pattern string) reggen.Generator {
	generator, err := reggen.NewGenerator(pattern)
	if err != nil {
		panic(err)
	}
	return *generator
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
