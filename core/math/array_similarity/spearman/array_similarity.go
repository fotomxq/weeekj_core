package CoreMathArraySimilaritySpearman

import (
	"math"
	"sort"
)

func ArraySimilarity(x, y []float64) float64 {
	if len(x) != len(y) {
		panic("两组参数的长度必须相同")
	}
	n := len(x)
	if n <= 1 {
		return 1.0
	}
	// 获取 x 和 y 的排名
	ranksX := rank(x)
	ranksY := rank(y)
	// 计算排名差的平方和
	var dSquaredSum float64
	for i := 0; i < n; i++ {
		dSquaredSum += math.Pow(ranksX[i]-ranksY[i], 2)
	}
	// 计算斯皮尔曼相关系数
	spearmanCoeff := 1 - (6*dSquaredSum)/(float64(n)*(float64(n)*float64(n)-1))
	return spearmanCoeff
}

// rank 返回参数的排名
func rank(arr []float64) []float64 {
	// 创建一个索引切片
	indices := make([]int, len(arr))
	for i := range indices {
		indices[i] = i
	}
	// 根据数组值排序索引
	sort.Slice(indices, func(i, j int) bool {
		return arr[indices[i]] < arr[indices[j]]
	})
	// 生成排名
	ranks := make([]float64, len(arr))
	rank := 1.0
	for i, idx := range indices {
		if i > 0 && arr[idx] == arr[indices[i-1]] {
			// 处理重复值，赋予相同的排名
			ranks[idx] = ranks[indices[i-1]]
		} else {
			ranks[idx] = rank
			rank++
		}
	}
	return ranks
}
