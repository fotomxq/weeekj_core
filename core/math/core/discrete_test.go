package CoreMathCore

import "testing"

func TestDiscreteMean(t *testing.T) {
	v := []float64{51, 55, 60, 65, 70}
	t.Log(DiscreteMean(v))
}

func TestDiscreteVariance(t *testing.T) {
	v := []float64{51, 55, 60, 65, 70}
	t.Log(DiscreteVariance(v))
}

func TestDiscreteStd(t *testing.T) {
	v := []float64{51, 55, 60, 65, 70}
	t.Log(DiscreteStd(v))
}
