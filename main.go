package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/lucasjones/reggen"
)

func main() {
	//timer
	// start := time.Now()

	configPath := getConfigPath()
	fileName, rowsLimit, columns := readConfig(configPath)
	fmt.Printf("Writing %d rows to %v with column names: %v \n", rowsLimit, fileName, columns)

	//write column headings
	os.Remove(fileName)
	columnNames := getMapKeys(columns)
	writeToFile(strings.Join(columnNames[:], ",")+"\n", fileName)

	cacheDict, _ := getCacheTables(columns)
	fmt.Printf("%v", cacheDict)
	// //run goroutine jobs to generate rows
	// goroutineSize := runtime.NumCPU() * 4
	// fmt.Printf("Go Routine Size: %d \n", goroutineSize)
	// rowDataChan := make(chan string)
	// for i := 1; i < goroutineSize; i++ {
	// 	go generateRows(columnPatterns, rowDataChan)
	// }

	// //write data to files
	// var rows []string
	// writeThreshold := 500000
	// writtenRows := 0
	// for {
	// 	rows = append(rows, <-rowDataChan)
	// 	if len(rows) >= min(rowsLimit-writtenRows, writeThreshold) {
	// 		writeToFile(strings.Join(rows[:], "\n"), fileName)
	// 		writtenRows += min(rowsLimit-writtenRows, writeThreshold)
	// 		progress := float64(writtenRows) / float64(rowsLimit) * 100
	// 		fmt.Printf("%.2f %% : Written %d rows to file %v \n", progress, writtenRows, fileName)
	// 		rows = rows[:0]
	// 		if writtenRows >= rowsLimit {
	// 			break
	// 		}
	// 	}
	// }
	// close(rowDataChan)
	// fmt.Printf("Total Time used: %v", time.Since(start))
}

func getCacheTables(columns map[string]interface{}) (tables map[string][]string, mappings map[string]interface{}) {
	sortedByGroupMap := make(map[int][]interface{})
	tables = make(map[string][]string)
	for k, v := range columns {
		//check if the columns are part of a group
		_, foundKey := v.(map[string]interface{})["Group"]
		if !foundKey {
			//check if the columns has size property so we need to generate enough unique random examples
			_, foundKey := v.(map[string]interface{})["Size"]
			if foundKey {
				size := int(v.(map[string]interface{})["Size"].(float64))
				generator := createGenerator(v.(map[string]interface{})["Pattern"].(string))
				for len(tables[k]) < size {
					//generate enough unique example
					tables[k] = append(tables[k], generator.Generate(math.MaxInt8))
					if len(tables[k]) >= size {
						tables[k] = unique(tables[k])
					}
				}
			}
			continue
		}
		//group columns if they are in the same group
		group := int(v.(map[string]interface{})["Group"].(float64))
		column := make(map[string]interface{})
		column[k] = v
		sortedByGroupMap[group] = append(sortedByGroupMap[group], column)
	}
	for _, v := range sortedByGroupMap {
		//sort the group by size in asc
		sort.Slice(v, func(i, j int) bool {
			//get the map keys
			iKey := getMapKeys(v[i].(map[string]interface{}))[0]
			jKey := getMapKeys(v[j].(map[string]interface{}))[0]
			//check if they have size property, if not then they will be the biggest
			_, found := v[i].(map[string]interface{})[iKey].(map[string]interface{})["Size"]
			if !found {
				return false
			}
			_, found = v[j].(map[string]interface{})[jKey].(map[string]interface{})["Size"]
			if !found {
				return true
			}
			return int(v[i].(map[string]interface{})[iKey].(map[string]interface{})["Size"].(float64)) < int(v[j].(map[string]interface{})[jKey].(map[string]interface{})["Size"].(float64))
		})
	}
	return tables, mappings
}

func unique(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func getMapKeys(item map[string]interface{}) (keys []string) {
	for k, _ := range item {
		keys = append(keys, k)
	}
	return keys
}

func getConfigPath() string {
	//load config
	flag.Parse()
	var configPath string
	if flag.NArg() == 0 {
		fmt.Printf("Could not find config file, loading default config.json file \n")
		_, err := os.Stat("config.json")
		if err != nil {
			panic("Could not find default config.json file. Please pass in config file or create default config.json file \n")
		}
		configPath = "config.json"
	} else {
		configPath = flag.Arg(0)
	}
	return configPath
}

func readConfig(configPath string) (fileName string, rowsLimit int, columns map[string]interface{}) {
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
	columns = config["columns"].(map[string]interface{})
	if fileName == "" || rowsLimit <= 0 || columns == nil {
		panic("Lack of information to generate test data")
	}
	// for k, v := range columns {
	// 	columnNames = append(columnNames, k)
	// 	columnPatterns = append(columnPatterns, v.(map[string]interface{}))
	// }
	return fileName, rowsLimit, columns
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
