package CoreMathArima

import (
	"errors"
	"github.com/sajari/regression"
)

func SimpleArima(data []float64, n int) ([]float64, error) {
	if len(data) == 0 {
		return nil, errors.New("data cannot be empty")
	}
	if n <= 0 {
		return nil, errors.New("n must be greater than 0")
	}
	var r regression.Regression
	r.SetObserved("Sales")
	r.SetVar(0, "Time")
	// Train the model
	for i, d := range data {
		r.Train(regression.DataPoint(d, []float64{float64(i)}))
	}
	// Fit the model
	err := r.Run()
	if err != nil {
		return nil, err
	}
	// Predict the future values
	forecast := make([]float64, n)
	for i := 0; i < n; i++ {
		forecast[i], err = r.Predict([]float64{float64(len(data) + i)})
		if err != nil {
			return nil, err
		}
	}
	return forecast, nil
}
