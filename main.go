package main

import (
	"fmt"
	"log"

	gnss "github.com/mothergoose31/GNNS-GO/GNSS"
)

func main() {
	// Example of parsing a RINEX file
	filename := "./brdc2050.24g"

	// Parse the RINEX file
	header, ephemerides, err := gnss.ParseRINEXFile(filename)
	if err != nil {
		log.Fatalf("Error parsing RINEX file: %v", err)
	}

	// Print summary of parsed data
	fmt.Printf("Successfully parsed RINEX file: %s\n", filename)
	fmt.Printf("Number of ephemerides: %d\n", len(ephemerides))

	// Print header details
	fmt.Println("\nRINEX Header Details:")
	gnss.PrintRINEXHeader(*header)

	// Print first few ephemerides as example
	fmt.Println("\nFirst 5 Ephemerides:")
	for i, eph := range ephemerides {
		if i >= 5 {
			break
		}
		fmt.Println("====================================")
		gnss.PrintRINEXEphemeris(eph)
	}
}
