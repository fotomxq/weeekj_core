package CoreMathCore

import "math"

// DiscreteMean 期望
// 预期平均值
func DiscreteMean(v []float64) float64 {
	var res float64 = 0
	var n int = len(v)
	for i := 0; i < n; i++ {
		res += v[i]
	}
	return res / float64(n)
}

// DiscreteVariance 方差
// 直接识别一组数据的波动性
func DiscreteVariance(v []float64) float64 {
	var res float64 = 0
	var m = DiscreteMean(v)
	var n int = len(v)
	for i := 0; i < n; i++ {
		res += (v[i] - m) * (v[i] - m)
	}
	return res / float64(n-1)
}

// DiscreteStd 标准差
// 基于方差基础上，更符合认知
// 便于识别波动剧烈性
func DiscreteStd(v []float64) float64 {
	return math.Sqrt(DiscreteVariance(v))
}
