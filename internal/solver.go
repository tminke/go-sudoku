package internal

import (
	"log"
	"time"
)

const MaxIterations = 100 // Maximum number of passes through the algorithm

// Solver contains the basic state used when solving a Grid.
type Solver struct {
	maxIterations int
	verbose       bool
}

// NewSolver returns a Solver with the specified configuration.
func NewSolver(maxIterations int, verbose bool) *Solver {
	return &Solver{
		maxIterations: maxIterations,
		verbose:       verbose,
	}
}

// Solve does an in-place update to the specified Grid by setting values and
// iterating until complete or max iterations reached.
func (s *Solver) Solve(grid *Grid) {

	// Track Solve Time
	startTime := time.Now()

	// Loop until solved or MaxIterations reached
	iteration := 0
	for {

		// Track completion based on whether the Grid was updated)
		updated := false

		// Track Iterations
		iteration = iteration + 1
		if s.verbose {
			log.Printf("\n----- Iteration %d -----", iteration)
		}

		// Set any Cells where only 1 value is still possible
		updated = s.setSinglePossibleValueInGrid(grid)

		// Set any Cells where a possible value MUST belong for that Row
		// because it has been eliminated from all other Cells in the Row
		updated = updated || s.setOnlyPossibleValueInRow(grid)

		// Set any Cells where a possible value MUST belong for that Column
		// because it has been eliminated from all other Cells in the Column
		updated = updated || s.setOnlyPossibleValueInCol(grid)

		// Set any Cells where a possible value MUST belong for that Group
		// because it has been eliminated from all other Cells in the Group
		updated = updated || s.setOnlyPossibleValueInGroup(grid)

		// If no further updates were made then it should be solved!
		if !updated {
			break
		}

		// If we hit the max iterations then stop
		if iteration >= MaxIterations {
			log.Printf("Reached maximum iterations (%d) without solving", MaxIterations)
			break
		}
	}

	// Track solve time and log completion stats
	solveTime := time.Since(startTime)
	log.Printf("Finished solving in %d iterations over %s", iteration, solveTime.String())
}

// setSinglePossibleValueInGrid updates the Grid by Setting the value of
// any Cell which has only a single possible value remaining.
func (s *Solver) setSinglePossibleValueInGrid(grid *Grid) bool {

	// Track whether any updates wer made to the Grid
	updated := false

	// Loop over all the Cells in the Grid
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {

			// If the Cell only has a single remaining possible value, then set it!
			possibleValues := grid.GetCell(row, col).GetPossibleValues()
			if len(possibleValues) == 1 {
				value := possibleValues[0]
				s.logSetValueReason(row, col, value, "Only one possible value remaining for cell")
				grid.SetValue(row, col, value)
				updated = true
			}
		}
	}

	// Return Grid updated status
	return updated
}

// setOnlyPossibleValueInRow updates the Grid by Setting the value of any Cell
// which is the only remaining Cell in its Row with a particular possible Value.
func (s *Solver) setOnlyPossibleValueInRow(grid *Grid) bool {

	// Track whether any updates wer made to the Grid
	updated := false

	// Loop over the 9 Rows
	for row := 0; row < 9; row++ {

		// Loop over the possible Cell values (1-9)
		for value := 1; value <= 9; value++ {
			valueCount := 0
			valueCol := -1

			// Loop over the 9 Columns for the current Row
			for col := 0; col < 9; col++ {

				// Check the Cell's possible values for the current value
				if grid.GetCell(row, col).IsPossibleValue(value) {
					valueCount = valueCount + 1
					valueCol = col
				}

				// If the value already exists then cease further checks
				if valueCount > 1 {
					valueCol = -1
					break
				}
			}

			// If only a single Cell in the Row has the possible value, then set it!
			if valueCount == 1 {
				s.logSetValueReason(row, valueCol, value, "Only cell in the row with this possible value")
				grid.SetValue(row, valueCol, value)
				updated = true
			}
		}
	}

	// Return Grid updated status
	return updated
}

// setOnlyPossibleValueInCol updates the Grid by Setting the value of any Cell
// which is the only remaining Cell in its Column with a particular possible Value.
func (s *Solver) setOnlyPossibleValueInCol(grid *Grid) bool {

	// Track whether any updates wer made to the Grid
	updated := false

	// Loop over the 9 Columns
	for col := 0; col < 9; col++ {

		// Loop over the possible Cell values (1-9)
		for value := 1; value <= 9; value++ {
			valueCount := 0
			valueRow := -1

			// Loop over the 9 Rows for the current Column
			for row := 0; row < 9; row++ {

				// Check the Cell's possible values for the current value
				if grid.GetCell(row, col).IsPossibleValue(value) {
					valueCount = valueCount + 1
					valueRow = row
				}

				// If the value already exists then cease further checks
				if valueCount > 1 {
					valueRow = -1
					break
				}
			}

			// If only a single Cell in the Column has the possible value, then set it!
			if valueCount == 1 {
				s.logSetValueReason(valueRow, col, value, "Only cell in the column with possible value")
				grid.SetValue(valueRow, col, value)
				updated = true
			}
		}
	}

	// Return Grid updated status
	return updated
}

// setOnlyPossibleValueInGroup updates the Grid by Setting the value of any Cell
// which is the only remaining Cell in its Group with a particular possible Value.
func (s *Solver) setOnlyPossibleValueInGroup(grid *Grid) bool {

	// Track whether any updates wer made to the Grid
	updated := false

	// Loop over the 9 (sub) Groups selecting the upper-left Cell of each Group
	for groupRow := 0; groupRow <= 6; groupRow = groupRow + 3 {
		for groupCol := 0; groupCol <= 6; groupCol = groupCol + 3 {

			// Loop over the possible Cell values (1-9)
			for value := 1; value <= 9; value++ {
				valueCount := 0
				valueRow := -1
				valueCol := -1

			outerLoop:

				// Loop over the 9 Cells in the Group
				for groupCellRow := groupRow; groupCellRow < groupRow+3; groupCellRow++ {
					for groupCellCol := groupCol; groupCellCol < groupCol+3; groupCellCol++ {

						// Check the Cell's possible values for the current value
						if grid.GetCell(groupCellRow, groupCellCol).IsPossibleValue(value) {
							valueCount = valueCount + 1
							valueRow = groupCellRow
							valueCol = groupCellCol
						}

						// If the value already exists then cease further checks
						if valueCount > 1 {
							valueRow = -1
							valueCol = -1
							break outerLoop
						}
					}
				}

				// If only a single Cell in the Group has the possible value, then set it!
				if valueCount == 1 {
					s.logSetValueReason(valueRow, valueCol, value, "Only cell in the group with possible value")
					grid.SetValue(valueRow, valueCol, value)
					updated = true
				}
			}
		}
	}

	// Return Grid updated status
	return updated
}

// logSetValueReason logs a specific SetValue() decision with associated reason.
func (s *Solver) logSetValueReason(row int, col int, value int, reason string) {
	if s.verbose {
		log.Printf("SetValue  [%d,%d] --> %d     %s", row, col, value, reason)
	}
}
