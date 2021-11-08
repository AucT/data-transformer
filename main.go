package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Configuration struct {
	InputFileName           string //your table data
	DataSourceFileName      string //data source to make map and process your table later with
	NotMappedOutputFileName string //your table data with only rows that not processed
	MappedOutputFileName    string //your table data with only rows that did processed
	CombinedOutputFileName  string //your table data with all rows. Not processed will have blank new column
	InputColumn             string //column in your table to map with
	OutputColumn            string //mapped column to add to your table
	ConvertToFullImdbString bool   //convert 123 to tt0000123 in output tables
}

var (
	config Configuration
	csvMap = make(map[int]int)
)

func main() {

	start := time.Now()

	initialize()
	mapCsv()
	generateNewCsv()

	elapsed := time.Since(start)
	fmt.Println("Time took ", elapsed)
}

func initialize() {
	inputFileName := flag.String("inputFileName", "", "")
	dataSourceFileName := flag.String("dataSourceFileName", "", "")
	notMappedOutputFileName := flag.String("notMappedOutputFileName", "", "")
	mappedOutputFileName := flag.String("mappedOutputFileName", "", "")
	combinedOutputFileName := flag.String("combinedOutputFileName", "", "")
	inputColumn := flag.String("inputColumn", "", "")
	outputColumn := flag.String("outputColumn", "", "")
	configFileName := flag.String("config", "data-transformer.json", "data-transformer.json")

	flag.Parse()

	readConfig(*configFileName)

	if len(*inputFileName) != 0 {
		config.InputFileName = *inputFileName
	}
	if len(*dataSourceFileName) != 0 {
		config.DataSourceFileName = *dataSourceFileName
	}
	if len(*notMappedOutputFileName) != 0 {
		config.NotMappedOutputFileName = *notMappedOutputFileName
	}
	if len(*mappedOutputFileName) != 0 {
		config.MappedOutputFileName = *mappedOutputFileName
	}

	if len(*combinedOutputFileName) != 0 {
		config.CombinedOutputFileName = *combinedOutputFileName
	}

	if len(*inputColumn) != 0 {
		config.InputColumn = *inputColumn
	}
	if len(*outputColumn) != 0 {
		config.OutputColumn = *outputColumn
	}
}

func readConfig(configFileName string) {
	file, err := os.Open(configFileName)
	defer file.Close()

	if err != nil {
		fmt.Println("Can't open "+configFileName, err)

		ex, err := os.Executable()
		checkError("Can't find executable path ", err)
		exPath := filepath.Dir(ex)

		filePath := filepath.Join(exPath, configFileName)
		fmt.Println("Trying to open in executable folder " + filePath)

		file, err = os.Open(filePath)
		checkError("Can't open "+filePath, err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	checkError("Error parsing config:", err)
}
func mapCsv() {
	csvFile, err := os.Open(config.DataSourceFileName)

	defer csvFile.Close()

	if err != nil {
		fmt.Println("Can't open "+config.DataSourceFileName, err)

		ex, err := os.Executable()
		checkError("Can't find executable path ", err)
		exPath := filepath.Dir(ex)

		filePath := filepath.Join(exPath, config.DataSourceFileName)
		fmt.Println("Trying to open in executable folder " + filePath)
		csvFile, err = os.Open(filePath)
		checkError("Can't open "+filePath, err)
	}

	fmt.Println("Successfully Opened CSV file")

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	checkError("Can't read "+config.DataSourceFileName, err)

	var columnKey int
	var columnValue int

	for i, column := range csvLines[0] {
		if column == config.InputColumn {
			columnKey = i
		}
		if column == config.OutputColumn {
			columnValue = i
		}
	}

	for _, line := range csvLines {
		key, _ := strconv.Atoi(line[columnKey])
		csvMap[key], _ = strconv.Atoi(line[columnValue])
	}
}

func generateNewCsv() {
	csvFile, err := os.Open(config.InputFileName)
	defer csvFile.Close()
	checkError("Can't open "+config.InputFileName, err)

	fmt.Println("Successfully Opened " + config.InputFileName)

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	checkError("Can't read "+config.InputFileName, err)

	var columnKey int

	for i, column := range csvLines[0] {
		if column == config.InputColumn {
			columnKey = i
			break
		}
	}
	var mappedData [][]string
	var notMappedData [][]string
	var combinedData [][]string

	mappedData = append(mappedData, append(csvLines[0], config.OutputColumn))
	notMappedData = append(notMappedData, append(csvLines[0], config.OutputColumn))
	combinedData = append(combinedData, append(csvLines[0], config.OutputColumn))

	for i, line := range csvLines {
		if i == 0 {
			continue
		}
		key, _ := strconv.Atoi(line[columnKey])

		value := ""
		if csvMap[key] != 0 {
			if config.ConvertToFullImdbString {
				value = fmt.Sprintf("tt%07d", csvMap[key])
			} else {
				value = strconv.Itoa(csvMap[key])
			}
			mappedData = append(mappedData, append(line, value))
		} else {
			notMappedData = append(notMappedData, append(line, value))
		}
		combinedData = append(combinedData, append(line, value))
	}
	dumpCsv(mappedData, config.MappedOutputFileName)
	dumpCsv(notMappedData, config.NotMappedOutputFileName)
	dumpCsv(combinedData, config.CombinedOutputFileName)
}

func dumpCsv(data [][]string, fileName string) {
	file, err := os.Create(fileName)
	checkError("Cannot create file "+fileName, err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		checkError("Cannot write to file "+fileName, err)
	}
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
