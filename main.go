package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// CLI Arguments
	filePath := flag.String("file", "", "Path to the input file")
	textInput := flag.String("text", "", "String input for word count")
	outputFormat := flag.String("output", "txt", "Output format: txt, json, csv")
	configPath := flag.String("config", "", "Path to configuration file")
	printToTerminal := flag.Bool("print", false, "Print word count to the terminal")
	sheetName := flag.String("sheet", "", "Sheet name (for .xls files)")
	columnName := flag.String("column", "", "Column name (for .csv/.xls files)")

	// Parse initial flags
	flag.Parse()

	// If no file or text was provided, prompt the user
	if *filePath == "" && *textInput == "" {
		fmt.Println("No file or text input provided. Please enter your arguments below.")
		fmt.Println("For example:")
		fmt.Println("--file test_files/sample.txt --print")
		fmt.Println("or")
		fmt.Println("--text \"Hello world\" --print")

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter arguments: ")
		argLine, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			os.Exit(1)
		}

		argLine = strings.TrimSpace(argLine)
		if argLine == "" {
			fmt.Println("No input provided. Exiting.")
			os.Exit(1)
		}

		// Split arguments by space and re-parse
		args := strings.Split(argLine, " ")

		// Create a new FlagSet to parse the new arguments
		fs := flag.NewFlagSet("interactive", flag.ExitOnError)
		filePath = fs.String("file", "", "Path to the input file")
		textInput = fs.String("text", "", "String input for word count")
		outputFormat = fs.String("output", "txt", "Output format: txt, json, csv")
		configPath = fs.String("config", "", "Path to configuration file")
		printToTerminal = fs.Bool("print", false, "Print word count to the terminal")
		sheetName = fs.String("sheet", "", "Sheet name (for .xls files)")
		columnName = fs.String("column", "", "Column name (for .csv/.xls files)")

		// Parse the new arguments
		fs.Parse(args)

		// After re-parsing, continue with logic below
	}

	// After interactive mode, proceed as before.
	var config map[string][]string
	if *configPath != "" {
		config = loadConfig(*configPath)
	}

	var wordCount int
	if *filePath != "" {
		wordCount = processFile(*filePath, *sheetName, *columnName, config)
	} else if *textInput != "" {
		wordCount = countWords(preprocessText(*textInput, config))
	} else {
		fmt.Println("Error: Please provide a file or text input.")
		os.Exit(1)
	}

	if *printToTerminal {
		fmt.Printf("Word Count: %d\n", wordCount)
	} else {
		switch *outputFormat {
		case "txt":
			writeOutputTxt(wordCount, "output.txt")
		case "json":
			writeOutputJSON(map[string]int{"word_count": wordCount}, "output.json")
		case "csv":
			writeOutputCSV(map[string]int{"word_count": wordCount}, "output.csv")
		default:
			fmt.Println("Error: Invalid output format. Choose txt, json, or csv.")
		}
	}
}
