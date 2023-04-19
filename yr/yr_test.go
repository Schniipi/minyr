package yr

import (
	"bufio"
	"encoding/csv"
        "math"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestLineCount(t *testing.T) {
	const outputFileName = "kjevik-temp-fahr-20220318-20230318.csv"

	// Open the output CSV file
	file, err := os.Open(outputFileName)
	if err != nil {
		t.Fatalf("error opening output file: %v", err)
	}
	defer file.Close()

	// Create a scanner for reading the output CSV file
	scanner := bufio.NewScanner(file)

	// Check the number of lines in the output file
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	if lineCount != 16756 {
		t.Errorf("unexpected number of lines in the output file: got %d, want 16756", lineCount)
	}
}

func TestTemperatureConversion(t *testing.T) {
	testCases := []struct {
		inputLine  string
		outputLine string
	}{
		{"Kjevik;SN39040;18.03.2022 01:50;6", "Kjevik;SN39040;18.03.2022 01:50;42.8"},
		{"Kjevik;SN39040;07.03.2023 18:20;0", "Kjevik;SN39040;07.03.2023 18:20;32.0"},
		{"Kjevik;SN39040;08.03.2023 02:20;-11", "Kjevik;SN39040;08.03.2023 02:20;12.2"},
	}

	for _, tc := range testCases {
		inputFields := strings.Split(tc.inputLine, ";")
		inputCelsius, _ := strconv.ParseFloat(inputFields[3], 64)
		inputFahrenheit := yr.ConvertCelsiusToFahrenheit(inputCelsius)

		outputFields := strings.Split(tc.outputLine, ";")
		outputFahrenheit, _ := strconv.ParseFloat(outputFields[3], 64)

		if inputFahrenheit != outputFahrenheit {
			t.Errorf("unexpected temperature conversion: got %s, want %s", inputFields[3], outputFields[3])
		}
	}
}

func TestDataText(t *testing.T) {
	const outputFileName = "kjevik-temp-fahr-20220318-20230318.csv"

	// Open the output CSV file
	file, err := os.Open(outputFileName)
	if err != nil {
		t.Fatalf("error opening output file: %v", err)
	}
	defer file.Close()

	// Create a CSV reader to read from the output file
	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'

	// Read to the last line
	var lastLine []string
	for {
		line, err := csvReader.Read()
		if err != nil {
			break
		}
		lastLine = line
	}

	// Test data text line
	wantDataText := "Data er basert på gyldig data (as of 18.03.2023) (CC BY 4.0) fra Meteorologisk institutt (MET); endringen er gjort av Oliver Aaron Berg Johnston"
	gotDataText := lastLine[0]
	if gotDataText != wantDataText {
		t.Errorf("unexpected data text: got %q, want %q", gotDataText, wantDataText)
	}
}

func TestAverageTemp(t *testing.T) {
	wantAvg := 8.56
	avgTemp, err := yr.AverageTemp("c")
	if err != nil {
		t.Fatalf("error calculating average temperature: %v", err)
	}

	// Compare the calculated average temperature with the expected value
	if math.Abs(avgTemp-wantAvg) > 1e-2 {
		t.Errorf("unexpected average temperature: got %.2f, want %.2f", avgTemp, wantAvg)
	}
}


