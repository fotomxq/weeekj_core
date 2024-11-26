package CoreMathCore

// Normalize 归一化数据
func Normalize(data []float64) (result []float64) {
	if len(data) < 1 {
		return
	}
	//获取最大值和最小值
	var maxD float64 = -1
	var minD float64 = -1
	for _, v := range data {
		if maxD == -1 || v > maxD {
			maxD = v
		}
		if minD == -1 || v < minD {
			minD = v
		}
	}
	//归一化
	for _, v := range data {
		result = append(result, (v-minD)/(maxD-minD))
	}
	return
}
