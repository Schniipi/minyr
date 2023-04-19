package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/Schniipi/minyr/yr"
)

func main() {
        // Brukeren skriver "minyr" og med annen input blir programmet avsluttet
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Type 'minyr' to continue or anything else to exit:")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "minyr" {
		fmt.Println("Invalid input. Exiting.")
		return
	}

	for {
        // Brukeren kan velge mellom å konvertere, finne average eller å avslutte
		fmt.Println("Choose an option: 'convert', 'average', or 'exit':")
		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "convert":
			err := yr.ConvertTemp()
			if err != nil {
				fmt.Printf("Error converting temperatures: %v\n", err)
			} else {
				fmt.Println("Temperatures successfully converted.")
			}
		case "average":
			fmt.Println("Choose a unit: 'c' for Celsius or 'f' for Fahrenheit:")
			unit, _ := reader.ReadString('\n')
			unit = strings.TrimSpace(unit)

			average, err := yr.AverageTemp(unit)
			if err != nil {
				fmt.Printf("Error calculating average temperature: %v\n", err)
			} else {
				fmt.Printf("The average temperature is %.1f degrees %s\n", average, strings.ToUpper(unit))
			}
		case "exit":
			fmt.Println("Exiting the program.")
			return
		default:
			fmt.Println("Invalid option. Try again.")
		}
	}
}

