package qc

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	curve, _ := NewQuadCurveFromTriplePoints(1.0, 5.0, 2.0, 3.0, 3.0, 3.0)
	assert.Equal(t, 1.0, curve.A)
	assert.Equal(t,-5.0, curve.B)
	assert.Equal(t, 9.0, curve.C)
	assert.Equal(t, int64(36), curve.CalcArea(6))
}
