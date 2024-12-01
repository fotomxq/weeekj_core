package CoreMathCore

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
)

// GetScoreLLMH 根据X和Y的值以及它们的中位数来计算得分
/**
用途：
可以将XY轴两组数据，输出为一组数据，形成归一化处理

得分分布：
第一象限、第二象限得分低；第三象限得分中；第四象限得分高
|  低  |  低  |
|  中  |  高  |
*/
func GetScoreLLMH(X, Y []float64) []float64 {
	if len(X) != len(Y) {
		return []float64{}
	}
	medX := GetMid(X)
	medY := GetMid(Y)
	scores := make([]float64, len(X))
	// 确定得分范围
	maxX, minX := maxAndMin(X)
	maxY, minY := maxAndMin(Y)
	if maxX == 0 || maxY == 0 {
		return scores
	}
	// 根据X和Y与中位数的相对位置来计算得分
	for i := 0; i < len(X); i++ {
		//Y值越低得分越小，X值越高得分越小
		// 根据这两个值与中位数的距离来分配得分
		yScore := 100 * (Y[i] - minY) / (maxY - minY)
		xScore := 100 * (maxX - X[i]) / (maxX - minX)
		// 结合X和Y的得分，这里简单使用平均值
		// 可以根据需要调整权重或结合方式
		scores[i] = (yScore + xScore) / 2
		// 根据象限调整得分
		if Y[i] < medY {
			if X[i] < medX {
				// 第三象限：中，降低得分
				scores[i] *= 0.5
			} else {
				// 第四象限：高，提高得分
				scores[i] = 100 - scores[i]*0.5
			}
		} else {
			// 第一或第二象限：低，保持或稍微降低得分
			scores[i] = scores[i] * 0.5
		}
	}
	for k, v := range scores {
		scores[k] = CoreFilter.RoundToTwoDecimalPlaces(v)
	}
	return scores
}

// GetScoreWeightedSum 计算加权输出上级得分
func GetScoreWeightedSum(indicators [][]float64, weights []float64) ([]float64, error) {
	// 检查输入的有效性
	if len(indicators) == 0 {
		return nil, errors.New("no indicators")
	}
	numSamples := len(indicators[0])
	if numSamples == 0 {
		return nil, errors.New("no samples")
	}
	if len(weights) != len(indicators) {
		return nil, errors.New("invalid weights")
	}
	// 初始化上级风险指标数组
	compositeIndicator := make([]float64, numSamples)
	// 计算加权和
	for i := 0; i < numSamples; i++ {
		var sum float64
		for j, indicator := range indicators {
			sum += indicator[i] * weights[j]
		}
		compositeIndicator[i] = sum
	}
	return compositeIndicator, nil
}
