package main

import (
	"fmt"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func saveToExcel(aveValues []float64, t time.Time) error {
	f := excelize.NewFile()

	// Set internal interval for excel
	eInterval := 0

	// Sample input for excel. ignore
	f.SetCellValue("Sheet1", "A30", "sample")

	// set row for 14 hours
	for row := 1; row <= 14; row++ {
		// set column for 6 - 5min intervals
		for col := 'A'; col <= 'L'; col++ {
			cell := fmt.Sprintf("%c%d", col, row)
			cellwrite := fmt.Sprintf("%.2f", aveValues[eInterval])
			f.SetCellValue("Sheet1", cell, cellwrite)
			eInterval++
		}
	}

	fileTitle := fmt.Sprintf("5minsoutput %v-%02d-%02d.xlsx", t.Year(), int(t.Month()), t.Day())

	err := f.SaveAs(fileTitle)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
