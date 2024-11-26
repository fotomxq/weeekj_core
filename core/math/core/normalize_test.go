package CoreMathCore

import "testing"

func TestNormalize(t *testing.T) {
	data := []float64{2031, 322, 3323, 44, 235}
	result := Normalize(data)
	t.Log(result)
}
