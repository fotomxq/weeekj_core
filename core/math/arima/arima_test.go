package CoreMathArima

import (
	"testing"
)

func TestSimpleArima(t *testing.T) {
	data := []float64{5, 10, 15, 25, 30, 50, 70}
	n := 3
	forecast, err := SimpleArima(data, n)
	if err != nil {
		t.Error("Error:", err)
	} else {
		t.Log("Forecast:", forecast)
	}
}
