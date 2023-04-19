package yr

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
        "github.com/Schniipi/funtemps/tree/main/conv"
)

// ConvertCelsiusToFahrenheit converts Celsius to Fahrenheit
func ConvertCelsiusToFahrenheit(celsius float64) float64 {
	fahrenheit := conv.CelsiusToFahrenheit(celsius)
        return fahrenheit
}

// GenerateConvertedFile converts the input CSV file with Celsius temperature data
// to an output CSV file with Fahrenheit temperature data.
func ConvertTemp() error {
	const outputFileName = "kjevik-temp-fahr-20220318-20230318.csv"
	const inputFileName = "kjevik-temp-celsius-20220318-20230318.csv"

	// Check if the output file already exists
	if _, err := os.Stat(outputFileName); !os.IsNotExist(err) {
		var regenerate string
		fmt.Print("Output file already exists. Regenerate? (y/n): ")
		fmt.Scanln(&regenerate)
		if regenerate != "y" && regenerate != "Y" {
			fmt.Println("Exiting without generating new file.")
			return nil
		}
	}

	// Open the input CSV file
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		return fmt.Errorf("error opening input file: %w", err)
	}
	defer inputFile.Close()

	// Create a scanner for reading the input CSV file
	inputScanner := bufio.NewScanner(inputFile)

	// Create the output CSV file
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer outputFile.Close()

	// Create a CSV writer for writing to the output CSV file
	csvWriter := csv.NewWriter(outputFile)
	defer csvWriter.Flush()

	// Write the header row to the output CSV file
	if inputScanner.Scan() {
		firstLine := inputScanner.Text()
		if err = csvWriter.Write(strings.Split(firstLine, ";")); err != nil {
			return fmt.Errorf("error writing first line: %w", err)
		}
	}

	lineNo := 2
	// Loop through each line of the input CSV file
	for inputScanner.Scan() {
		// Split the line into fields
		fields := strings.Split(inputScanner.Text(), ";")

		// Break the loop if the maximum number of lines is exceeded
		if lineNo > 16755 {
			break
		}

		// Validate the input format
		if len(fields) != 4 {
			fmt.Printf("Error on line %d: Invalid input format.\n", lineNo)
			continue
		}

		// Check if the temperature field is empty
		tempField := fields[3]
		if tempField == "" {
			fmt.Printf("Error on line %d: Temperature field is empty.\n", lineNo)
			continue
		}

		// Parse the temperature value as a float64
		temperature, err := strconv.ParseFloat(tempField, 64)
		if err != nil {
			fmt.Printf("Error on line %d: %v\n", lineNo, err)
			continue
		}

		// Check if the temperature value is a valid float64
		if math.IsNaN(temperature) {
			fmt.Printf("Error on line %d: Temperature is not a valid float64 value.\n", lineNo)
			continue
		}



	fields[3] = fmt.Sprintf("%.1f", fahrenheit)
	if err := csvWriter.Write(fields); err != nil {
		return fmt.Errorf("error writing line to output file: %w", err)
	}

	lineNo++
}

dataText := "Data er basert på gyldig data (as of 18.03.2023) (CC BY 4.0) fra Meteorologisk institutt (MET); endringen er gjort av Oliver Aaron Berg Johnston"
err = csvWriter.Write([]string{dataText})
if err != nil {
	return fmt.Errorf("error writing data text to output file: %w", err)
}

func AverageTemp() (float64, error) {

// AverageTemperature calculates the average temperature in the provided CSV file.

var filename string
var tempColumn int
var delimeter rune


// Set the input file and temperature column based on the unit parameter
if unit == "c" {
	filename = "kjevik-temp-celsius-20220318-20230318.csv"
	tempColumn = 3
	delimeter = ';'
} else if unit == "f" {
	filename = "kjevik-temp-fahr-20220318-20230318.csv"
	tempColumn = 3
	delimeter = ';'
} else {
	return 0, errors.New("invalid temperature unit")
}

// Open the input CSV file
file, err := os.Open(filename)
if err != nil {
	return 0, fmt.Errorf("error opening input file: %w", err)
}
defer file.Close()

// Create a CSV reader to read from the input file
csvReader := csv.NewReader(file)
csvReader.Comma = delimeter

// Initialize variables to store the sum of temperatures and the count of lines
tempSum := 0.0
lineCount := 0

// Skip the header line
_, err = csvReader.Read()
if err != nil {
	return 0, fmt.Errorf("error reading header line: %w", err)
}

// Loop through each line of the input CSV file
for {
	fields, err := csvReader.Read()
	if err == io.EOF {
		break
	}
	if err != nil {
		return 0, fmt.Errorf("error reading line: %w", err)
	}

	// Extract the temperature value from the current line
	tempStr := fields[tempColumn]

	// Convert the temperature value to a float64
	temp, err := strconv.ParseFloat(tempStr, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing temperature value: %w", err)
	}

	// Add the temperature value to the sum and increment the line count
	tempSum += temp
	lineCount++
}

// Calculate the average temperature
averageTemp := tempSum / float64(lineCount)

return averageTemp, nil
}
