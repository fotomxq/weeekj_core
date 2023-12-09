package CoreMathLinearRegression

import (
	"fmt"
	"testing"
)

func TestLinearRegression(t *testing.T) {
	data := []float64{360, 351, 355, 325, 370}
	res := LinearRegression(data, 2)
	fmt.Println(res)
}
