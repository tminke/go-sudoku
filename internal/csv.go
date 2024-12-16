package internal

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// parseSudokuCsv returns the int values parsed from the specified CSV file,
// or an error should any formatting or content problems exist.
func parseSudokuCsv(csvFile string) ([9][9]int, error) {

	// Create an empty data set
	intData := [9][9]int{}

	// Attempt to open the CSV File
	file, err := os.Open(csvFile)
	if err != nil {
		return intData, err
	}
	defer file.Close()

	// Attempt to read all records (unprotected)
	csvReader := csv.NewReader(file)
	stringData, err := csvReader.ReadAll()
	if err != nil {
		return intData, err
	}

	// Perform some basic validation on the String data
	if len(stringData) != 9 {
		return intData, csvError(fmt.Sprintf("exactly nine rows expected, encountered %d", len(stringData)))
	}

	// Convert String data to Ints
	for row, rowData := range stringData {
		if len(rowData) != 9 {
			return intData, csvError(fmt.Sprintf("exactly nine cols expected, encountered %d on row %d", len(rowData), row))
		}
		for col, stringValue := range rowData {
			stringValue = strings.TrimSpace(stringValue)
			if stringValue == "-" {
				intData[row][col] = 0
			} else {
				intValue, err := strconv.Atoi(stringValue)
				if err != nil {
					return intData, csvError(fmt.Sprintf("encountered unsupported value '%s' must be one of -,1,2,3,4,5,6,7,8,9: err=%+v", stringValue, err))
				}
				if intValue < 1 || intValue > 9 {
					return intData, csvError(fmt.Sprintf("encountered invalid value '%d' must be one of 1,2,3,4,5,6,7,8,9", intValue))
				}
				intData[row][col] = intValue
			}
		}
	}

	// Return Success
	return intData, nil
}

func csvError(reason string) error {
	return fmt.Errorf("sudoku CSV file format error: %s", reason)
}
