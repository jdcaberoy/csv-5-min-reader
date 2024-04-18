package main

import (
	"os"
)

func main() {
	file, err := os.Open("data.CSV")
	startGui(file, err)
}
