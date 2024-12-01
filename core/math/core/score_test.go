package CoreMathCore

import (
	"testing"
)

func TestGetScoreLLMH(t *testing.T) {
	/**
	|  低  |  低  |
	|  中  |  高  |
	*/
	X := []float64{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 0}
	Y := []float64{10.0, 20.0, 5.0, 4.0, 15.0, 25.0, 15}
	// 计算得分
	scores := GetScoreLLMH(X, Y)
	t.Log(scores)
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
