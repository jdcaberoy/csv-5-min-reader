package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
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

	reader := csv.NewReader(file)
	if _, err := reader.Read(); err != nil { // skip header row
		log.Fatal(err)
	}

	// Initialize Excel file
	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	f.SetActiveSheet(index)

	// Initialize time parameters
	startTime := time.Date(inputyear, time.Month(inputmonth), inputday, 5, 0, 0, 0, time.UTC)
	endTime := time.Date(inputyear, time.Month(inputmonth), inputday, 19, 0, 0, 0, time.UTC)
	interval := time.Minute * 5

	// Initialize row and column values for Excel
	row := 1
	col := 1

	// Initialize map to store values
	values := make(map[string]float64)

	// Loop through CSV file and store values in map
	for {
		rowData, err := reader.Read()
		if err != nil {
			break
		}

		t, err := time.Parse("01/02/2006 15:04:05", rowData[2])
		if err != nil {
			log.Fatal(err)
		}

		if t.Before(startTime) || t.After(endTime) {
			continue
		}

		value, err := strconv.ParseFloat(rowData[6], 64)
		if err != nil {
			log.Fatal(err)
		}

		key := t.Truncate(interval).Format("15:04:05")
		values[key] += value
	}

	// Loop through time interval and insert values into Excel
	for currentTime := startTime; currentTime.Before(endTime); currentTime = currentTime.Add(interval) {
		// Calculate end time for current interval
		// endTime := currentTime.Add(interval)

		// Get average value for current interval
		key := currentTime.Format("15:04:05")
		averageValue := values[key] / float64(interval/time.Minute)

		// Add value to Excel sheet
		cell := fmt.Sprintf("%s%d", excelize.ToAlphaString(col), row)
		f.SetCellValue("Sheet1", cell, averageValue)

		// Move to next column or row
		if col < 10 {
			col++
		} else {
			col = 1
			row++
		}
	}

	// Save Excel file
	if err := f.SaveAs("output.xlsx"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Excel file created successfully.")
	
	for {
		fmt.Println("Press CTRL+C to exit...")
		fmt.Scanln()
	}
}
