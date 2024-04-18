package CoreMathArima

import (
	"testing"
)

func TestSimpleArima(t *testing.T) {
	data := []float64{5, 10, 15, 25, 30, 50, 70, 60, 30, 20, 15, 30, 40, 60, 70, 80, 101, 90, 70, 60, 50, 30, 10, 5}
	n := 6
	forecast, err := SimpleArima(data, n)
	if err != nil {
		t.Error("Error:", err)
	} else {
		t.Log("Forecast:", forecast)
	}
}
