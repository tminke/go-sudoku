package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSolver(t *testing.T) {
	solver := NewSolver(99, true)
	assert.NotNil(t, solver)
	assert.Equal(t, 99, solver.maxIterations)
	assert.Equal(t, true, solver.verbose)
}

func TestSolve(t *testing.T) {

	// Manual hook for debugging
	verbose := false

	// Create a test Grid
	grid := testGrid()

	// Verbose Log Grid - Before
	if verbose {
		t.Logf("Before:\n\n%s\n", grid.String())
	}

	// Create a Solver and solve the puzzle!
	solver := NewSolver(100, verbose)
	solver.Solve(grid)

	// Verbose Log Grid - After
	if verbose {
		t.Logf("After:\n\n%s\n", grid.String())
	}

	// Verify all rows contain 1-9
	for row := 0; row < 9; row++ {
		values := [9]bool{}
		for col := 0; col < 9; col++ {
			value := grid.GetCell(row, col).GetValue()
			assert.False(t, values[value-1])
			values[value-1] = true
		}
		for index := 0; index < 9; index++ {
			assert.True(t, values[index])
		}
	}

	// Verify all columns contain 1-9
	for col := 0; col < 9; col++ {
		values := [9]bool{}
		for row := 0; row < 9; row++ {
			value := grid.GetCell(row, col).GetValue()
			assert.False(t, values[value-1])
			values[value-1] = true
		}
		for index := 0; index < 9; index++ {
			assert.True(t, values[index])
		}
	}

	// Verify all sub-groups contain 1-9
	for groupRow := 0; groupRow <= 6; groupRow = groupRow + 3 {
		for groupCol := 0; groupCol <= 6; groupCol = groupCol + 3 {
			values := [9]bool{}
			// Loop over the 9 Cells in the Group
			for groupCellRow := groupRow; groupCellRow < groupRow+3; groupCellRow++ {
				for groupCellCol := groupCol; groupCellCol < groupCol+3; groupCellCol++ {
					value := grid.GetCell(groupCellRow, groupCellCol).GetValue()
					assert.False(t, values[value-1])
					values[value-1] = true
				}
			}
			for index := 0; index < 9; index++ {
				assert.True(t, values[index])
			}
		}
	}
}
