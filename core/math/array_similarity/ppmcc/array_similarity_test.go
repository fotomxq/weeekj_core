package CoreMathArraySimilarityPPMCC

import "testing"

func TestArraySimilarity(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5, 6}
	y := []float64{1, 20, 3, 40, 5, 60}
	res := ArraySimilarity(x, y)
	// 1:  0.6412627329836124 = 64.12%
	t.Log("1: ", res)
	x2 := []float64{1, 20, 5, 41, 7, 62}
	y2 := []float64{1, 20, 3, 40, 5, 60}
	res = ArraySimilarity(x2, y2)
	// 2:  0.9992102527561709 = 99.92%
	t.Log("2: ", res)
}
