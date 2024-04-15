package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCalculate struct {
	mock.Mock
}

func (m *MockCalculate) Add(a, b int) int {
	return a - b
}

func TestCalculate(t *testing.T) {
	m := new(MockCalculate)
	m.On("Add", 1, 1).Return(0)

	result := Calculate(1, 1, m)

	assert.Equal(t, 0, result)

}
