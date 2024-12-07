package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

// Write output as TXT
func writeOutputTxt(wordCount int, outputPath string) {
	err := os.WriteFile(outputPath, []byte(fmt.Sprintf("Word Count: %d\n", wordCount)), 0644)
	if err != nil {
		fmt.Println("Error: Unable to write TXT output.")
		os.Exit(1)
	}
}

// Write output as JSON
func writeOutputJSON(data map[string]int, outputPath string) {
	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error: Unable to create JSON file.")
		os.Exit(1)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		fmt.Println("Error: Unable to write JSON output.")
		os.Exit(1)
	}
}

// Write output as CSV
func writeOutputCSV(data map[string]int, outputPath string) {
	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error: Unable to create CSV file.")
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers and data
	writer.Write([]string{"Metric", "Value"})
	writer.Write([]string{"Word Count", fmt.Sprintf("%d", data["word_count"])})
}
