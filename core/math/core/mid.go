package CoreMathCore

// GetQuantile3 获取一个浮点数的3分位数
// 反馈三组数据，分别对应最小值、中间值、最大值区间
func GetQuantile3(data float64) (min, mid, max float64) {
	//计算
	min = data * 0.25
	mid = data * 0.5
	max = data * 0.75
	//反馈
	return
}

// GetQuantile2 获取一个浮点数的2分位数
// 反馈两组数据，分别对应最小值、最大值区间
func GetQuantile2(data float64) (min, max float64) {
	//计算
	min = data * 0.25
	max = data * 0.75
	//反馈
	return
}

// GetMid 获取中位数
func GetMid(data float64) (mid float64) {
	//计算
	mid = data * 0.5
	//反馈
	return
}
