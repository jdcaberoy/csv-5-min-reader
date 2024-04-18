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

type timeLimits struct {
	startTime time.Time
	endTime   time.Time
}

func readData(file *os.File, tl timeLimits) map[time.Time][]float64 {
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
		// Read row headers
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
		if date.Before(tl.startTime) || date.After(tl.endTime) {
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
	return values
}

func averageValues(values map[time.Time][]float64, tl timeLimits) ([]string, []float64) {
	var eValues []float64

	var aveDataList []string

	for i := tl.startTime; i.Before(tl.endTime); i = i.Add(5 * time.Minute) {
		intervalStart := i
		intervalEnd := i.Add(5 * time.Minute)

		// fmt.Printf("DEBUG01: %v - %v\n", intervalStart, intervalEnd)

		if intervalEnd.After(tl.endTime) {
			intervalEnd = tl.endTime
		}

		valueList := values[intervalStart]

		// fmt.Printf("DEBUG02: %v len:%v\n", valueList, len(valueList))

		if len(valueList) == 0 {
			// fmt.Printf("%s - %s = 0.00\n", intervalStart.Format("15:04:05"), intervalEnd.Format("15:04:05"))
			aveData := fmt.Sprintf("%s - %s = 0.00\n", intervalStart.Format("15:04:05"), intervalEnd.Format("15:04:05"))
			aveDataList = append(aveDataList, aveData)
			eValues = append(eValues, 0.00)
			continue
		}

		sum := 0.0
		for _, value := range valueList {
			sum += value
		}
		avg := sum / float64(len(valueList))
		eValues = append(eValues, avg)

		fmt.Printf("%s - %s = %.2f\n", intervalStart.Format("15:04:05"), intervalEnd.Format("15:04:05"), avg)
		aveData := fmt.Sprintf("%s - %s = %.2f\n", intervalStart.Format("15:04:05"), intervalEnd.Format("15:04:05"), avg)
		aveDataList = append(aveDataList, aveData)
	}
	return aveDataList, eValues
}
