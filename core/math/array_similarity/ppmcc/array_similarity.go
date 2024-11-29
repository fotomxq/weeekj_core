package CoreMathArraySimilarityPPMCC

import (
	"errors"
	"math"
)

//用于求解两组相同长度的数据，相似度有多少
// 主要采用：皮尔森相关系数
/**
x := []float64{1, 2, 3, 4, 5, 6}
	y := []float64{1, 20, 3, 40, 5, 60}
	res, err := ArraySimilarity(x, y)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
*/

// ppmmccUtil 求解结构
type ppmmccUtil interface {
	// 求均值
	getAverageValue()
	// 求协方差
	getCovariance() float64
	// 求标准差
	getStandardDeviation() float64
	// 求相关系数
	getPPMCC() (float64, error)
}

type variableArray struct {
	// 样本数据 X
	X []float64
	// 样本数据 Y
	Y []float64
	// 样本大小
	Samplesize int
}

// 求均值
func (v *variableArray) getAverageValue() {
	var xSum float64 = 0
	var ySum float64 = 0
	for i := 0; i < v.Samplesize; i++ {
		xSum = xSum + v.X[i]
		ySum = ySum + v.Y[i]
	}
	// 新增一位用于保存均值
	v.X = append(v.X, xSum/float64(v.Samplesize))
	v.Y = append(v.Y, ySum/float64(v.Samplesize))
}

// 求协方差
func (v *variableArray) getCovariance() float64 {
	var res float64 = 0
	for i := 0; i < v.Samplesize; i++ {
		res = res + (v.X[i]-v.X[v.Samplesize])*(v.Y[i]-v.Y[v.Samplesize])
	}
	return res
}

// 求标准差
func (v *variableArray) getStandardDeviation() float64 {
	var xRes float64 = 0
	var yRes float64 = 0
	for i := 0; i < v.Samplesize; i++ {
		xRes = xRes + (v.X[i]-v.X[v.Samplesize])*(v.X[i]-v.X[v.Samplesize])
		yRes = yRes + (v.Y[i]-v.Y[v.Samplesize])*(v.Y[i]-v.Y[v.Samplesize])
	}
	// math.Sqrt() 开根号函数
	var res float64 = math.Sqrt(xRes) * math.Sqrt(yRes)
	return res
}

// 求相关系数
func (v *variableArray) getPPMCC() (float64, error) {
	if len(v.X) != len(v.Y) {
		return 0, errors.New("the length of the two arrays is different")
	}
	v.getAverageValue()
	return v.getCovariance() / v.getStandardDeviation(), nil
}

// ArraySimilarity 求数据相似度
func ArraySimilarity(x []float64, y []float64) float64 {
	var ppmccUtil ppmmccUtil = &variableArray{
		X:          x,
		Y:          y,
		Samplesize: len(x),
	}
	ppmcc, err := ppmccUtil.getPPMCC()
	if err != nil {
		return 0
	}
	return ppmcc
}
