package CoreMathRFM

import CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"

/**
RFM模型
1. 用于测量客户价值和客户分群的一种方法
2. 本模块提供基础运算逻辑，根据需求调用即可
*/

type Core struct {
	//权重列表
	weightList []Weight
	//数据的最小值
	rMin float64
	fMin float64
	mMin float64
	//数据的最大值
	rMax float64
	fMax float64
	mMax float64
}

type Weight struct {
	//编号
	// 可用于系统固化ID、排序、具体业务逻辑等
	Number int64
	//权重值
	R float64
	F float64
	M float64
}

// SetWeight 设置权重
func (t *Core) SetWeight(weightList []Weight) {
	t.weightList = weightList
}

// GetWeight 获取权重
func (t *Core) GetWeight(num int64) Weight {
	if len(t.weightList) < 1 {
		return Weight{
			Number: num,
			R:      0.3,
			F:      0.3,
			M:      0.6,
		}
	}
	for k := 0; k < len(t.weightList); k++ {
		v := t.weightList[k]
		if v.Number == num {
			return v
		}
		if num > v.Number {
			return v
		}
	}
	return t.weightList[len(t.weightList)-1]
}

// SetDataRange 设置数据范围
func (t *Core) SetDataRange(rMin float64, fMin float64, mMin float64, rMax float64, fMax float64, mMax float64) {
	t.rMin = rMin
	t.fMin = fMin
	t.mMin = mMin
	t.rMax = rMax
	t.fMax = fMax
	t.mMax = mMax
}

// GetScoreByWeight 获取分数
func (t *Core) GetScoreByWeight(recency float64, frequency float64, monetary float64, widthNum int64) (score float64) {
	//获取权重
	var weight Weight
	weight = t.GetWeight(widthNum)
	//计算RFM得分
	score = t.GetScore(recency, frequency, monetary, weight.R, weight.F, weight.M, t.rMin, t.fMin, t.mMin, t.rMax, t.fMax, t.mMax)
	//反馈
	return
}

// GetScore 获取分数底层方法
// 注意，R指标方向和F、M指标相同
func (t *Core) GetScore(recency float64, frequency float64, monetary float64, weightR float64, weightF float64, weightM float64, minValR float64, minValF float64, minValM float64, maxValR float64, maxValF float64, maxValM float64) (score float64) {
	//检查意外值
	var r, f, m float64
	//归一化处理
	if (maxValR - minValR) != 0 {
		r = (recency - minValR) / (maxValR - minValR)
	}
	if (maxValF - minValF) != 0 {
		f = (frequency - minValF) / (maxValF - minValF)
	}
	if (maxValM - minValM) != 0 {
		m = (monetary - minValM) / (maxValM - minValM)
	}
	//计算RFM得分
	score = weightR*r + weightF*f + weightM*m
	//保留小数点2位
	score = CoreFilter.RoundToTwoDecimalPlaces(score)
	//反馈
	return
}

// GetScoreR2 反转R指标
// R指标与FM不同，尤其可用于
func (t *Core) GetScoreR2(recency float64, frequency float64, monetary float64, weightR float64, weightF float64, weightM float64, minValR float64, minValF float64, minValM float64, maxValR float64, maxValF float64, maxValM float64) (score float64) {
	//检查意外值
	var r, f, m float64
	//归一化处理
	if (maxValR - minValR) != 0 {
		r = (maxValR - recency) / (maxValR - minValR)
	}
	if (maxValF - minValF) != 0 {
		f = (frequency - minValF) / (maxValF - minValF)
	}
	if (maxValM - minValM) != 0 {
		m = (monetary - minValM) / (maxValM - minValM)
	}
	//计算RFM得分
	score = weightR*r + weightF*f + weightM*m
	//保留小数点2位
	score = CoreFilter.RoundToTwoDecimalPlaces(score)
	//反馈
	return
}
