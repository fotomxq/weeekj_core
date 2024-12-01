package CoreMathCore

import "testing"

func TestLinearCombination(t *testing.T) {
	X := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	Y := []float64{10.0, 20.0, 30.0, 40.0, 50.0}
	// 计算线性组合
	linearCombined := LinearCombination(X, Y, 0.5, 0.5)
	t.Log(linearCombined)
}

func TestFeatureCrossing(t *testing.T) {
	X := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	Y := []float64{10.0, 20.0, 30.0, 40.0, 50.0}
	// 计算特征交叉
	crossedFeatures := FeatureCrossing(X, Y)
	t.Log(crossedFeatures)
}
