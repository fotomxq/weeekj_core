package CoreMathCore

import (
	"fmt"
	"testing"
)

func TestGetScoreHMLM(t *testing.T) {
	/**
	|  中  |  高  |
	|  低  |  中  |
	*/
	X := []float64{7, 8, 9, 1, 2, 3, 1, 2, 3, 1, 2, 3}
	Y := []float64{7, 8, 9, 7, 8, 9, 7, 8, 9, 1, 2, 3}
	// 计算得分
	scores, medX, medY, medP := GetScoreHMLM(X, Y)
	t.Log("x: ", X)
	t.Log("y: ", Y)
	t.Log("medX: ", medX)
	t.Log("medY: ", medY)
	t.Log("medP: ", medP)
	t.Log("scores: ", scores)
}

func TestGetScoreWeightedSum(t *testing.T) {
	indicators := [][]float64{
		{1.0, 2.0, 3.0}, // 指标1
		{4.0, 5.0, 6.0}, // 指标2
		{7.0, 8.0, 9.0}, // 指标3
	}
	weights := []float64{0.2, 0.3, 0.5} // 对应的指标权重
	// 计算上级风险指标
	compositeIndicator, err := GetScoreWeightedSum(indicators, weights)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(compositeIndicator)
}

func TestClassifyEqualWidth(t *testing.T) {
	data := []float64{0.1, 0.4, 0.35, 0.8, 0.45, 0.9, 0.2, 0.75, 0.6, 0.15}
	numBins := 3
	bins, err := ClassifyEqualWidth(data, numBins)
	if err != nil {
		t.Error(err)
		return
	}
	for i, bin := range bins {
		t.Log(fmt.Sprintf("Bin %d: %v", i, bin))
	}
}
