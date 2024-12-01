package CoreMathCore

import "sort"

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
func GetMid(data []float64) (result float64) {
	n := len(data)
	if n == 0 {
		return
	}
	sorted := make([]float64, len(data))
	copy(sorted, data)
	sort.Float64s(sorted)
	mid := n / 2
	if n%2 == 0 {
		return (sorted[mid-1] + sorted[mid]) / 2
	}
	return sorted[mid]
}

// maxAndMin 获取数组中的最大值和最小值
func maxAndMin(arr []float64) (float64, float64) {
	arrMax, arrMin := arr[0], arr[0]
	for _, v := range arr {
		if v > arrMax {
			arrMax = v
		}
		if v < arrMin {
			arrMin = v
		}
	}
	return arrMax, arrMin
}
