package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	// Prompt user to input specific date
	var inputDate string
	fmt.Print("Enter date (dd/mm/yyyy): ")
	fmt.Scanln(&inputDate)
	if inputDate == "" {
		inputDate = time.Now().Format("02/01/2006")
	}

	// Open CSV file and read data
	csvFile, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	headers, err := csvReader.Read()
	if err != nil {
		panic(err)
	}

	// Find index of "Date" and "Value" columns
	dateIndex := -1
	valueIndex := -1
	for i, header := range headers {
		switch header {
		case "Date":
			dateIndex = i
		case "Value":
			valueIndex = i
		}
	}
	if dateIndex == -1 || valueIndex == -1 {
		panic("Missing Date or Value column in CSV file")
	}

	// Define start and end times for 5-minute intervals
	startTime := "05:00:00.000"
	endTime := "19:00:00.000"

	// Initialize variables to store total and count for each 5-minute interval
	intervalTotal := 0.0
	intervalCount := 0

	// Loop through each row in the CSV file
	for {
		record, err := csvReader.Read()
		if err != nil {
			break
		}

		// Parse the row's date string into a time.Time value
		rowTime, err := time.Parse("02/01/2006 15:04:05", record[dateIndex])
		if err != nil {
			panic(err)
		}

		// Check if the row's date matches the input date
		if inputDate == rowTime.Format("02/01/2006") {

			// Extract the time component from the row's date string
			rowTimeString := rowTime.Format("15:04:05.000")

			// Check if the row's time is within the 5-minute interval
			if startTime <= rowTimeString && rowTimeString <= endTime {

				// Add the row's value to the total and increment the count
				value, err := strconv.ParseFloat(record[valueIndex], 64)
				if err != nil {
					panic(err)
				}
				intervalTotal += value
				intervalCount++

			} else if intervalCount > 0 {
				// If the row's time is outside the 5-minute interval, calculate the average and print it
				intervalAverage := intervalTotal / float64(intervalCount)
				fmt.Printf("Average value for %s-%s: %.2f\n", startTime, endTime, intervalAverage)
				intervalTotal = 0.0
				intervalCount = 0
			}

			// Update start time for next interval
			startTime = rowTimeString
		}
	}

	// Calculate the average for the last 5-minute interval (if any)
	if intervalCount > 0 {
		intervalAverage := intervalTotal / float64(intervalCount)
		fmt.Printf("Average value for %s-%s: %.2f\n", startTime, endTime, intervalAverage)
	}
}