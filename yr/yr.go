package yr

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
        "github.com/Schniipi/funtemps/conv"
)

// ConvertCelsiusToFahrenheit converts Celsius to Fahrenheit
func ConvertCelsiusToFahrenheit(celsius float64) float64 {
	fahrenheit := conv.CelsiusToFahrenheit(celsius)
        return fahrenheit
}

// ConvertTemp converts the input CSV file with Celsius temperature data
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
		temperature, err := parseFloatWithComma(tempField)
		if err != nil {
			fmt.Printf("Error on line %d: %v\n", lineNo, err)
			continue
		}

		// Check if the temperature value is a valid float64
		if math.IsNaN(temperature) {
			fmt.Printf("Error on line %d: Temperature is not a valid float64 value.\n", lineNo)
			continue
		}

                fahrenheit:= ConvertCelsiusToFahrenheit(temperature)



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
        return nil
}

func parseFloatWithComma(s string) (float64, error) {
    s = strings.Replace(s, ",", ".", -1)
    return strconv.ParseFloat(s, 64)
}


func AverageTemp(unit string) (float64, error) {
	// Set the appropriate filename, temperature column, and delimiter based on the temperature unit.
	var filename string
	var tempColumn int
	var delimeter rune

	if unit == "c" {
		filename = "kjevik-temp-celsius-20220318-20230318.csv"
		tempColumn = 3
		delimeter = ';'
	} else if unit == "f" {
		filename = "kjevik-temp-fahr-20220318-20230318.csv"
		tempColumn = 3
		delimeter = ';'
	} else {
		return 0, fmt.Errorf("invalid temperature unit: %s", unit)
	}

	// Open the file for reading and ensure it is closed when the function exits.
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Create a new CSV reader with the correct delimiter.
	reader := csv.NewReader(file)
	reader.Comma = delimeter
	// Allow a variable number of fields per record.
	reader.FieldsPerRecord = -1

	var total float64
	var line int

	// Loop through the lines in the CSV file.
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, err
		}

                line++

		// Skip the first line (header).
		if line < 2 || line > 16755 {
			continue
		}

		// Check if the record has the required number of fields.
		if len(record) <= tempColumn {
			return 0, fmt.Errorf("invalid data in file %s at line %d, column %d", filename, line, tempColumn)
		}

		// Parse the temperature value and report an error if it fails.
		temp, err := strconv.ParseFloat(record[tempColumn], 64)
		if err != nil {
			return 0, err
		}

		// Add the temperature value to the total and increment the count.
		total += temp
	}

	// If there's no temperature data, report an error.
	if line == 1 {
		return 0, fmt.Errorf("no temperature data found in file %s", filename)
	}

	// Calculate the average temperature by dividing the total by the number of temperature values.
	return total / float64(line-1), nil
}

