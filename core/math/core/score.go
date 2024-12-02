package CoreMathCore

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	"sort"
)

// GetScoreHMLM 根据X和Y的值以及它们的中位数来计算得分
/**
用途：
可以将XY轴两组数据，输出为一组数据，形成归一化处理

得分分布：
1-4代表象限位置
|  2中  |  1高  |
|  3低  |  4中  |

使用方法：
1. 根据得分结果，可以将数据转换为坐标，进行可视化展示，必定为线性值
2. medX和medY可以用于对原始数据进行区域识别
3. medP为最终得分的中位数，可以用于对得分进行区域识别，例如将区域切分为2等分；
	如果需切分3等分，可medP=0~0.25; medP=0.25~0.75; medP=0.75~1
					medP*0.25 ; medP*0.25~medP*0.75; medP=0.75
	第二种切分方式为另外一种3等分，medP=0~0.33; medP=0.33~0.66; medP=0.66~1
4. 如果采用medXY切分，那么更容易业务人员理解；如果采用medP切分，比较倾向于算法领域。可根据实际需求切分
*/
func GetScoreHMLM(X, Y []float64) (result []float64, medX, medY, medP float64) {
	if len(X) != len(Y) {
		return
	}
	medX = GetMid(X)
	medY = GetMid(Y)
	scores := make([]float64, len(X))
	// 确定得分范围
	maxX, minX := maxAndMin(X)
	maxY, minY := maxAndMin(Y)
	if maxX == 0 || maxY == 0 {
		result = scores
		return
	}
	// 根据X和Y与中位数的相对位置来计算得分
	for i := 0; i < len(X); i++ {
		//Y值越低得分越小，X值越高得分越小
		// 根据这两个值与中位数的距离来分配得分
		xScore := 100 * ((X[i] - minX) / (maxX - minX))
		yScore := 100 * ((Y[i] - minY) / (maxY - minY))
		// 结合X和Y的得分，使用平均值，根据需要调整权重或结合方式
		scores[i] = (yScore + xScore) / 2
		// 根据象限调整得分
		//if X[i] >= medX {
		//	if Y[i] >= medY {
		//		//1
		//		scores[i] = scores[i] * 0.75
		//	} else {
		//		//4
		//		scores[i] = scores[i] * 0.05
		//	}
		//} else {
		//	if Y[i] < medY {
		//		//2
		//		scores[i] = scores[i]
		//	} else {
		//		//3
		//		scores[i] = scores[i]
		//	}
		//}
	}
	for k, v := range scores {
		scores[k] = CoreFilter.RoundToTwoDecimalPlaces(v)
	}
	result = scores
	if len(result) > 0 {
		medP = GetMid(result)
	}
	return
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
	// 初始化上级指标数组
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

// ClassifyEqualWidth 等宽区间划分算法
// data: 数据
// numBins: 切分数量
/**
使用方法：
1. 使用本方法输出数据，然后对应数据的值就是切分单元
2. 切分单元的key采用[0]~[n]为基准
*/
func ClassifyEqualWidth(data []float64, numBins int) ([][]float64, error) {
	//检查参数
	if len(data) == 0 {
		return nil, fmt.Errorf("data is empty")
	}
	if numBins <= 0 {
		return nil, fmt.Errorf("number of bins must be greater than 0")
	}
	// 排序数据
	sort.Float64s(data)
	// 计算区间宽度
	minVal := data[0]
	maxVal := data[len(data)-1]
	binWidth := (maxVal - minVal) / float64(numBins)
	// 初始化区间
	bins := make([][]float64, numBins)
	for i := range bins {
		bins[i] = []float64{}
	}
	// 划分数据到区间
	for _, val := range data {
		binIndex := int((val - minVal) / binWidth)
		if binIndex >= numBins {
			binIndex = numBins - 1
		}
		bins[binIndex] = append(bins[binIndex], val)
	}
	return bins, nil
}
