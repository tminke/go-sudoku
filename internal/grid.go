package internal

import (
	"fmt"
	"sync"
)

// Grid represents the current state of the Sudoku board maintianing
// the known and possible values for each cell.  It is accessed by
// Grid[row][col] coordinates.
type Grid struct {
	cells [9][9]*Cell
	mutex sync.RWMutex
}

// NewGrid returns an initialized Grid with all Cells unkown, suitable to
// start calling SetValue() on.
func NewGrid() *Grid {
	cells := [9][9]*Cell{}
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			cells[row][col] = NewCell()
		}
	}
	return &Grid{cells: cells}
}

// NewGridFromCsv returns a Grid initialized from the content in the
// specified CSV file, or an error if the format / content are invalid.
func NewGridFromCsv(csvFile string) (*Grid, error) {

	// Attempt To Parse The Specified Sudoku CSV File
	csvData, err := parseSudokuCsv(csvFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Sudoku CSV file '%s': err = %w", csvFile, err)
	}

	// Create a new starting Grid
	grid := NewGrid()

	// Initialze the Grid from the CSV data
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if csvData[row][col] > 0 {
				value := csvData[row][col]
				grid.SetValue(row, col, value)
			}
		}
	}

	// Return The Initialized Grid
	return grid, nil
}

// String returns a "box-drawing" string representing the current state of the
// Grid suitable for display.
func (g *Grid) String() string {

	// Constant Unicode Grid Border / Separators
	resetColor := "\033[0m"
	borderColor := "\033[34m" // Blue
	topBorder := borderColor + "\u250F\u2501\u2501\u2501\u252F\u2501\u2501\u2501\u252F\u2501\u2501\u2501\u2533\u2501\u2501\u2501\u252F\u2501\u2501\u2501\u252F\u2501\u2501\u2501\u2533\u2501\u2501\u2501\u252F\u2501\u2501\u2501\u252F\u2501\u2501\u2501\u2513\n" + resetColor
	btmBorder := borderColor + "\u2517\u2501\u2501\u2501\u2537\u2501\u2501\u2501\u2537\u2501\u2501\u2501\u253B\u2501\u2501\u2501\u2537\u2501\u2501\u2501\u2537\u2501\u2501\u2501\u253B\u2501\u2501\u2501\u2537\u2501\u2501\u2501\u2537\u2501\u2501\u2501\u251B\n" + resetColor
	lightSeparator := borderColor + "\u2520\u2500\u2500\u2500\u253c\u2500\u2500\u2500\u253c\u2500\u2500\u2500\u2542\u2500\u2500\u2500\u253c\u2500\u2500\u2500\u253c\u2500\u2500\u2500\u2542\u2500\u2500\u2500\u253c\u2500\u2500\u2500\u253c\u2500\u2500\u2500\u2528\n" + resetColor
	heavySeparator := borderColor + "\u2523\u2501\u2501\u2501\u253F\u2501\u2501\u2501\u253F\u2501\u2501\u2501\u254B\u2501\u2501\u2501\u253F\u2501\u2501\u2501\u253F\u2501\u2501\u2501\u254B\u2501\u2501\u2501\u253F\u2501\u2501\u2501\u253F\u2501\u2501\u2501\u252B\n" + resetColor

	// Start with the Top border ; )
	gridString := topBorder

	// Loop over all Grid Rows appending content and separators
	for row := 0; row < 9; row++ {

		// Format the Row's data with appropriate Cell dividors (light, heavy)
		rowString := borderColor + "\u2503" + resetColor // Heavy Vertical Bar
		for col := 0; col < 9; col++ {
			verticalSeparator := borderColor + "\u2502" + resetColor // Light Vertical Bar
			if (col+1)%3 == 0 {
				verticalSeparator = borderColor + "\u2503" + resetColor // Heavy Vertical Bar
			}
			rowString = fmt.Sprintf("%s %s %s", rowString, g.cells[row][col].GetValueString(), verticalSeparator)
		}

		// Determine the appropriate Separator line for each Row (light, heavy, none)
		horizontalSeparator := lightSeparator
		if row == 8 {
			horizontalSeparator = ""
		} else if (row+1)%3 == 0 {
			horizontalSeparator = heavySeparator
		}

		// Append the Row and subsequent Separator line
		gridString = fmt.Sprintf("%s%s\n%s", gridString, rowString, horizontalSeparator)
	}

	// Finally add the bottom border
	gridString = gridString + btmBorder

	// Return the result!
	return gridString
}

// GetCell returns the Cell at the specified row/col.
func (g *Grid) GetCell(row int, col int) *Cell {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	return g.cells[row][col]
}

// SetValue marks a Cell with the spcified value and updates the possible
// values of all Cells in the associated Rows, Column, and Group.
func (g *Grid) SetValue(row int, col int, value int) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.cells[row][col].SetValue(value)
	g.eliminateValueFromRow(row, value)
	g.eliminateValueFromCol(col, value)
	g.eliminateValueFromGroup(row, col, value)
}

func (g *Grid) eliminateValueFromRow(row int, value int) {
	for col := 0; col < 9; col++ {
		g.cells[row][col].EliminateValue(value)
	}
}

func (g *Grid) eliminateValueFromCol(col int, value int) {
	for row := 0; row < 9; row++ {
		g.cells[row][col].EliminateValue(value)
	}
}

func (g *Grid) eliminateValueFromGroup(row int, col int, value int) {
	groupRow := int(row/3) * 3 // Starting Row Index Of Group
	groupCol := int(col/3) * 3 // Starting Col Index Of Group
	for row := groupRow; row < groupRow+3; row++ {
		for col := groupCol; col < groupCol+3; col++ {
			g.cells[row][col].EliminateValue(value)
		}
	}
}
