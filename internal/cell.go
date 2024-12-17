package internal

import (
	"strconv"
	"sync"
)

// Cell maintains the current state of an individual Cell on the Sudoko board.
type Cell struct {
	value    int          // The known value of the Cell (0 indicates unknown)
	possible [9]bool      // Values still possible for this Cell
	mutex    sync.RWMutex // Protect for potential parallel access
}

// NewCell returns an initialized Cell where all values are still possible.
func NewCell() *Cell {
	return &Cell{
		value:    0,                                                             // Zero indicates no value currently set
		possible: [9]bool{true, true, true, true, true, true, true, true, true}, // All values are possible to start
	}
}

// GetValue returns the current known value of the Cell or 0 if unknown.
func (c *Cell) GetValue() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.value
}

// GetValueString returns a string representation of the current value or
// and empty string " " if unknown.
func (c *Cell) GetValueString() string {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.value == 0 {
		return " "
	}
	return strconv.Itoa(c.value)
}

// GetPossibleValues returns an array of possible values in order.  The indexes in the
// array are not relevant and the array contains the "values" that are still possible.
func (c *Cell) GetPossibleValues() []int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	possibleValues := []int{}
	for i := 0; i < 9; i++ {
		if c.possible[i] {
			possibleValues = append(possibleValues, i+1)
		}
	}
	return possibleValues
}

// IsPossibleValue is a convenience function that returns a boolean indication
// of whether the specified value is still possible.
func (c *Cell) IsPossibleValue(value int) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.possible[value-1]
}

// SetValue sets the Cell to the specified value and removes all possible values.
func (c *Cell) SetValue(value int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.value = value
	for i := 0; i < 9; i++ {
		c.possible[i] = false
	}
}

// EliminateValue will mark the specified value as no longer possible for the Cell.
func (c *Cell) EliminateValue(value int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.possible[value-1] = false
}
