package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	allPossibleValues  = [9]bool{true, true, true, true, true, true, true, true, true}
	noPossibleValues   = [9]bool{false, false, false, false, false, false, false, false, false}
	evenPossibleValues = [9]bool{false, true, false, true, false, true, false, true, false}
)

func TestNewCell(t *testing.T) {
	cell := NewCell()
	assert.NotNil(t, cell)
	assert.Equal(t, 0, cell.value)
	assert.Equal(t, allPossibleValues, cell.possible)
}

func TestCell_GetValue(t *testing.T) {
	cell := &Cell{value: 7}
	assert.Equal(t, 7, cell.GetValue())
}

func TestCell_GetValueString(t *testing.T) {
	cell := &Cell{value: 7}
	assert.Equal(t, "7", cell.GetValueString())
	cell.SetValue(0)
	assert.Equal(t, " ", cell.GetValueString())
}

func TestCell_GetPossibleValues(t *testing.T) {
	cell := &Cell{possible: evenPossibleValues}
	assert.Equal(t, []int{2, 4, 6, 8}, cell.GetPossibleValues())
}

func TestCell_IsPossibleValue(t *testing.T) {
	cell := &Cell{possible: evenPossibleValues}
	assert.True(t, cell.IsPossibleValue(2))
	assert.True(t, cell.IsPossibleValue(4))
	assert.True(t, cell.IsPossibleValue(6))
	assert.True(t, cell.IsPossibleValue(8))
	assert.False(t, cell.IsPossibleValue(1))
	assert.False(t, cell.IsPossibleValue(3))
	assert.False(t, cell.IsPossibleValue(5))
	assert.False(t, cell.IsPossibleValue(7))
	assert.False(t, cell.IsPossibleValue(9))
}

func TestCell_SetValue(t *testing.T) {
	value := 7
	cell := NewCell()
	assert.Equal(t, 0, cell.value)
	assert.Equal(t, allPossibleValues, cell.possible)
	cell.SetValue(value)
	assert.Equal(t, value, cell.value)
	assert.Equal(t, noPossibleValues, cell.possible)
}

func TestCell_EliminateValue(t *testing.T) {
	cell := &Cell{possible: evenPossibleValues}
	cell.EliminateValue(4)
	assert.Equal(t, [9]bool{false, true, false, false, false, true, false, true, false}, cell.possible)
	cell.EliminateValue(8)
	assert.Equal(t, [9]bool{false, true, false, false, false, true, false, false, false}, cell.possible)
}
