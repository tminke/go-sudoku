package internal

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSudokuCsv(t *testing.T) {

	// Test Data
	initialData := [9][9]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	// Define the test cases
	testCases := map[string]struct {
		fileContent []string
		expectData  [9][9]int
		expectErr   string
	}{
		"Not Enough Rows Error": {
			fileContent: []string{ // 8 Rows
				"1,2,3,4,5,6,7,8,9\n",
				"1,2,3,4,5,6,7,8,9\n",
				"1,2,3,4,5,6,7,8,9\n",
				"1,2,3,4,5,6,7,8,9\n",
				"1,2,3,4,5,6,7,8,9\n",
				"1,2,3,4,5,6,7,8,9\n",
				"1,2,3,4,5,6,7,8,9\n",
				"1,2,3,4,5,6,7,8,9\n",
			},
			expectData: initialData,
			expectErr:  "exactly nine rows expected",
		},
		"Too Many Rows Error": {
			fileContent: []string{ // 10 Rows
				"1,2,3,4,5,6,7,8,9\n",
				"1,2,3,4,5,6,7,8,9\n",
				"1,2,3,4,5,6,7,8,9\n",
				"1,2,3,4,5,6,7,8,9\n",
				"1,2,3,4,5,6,7,8,9\n",
				"1,2,3,4,5,6,7,8,9\n",
				"1,2,3,4,5,6,7,8,9\n",
				"1,2,3,4,5,6,7,8,9\n",
				"1,2,3,4,5,6,7,8,9\n",
				"1,2,3,4,5,6,7,8,9\n",
			},
			expectData: initialData,
			expectErr:  "exactly nine rows expected",
		},
		"Not Enough Cols Error": {
			fileContent: []string{ // 8 Cols
				"1,2,3,4,5,6,7,8\n",
				"1,2,3,4,5,6,7,8\n",
				"1,2,3,4,5,6,7,8\n",
				"1,2,3,4,5,6,7,8\n",
				"1,2,3,4,5,6,7,8\n",
				"1,2,3,4,5,6,7,8\n",
				"1,2,3,4,5,6,7,8\n",
				"1,2,3,4,5,6,7,8\n",
				"1,2,3,4,5,6,7,8\n",
			},
			expectData: initialData,
			expectErr:  "exactly nine cols expected",
		},
		"Too Many Cols Error": {
			fileContent: []string{ // 10 Cols
				"1,2,3,4,5,6,7,8,9,10\n",
				"1,2,3,4,5,6,7,8,9,10\n",
				"1,2,3,4,5,6,7,8,9,10\n",
				"1,2,3,4,5,6,7,8,9,10\n",
				"1,2,3,4,5,6,7,8,9,10\n",
				"1,2,3,4,5,6,7,8,9,10\n",
				"1,2,3,4,5,6,7,8,9,10\n",
				"1,2,3,4,5,6,7,8,9,10\n",
				"1,2,3,4,5,6,7,8,9,10\n",
			},
			expectData: initialData,
			expectErr:  "exactly nine cols expected",
		},
		"Non Int": {
			fileContent: []string{
				"a,2,3,4,5,6,7,8,9\n", // a is not an integer
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
			},
			expectData: initialData,
			expectErr:  "encountered unsupported value 'a'",
		},
		"Invalid Large Int": {
			fileContent: []string{
				"10,2,3,4,5,6,7,8,9\n", // 10 is out of range 1-9
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
			},
			expectData: initialData,
			expectErr:  "encountered invalid value '10'",
		},
		"Invalid Small Int": {
			fileContent: []string{
				"0,2,3,4,5,6,7,8,9\n", // 0 is out of range 1-9
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
				"-,-,-,-,-,-,-,-,-\n",
			},
			expectData: initialData,
			expectErr:  "encountered invalid value '0'",
		},
		"Valid Data": {
			fileContent: []string{ // Valid Content (for parsing if not sudoku :)
				"1,2,3,4,5,6,7,8,9\n",
				"-,-,-,-,-,-,-,-,-\n",
				"1,2,3,4,5,6,7,8,9\n",
				"-,-,-,-,-,-,-,-,-\n",
				"1,2,3,4,5,6,7,8,9\n",
				"-,-,-,-,-,-,-,-,-\n",
				"1,2,3,4,5,6,7,8,9\n",
				"-,-,-,-,-,-,-,-,-\n",
				"1,2,3,4,5,6,7,8,9\n",
			},
			expectData: [9][9]int{
				{1, 2, 3, 4, 5, 6, 7, 8, 9},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{1, 2, 3, 4, 5, 6, 7, 8, 9},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{1, 2, 3, 4, 5, 6, 7, 8, 9},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{1, 2, 3, 4, 5, 6, 7, 8, 9},
				{0, 0, 0, 0, 0, 0, 0, 0, 0},
				{1, 2, 3, 4, 5, 6, 7, 8, 9},
			},
			expectErr: "",
		},
	}

	// Execute the test cases
	for testCaseName, testCase := range testCases {
		t.Run(testCaseName, func(t *testing.T) {

			// Create a temporary file in current directory
			file, err := os.CreateTemp("", "parse-sudoku-csv-*.csv")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(file.Name())

			// Write the content to the file
			for _, content := range testCase.fileContent {
				_, err = file.WriteString(content)
				if err != nil {
					t.Fatal(err)
				}
			}

			// Perform the test
			data, err := parseSudokuCsv(file.Name())

			// Verify the results
			assert.Equal(t, testCase.expectData, data)
			if testCase.expectErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Contains(t, err.Error(), testCase.expectErr)
			}
		})
	}
}
