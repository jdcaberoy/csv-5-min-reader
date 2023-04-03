package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	// Open CSV file
	file, err := os.Open("data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var inputmonth, inputday, inputyear int
	today := time.Now()
	fmt.Println("input month:")
	fmt.Scanln(&inputmonth)
	fmt.Println("input day:")
	fmt.Scanln(&inputday)
	fmt.Println("input year:")
	fmt.Scanln(&inputyear)
	if inputmonth == 0{
		today.Month()
	}
	if inputday == 0 {
		today.Day()
	}
	if inputyear ==0 {
		today.Year()
	}


	// Initialize start and end times
	startTime := time.Date(2023, 4, 3, 5, 0, 0, 0, time.UTC)
	endTime := time.Date(2023, 4, 3, 19, 0, 0, 0, time.UTC)

	// Initialize map to store values for each 5-minute interval
	values := make(map[time.Time][]float64)

	// Read CSV file
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.Comment = '#'
	reader.TrimLeadingSpace = true
	if _, err := reader.Read(); err != nil {
		log.Fatal(err)
	}

	for {
		// Read row
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		// Parse date
		dateString := row[2]
		date, err := time.Parse("02/01/2006 15:04:05", dateString)
		if err != nil {
			log.Fatal(err)
		}

		// Skip if date is outside of desired range
		if date.Before(startTime) || date.After(endTime) {
			continue
		}

		// Parse value
		value, err := strconv.ParseFloat(strings.TrimSpace(row[6]), 64)
		if err != nil {
			log.Fatal(err)
		}

		// Calculate interval start time
		intervalStart := time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), (date.Minute()/5)*5, 0, 0, date.Location())

		// Append value to interval
		values[intervalStart] = append(values[intervalStart], value)
	}

	// Calculate average for each interval
	for i := startTime; i.Before(endTime); i = i.Add(5 * time.Minute) {
		intervalStart := i
		intervalEnd := i.Add(5 * time.Minute)

		if intervalEnd.After(endTime) {
			intervalEnd = endTime
		}

		valueList := values[intervalStart]

		if len(valueList) == 0 {
			fmt.Printf("%s - %s = 0.00\n", intervalStart.Format("15:04:05"), intervalEnd.Format("15:04:05"))
			continue
		}

		sum := 0.0
		for _, value := range valueList {
			sum += value
		}
		avg := sum / float64(len(valueList))

		fmt.Printf("%s - %s = %.2f\n", intervalStart.Format("15:04:05"), intervalEnd.Format("15:04:05"), avg)
	}
}
