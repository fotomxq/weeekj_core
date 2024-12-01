package CoreMathCore

// ZScoreXY Z-score标准化XY轴两组浮点数
func ZScoreXY(x float64, y float64) (x1 float64, y1 float64) {
	//计算
	x1 = (x - y) / y
	y1 = (y - x) / x
	//反馈
	return
}

// LinearCombination 计算线性组合关系
func LinearCombination(X, Y []float64, weightX, weightY float64) []float64 {
	// 检查X和Y的长度是否相等
	if len(X) != len(Y) {
		panic("x and y must have the same length")
	}
	// 初始化结果数组
	result := make([]float64, len(X))
	// 计算线性组合
	for i := 0; i < len(X); i++ {
		result[i] = weightX*X[i] + weightY*Y[i]
	}
	return result
}

// FeatureCrossing 计算X和Y的乘积来实现特征交叉
func FeatureCrossing(X, Y []float64) []float64 {
	// 检查X和Y的长度是否相等
	if len(X) != len(Y) {
		panic("x and y must have the same length")
	}
	// 初始化结果数组
	result := make([]float64, len(X))
	// 计算特征交叉（乘积）
	for i := 0; i < len(X); i++ {
		result[i] = X[i] * Y[i]
	}
	return result
}
