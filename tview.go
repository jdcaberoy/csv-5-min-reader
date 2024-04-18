package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rivo/tview"
)

func startGui(file *os.File, err error) {
	app := tview.NewApplication()
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	resultsBox := tview.NewTextView().SetScrollable(true)
	exitForm := tview.NewForm().
		AddButton("EXIT", func() {
			file.Close()
			app.Stop()
			return
		})

	// set msg for problems on file open
	if (err != nil) || (file == nil) {
		resultsBox.SetText("NO DATA.CSV FOUND\nCHECK FILES AND RESTART")
		flex.AddItem(resultsBox, 0, 1, true)
		flex.AddItem(exitForm, 0, 1, true)
		if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
			panic(err)
		}
		return
	}

	form := tview.NewForm()

	form.AddInputField("month", "", 20, tview.InputFieldInteger, nil)
	form.AddInputField("day", "", 20, tview.InputFieldInteger, nil)
	form.AddInputField("year", "2024", 20, tview.InputFieldInteger, nil)
	form.AddButton("Extract", func() {
		tl := setDate(form)
		if int(tl.startTime.Month()) > 6 {
			app.Stop()
			return
		}
		values := readData(file, tl)
		textData, aveValues := averageValues(values, tl)
		resultsBox.SetText(fmt.Sprint(strings.Join(textData, "")))
		if err := saveToExcel(aveValues, tl.startTime); err != nil {
			resultsBox.SetText(fmt.Sprintf("Error in saving excel: %v", err))
		}

		// Reset the file position to the beginning
		_, err := file.Seek(0, 0)
		if err != nil {
			resultsBox.SetText(fmt.Sprintf("Error resetting file position: %v", err))
		}

		// resultsBox.SetBorder(true)
	})
	flex.AddItem(form, 0, 1, true)
	flex.AddItem(resultsBox, 0, 1, false)
	flex.AddItem(exitForm, 0, 1, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

func setDate(form *tview.Form) timeLimits {
	m := form.GetFormItem(0).(*tview.InputField).GetText()
	d := form.GetFormItem(1).(*tview.InputField).GetText()
	y := form.GetFormItem(2).(*tview.InputField).GetText()
	inputmonth, _ := strconv.Atoi(m)
	inputday, _ := strconv.Atoi(d)
	inputyear, _ := strconv.Atoi(y)
	startTime := time.Date(inputyear, time.Month(inputmonth), inputday, 5, 0, 0, 0, time.UTC)
	endTime := time.Date(2024, time.Month(inputmonth), inputday, 19, 0, 0, 0, time.UTC)
	return timeLimits{
		startTime: startTime,
		endTime:   endTime,
	}
}
