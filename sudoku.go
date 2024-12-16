package main

import (
	"flag"
	"log"

	sudoku "github.com/tminke/go-sudoku/internal"
)

func main() {

	// Parse Flags
	csvFile := flag.String("file", "sudoku.csv", "Path/Name of the CSV file containing the sudoku puzzle (default = sudoku.csv).")
	maxIterations := flag.Int("iter", 50, "The maximum number of iterations before abandoning the solve (default = 50).")
	verbose := flag.Bool("verbose", false, "Whether or not to log the individual steps in the solve (default = false).")
	flag.Parse()

	// Create A Grid From The Specified Sudoku CSV File
	grid, err := sudoku.NewGridFromCsv(*csvFile)
	if err != nil {
		log.Fatalf("Failed to load CSV file: err=%+v", err)
	}
	log.Printf("Problem:\n\n%s\n", grid)

	// Create A New Sudoku Solver
	solver := sudoku.NewSolver(*maxIterations, *verbose)

	// Solve The Sudoku Puzzle & Log The Result
	solver.Solve(grid)
	log.Printf("Solution:\n\n%s\n", grid)
}
