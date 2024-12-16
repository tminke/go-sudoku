package internal

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGrid(t *testing.T) {
	grid := NewGrid()
	assert.NotNil(t, grid)
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			cell := grid.GetCell(row, col)
			assert.NotNil(t, cell)
			assert.Equal(t, 0, cell.GetValue())
		}
	}
}

func TestNewGridFromCsv(t *testing.T) {

	// Create a temporary file in current directory
	file, err := os.CreateTemp("", "new-grid-from-csv-*.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	// Write the content to the file
	fileContent := testGridCsv()
	for _, content := range fileContent {
		_, err = file.WriteString(content)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Perform The Test
	grid, err := NewGridFromCsv(file.Name())

	// Verify The Results
	assert.Equal(t, testGrid(), grid)
	assert.Nil(t, err)
}

func TestGridString(t *testing.T) {
	grid := testGrid()
	t.Logf("Grid:\n\n%s\n", grid.String())
	// Visual verification only ; )
}

func TestGridGetCell(t *testing.T) {
	grid := testGrid()
	assert.Equal(t, 2, grid.GetCell(0, 0).GetValue()) // Spot check values in test Grid
	assert.Equal(t, 0, grid.GetCell(0, 8).GetValue())
	assert.Equal(t, 0, grid.GetCell(8, 0).GetValue())
	assert.Equal(t, 7, grid.GetCell(8, 8).GetValue())
	assert.Equal(t, 5, grid.GetCell(3, 3).GetValue())
	assert.Equal(t, 6, grid.GetCell(6, 6).GetValue())
}

func TestGrid_SetValue(t *testing.T) {

	// Define The TestCases
	testCases := map[string]struct {
		row   int
		col   int
		value int
	}{
		"Step 1": {row: 3, col: 5, value: 1},
		"Step 2": {row: 6, col: 5, value: 5},
		"Step 3": {row: 6, col: 7, value: 1},
		"Step 4": {row: 3, col: 4, value: 7},
		"Step 5": {row: 0, col: 4, value: 5},
	}

	// Execute The TestCases
	for testCaseName, testCase := range testCases {
		t.Run(testCaseName, func(t *testing.T) {

			// Create a Grid to test
			grid := testGrid()

			// Perform The Test
			grid.SetValue(testCase.row, testCase.col, testCase.value)

			// Verify Cell State (Value is set and possible values cleared)
			assert.Equal(t, testCase.value, grid.GetCell(testCase.row, testCase.col).GetValue())
			assert.Equal(t, []int{}, grid.GetCell(testCase.row, testCase.col).GetPossibleValues())

			// Verify Row State (Value no longer possible in Row)
			for col := 0; col < 9; col++ {
				assert.NotContains(t, grid.GetCell(testCase.row, col).GetPossibleValues(), testCase.value)
			}

			// Verify Column State (Value no longer possible in Column)
			for row := 0; row < 9; row++ {
				assert.NotContains(t, grid.GetCell(row, testCase.col).GetPossibleValues(), testCase.value)
			}

			// Verify Group State (Value no longer possible in Group)
			groupRow := int(testCase.row/3) * 3
			groupCol := int(testCase.col/3) * 3
			for groupCellRow := groupRow; groupCellRow < groupRow+3; groupCellRow++ {
				for groupCellCol := groupCol; groupCellCol < groupCol+3; groupCellCol++ {
					assert.NotContains(t, grid.GetCell(groupCellRow, groupCellCol).GetPossibleValues(), testCase.value)
				}
			}
		})
	}
}

// testGrid returns a sample grid version of the
// samples/hard.csv puzzle for testing ; )
func testGrid() *Grid {

	// Create A New Grid
	grid := NewGrid()

	// Row 1
	grid.SetValue(0, 0, 2)
	grid.SetValue(0, 1, 6)
	grid.SetValue(0, 3, 1)
	grid.SetValue(0, 5, 4)

	// Row 2
	grid.SetValue(1, 6, 5)

	// Row 3
	grid.SetValue(2, 1, 8)
	grid.SetValue(2, 5, 7)
	grid.SetValue(2, 7, 2)
	grid.SetValue(2, 8, 9)

	// Row 4
	grid.SetValue(3, 0, 6)
	grid.SetValue(3, 3, 5)
	grid.SetValue(3, 7, 3)
	grid.SetValue(3, 8, 2)

	// Row 5
	grid.SetValue(4, 3, 9)
	grid.SetValue(4, 4, 6)
	grid.SetValue(4, 5, 3)
	grid.SetValue(4, 7, 4)

	// Row 6
	grid.SetValue(5, 0, 3)
	grid.SetValue(5, 2, 7)
	grid.SetValue(5, 3, 8)
	grid.SetValue(5, 4, 4)
	grid.SetValue(5, 5, 2)
	grid.SetValue(5, 6, 1)

	// Row 7
	grid.SetValue(6, 2, 8)
	grid.SetValue(6, 4, 9)
	grid.SetValue(6, 6, 6)

	// Row 8
	grid.SetValue(7, 1, 3)
	grid.SetValue(7, 2, 5)

	// Row 9
	grid.SetValue(8, 6, 2)
	grid.SetValue(8, 8, 7)

	// Return The Grid
	return grid
}

// testGridCsv returns an array of strings representing the equivalent CSV
// content of the Grid from testGrid().
func testGridCsv() [9]string {
	grid := testGrid()
	csv := [9]string{}
	for row := 0; row < 9; row++ {
		rowString := ""
		for col := 0; col < 9; col++ {
			valueString := "-"
			value := grid.GetCell(row, col).GetValue()
			if value >= 1 && value <= 9 {
				valueString = strconv.Itoa(value)
			}
			rowString = rowString + valueString
			if col < 8 {
				rowString = rowString + ","
			}
		}
		csv[row] = rowString + "\n"
	}
	return csv
}
