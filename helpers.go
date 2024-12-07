package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

// Count words in a string
func countWords(input string) int {
	words := strings.Fields(input)
	return len(words)
}

// Preprocess text based on configuration
func preprocessText(input string, config map[string][]string) string {
	if config == nil {
		return input
	}
	for _, tag := range config["ignore_tags"] {
		input = strings.ReplaceAll(input, tag, "")
	}
	for _, char := range config["remove_chars"] {
		input = strings.ReplaceAll(input, char, "")
	}
	return input
}

// Load configuration file
func loadConfig(path string) map[string][]string {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error: Unable to read configuration file.")
		os.Exit(1)
	}
	var config map[string][]string
	if err := json.Unmarshal(file, &config); err != nil {
		fmt.Println("Error: Invalid configuration file format.")
		os.Exit(1)
	}
	return config
}

// Process file input based on file type
func processFile(filePath, sheetName, columnName string, config map[string][]string) int {
	if strings.HasSuffix(filePath, ".txt") {
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error: Unable to read text file.")
			os.Exit(1)
		}
		return countWords(preprocessText(string(content), config))
	}

	if strings.HasSuffix(filePath, ".csv") {
		return processCSV(filePath, columnName, config)
	}

	if strings.HasSuffix(filePath, ".xls") {
		return processXLS(filePath, sheetName, columnName, config)
	}

	fmt.Println("Error: Unsupported file type.")
	os.Exit(1)
	return 0
}

// Process CSV files
func processCSV(filePath, columnName string, config map[string][]string) int {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error: Unable to open CSV file.")
		os.Exit(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error: Unable to read CSV file.")
		os.Exit(1)
	}

	// Find the target column index
	headers := records[0]
	columnIndex := -1
	for i, header := range headers {
		if header == columnName {
			columnIndex = i
			break
		}
	}
	if columnIndex == -1 {
		fmt.Println("Error: Specified column not found in CSV.")
		os.Exit(1)
	}

	wordCount := 0
	for _, row := range records[1:] {
		if columnIndex < len(row) {
			wordCount += countWords(preprocessText(row[columnIndex], config))
		}
	}
	return wordCount
}

// Process XLS files
func processXLS(filePath, sheetName, columnName string, config map[string][]string) int {
	if sheetName == "" || columnName == "" {
		fmt.Println("Error: Please specify both sheet and column for XLS files.")
		os.Exit(1)
	}

	file, err := excelize.OpenFile(filePath)
	if err != nil {
		fmt.Println("Error: Unable to open XLS file.")
		os.Exit(1)
	}

	rows, err := file.GetRows(sheetName)
	if err != nil {
		fmt.Println("Error: Unable to read sheet.")
		os.Exit(1)
	}
	if len(rows) == 0 {
		fmt.Println("Error: XLS sheet is empty.")
		os.Exit(1)
	}

	// Find the target column index
	headers := rows[0]
	columnIndex := -1
	for i, header := range headers {
		if header == columnName {
			columnIndex = i
			break
		}
	}
	if columnIndex == -1 {
		fmt.Println("Error: Specified column not found in XLS sheet.")
		os.Exit(1)
	}

	wordCount := 0
	for _, row := range rows[1:] {
		if columnIndex < len(row) {
			wordCount += countWords(preprocessText(row[columnIndex], config))
		}
	}
	return wordCount
}
