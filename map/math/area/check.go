package MapMathArea

import (
	"errors"
	MapMathArgs "github.com/fotomxq/weeekj_core/v5/map/math/args"
	MapMathConversion "github.com/fotomxq/weeekj_core/v5/map/math/conversion"
	"strings"
)

//对一组数据进行检查处理
// 用于多组面发生重叠时，能够排列出多个数据

type ArgsCheckXYInAreaList struct {
	//检查的点
	CheckPoint MapMathArgs.ParamsPoint
	//分区数据
	AreaDataList []ParamsArea
}

func CheckXYInAreaList(args *ArgsCheckXYInAreaList) ([]ParamsArea, error) {
	//初始化
	var result []ParamsArea
	//检查点在哪个区间出现
	var resultIn []ParamsArea
	for _, v := range args.AreaDataList {
		if CheckXYInArea(&ArgsCheckXYInArea{
			Point: args.CheckPoint, Area: v,
		}) == false {
			continue
		}
		resultIn = append(resultIn, v)
	}
	//如果没有数据，则说明该点均不在范围内，直接反馈失败
	if len(resultIn) < 1 {
		return nil, errors.New("not in range")
	}
	//完成后，重新遍历
	// 计算每个区域最小值和最大值的中心点
	type diffType struct {
		//序列ID
		ID int64
		//最小值 XY
		Min ParamsAreaPoint
		//最大值 XY
		Max ParamsAreaPoint
		//中心点
		Center ParamsAreaPoint
		//中心点相加后，与XY相加的差值，用于快速计算最小的，即最接近的值
		Diff float64
	}
	var resultDifference []diffType
	for _, v := range resultIn {
		//初始化结构
		vXY := diffType{
			v.ID,
			ParamsAreaPoint{
				Longitude: v.Points[0].Longitude,
				Latitude:  v.Points[0].Latitude,
			},
			ParamsAreaPoint{
				Longitude: v.Points[0].Longitude,
				Latitude:  v.Points[0].Latitude,
			},
			ParamsAreaPoint{
				Longitude: v.Points[0].Longitude,
				Latitude:  v.Points[0].Latitude,
			},
			(v.Points[0].Longitude + v.Points[0].Latitude) - (args.CheckPoint.Longitude + args.CheckPoint.Latitude),
		}
		//对比数据并替换
		for _, v2 := range v.Points {
			//如果更小
			if v2.Longitude+v2.Latitude < vXY.Min.Longitude+vXY.Min.Latitude {
				vXY.Min = v2
			}
			//如果更大
			if v2.Longitude+v2.Latitude > vXY.Max.Longitude+vXY.Max.Latitude {
				vXY.Max = v2
			}
			//计算中心点
			vXY.Center = ParamsAreaPoint{
				Longitude: vXY.Max.Longitude - vXY.Min.Longitude,
				Latitude:  vXY.Max.Latitude - vXY.Min.Latitude,
			}
			//计算相差
			vXY.Diff = (vXY.Center.Longitude + vXY.Center.Latitude) - (args.CheckPoint.Longitude + args.CheckPoint.Latitude)
		}
		//插入数据
		resultDifference = append(resultDifference, vXY)
	}
	// 根据中心点，计算距离最接近XY的是哪个
	// 如果多个相同，则返回多个值
	var resultDifference2 []diffType
	for _, v := range resultDifference {
		isOK := true
		for _, v2 := range resultDifference2 {
			if v.Diff > v2.Diff {
				isOK = false
			}
		}
		if isOK == true {
			resultDifference2 = append(resultDifference2, v)
		}
	}
	//抽取数据到新的集合内
	for _, v := range resultDifference2 {
		for _, v2 := range args.AreaDataList {
			if v.ID == v2.ID {
				result = append(result, v2)
			}
		}
	}
	//如果不存在数据，则返回失败
	if len(result) < 1 {
		return nil, errors.New("not in range, and result is error by finally")
	}
	//返回结果
	return result, nil
}

// ArgsCheckXYInArea 检查点和面的数学问题，点是否在面内
// 采用矩阵关系处理
// 测试耗时：10,000次 / 0.038秒
type ArgsCheckXYInArea struct {
	//点坐标
	Point MapMathArgs.ParamsPoint
	//分区范围
	Area ParamsArea
}

func CheckXYInArea(args *ArgsCheckXYInArea) bool {
	//如果分区数据不足，则推出
	if len(args.Area.Points) < 1 {
		return false
	}
	//如果数据格式不一致，则转化
	if strings.ToTitle(args.Point.PointType) != strings.ToTitle(args.Area.PointType) {
		conversionRes, err := MapMathConversion.Conversion(&MapMathConversion.ArgsConversion{
			SrcType:  args.Point.PointType,
			DestType: args.Area.PointType,
			Data: []MapMathConversion.ArgsConversionGPS{
				{
					Longitude: args.Point.Longitude,
					Latitude:  args.Point.Latitude,
				},
			},
		})
		if err != nil {
			return false
		}
		if len(conversionRes) < 1 {
			return false
		}
		args.Point.Longitude, args.Point.Latitude = conversionRes[0].Longitude, conversionRes[0].Latitude
	}
	//检查点是否在范围内
	count := 0
	x1, y1 := args.Area.Points[0].Longitude, args.Area.Points[0].Latitude
	x1Part := (y1 > args.Point.Latitude) || ((x1-args.Point.Longitude > 0) && (y1 == args.Point.Latitude))
	var a = []float64{x1, y1}
	var points [][]float64
	for _, v := range args.Area.Points {
		vF := []float64{
			v.Longitude,
			v.Latitude,
		}
		points = append(points, vF)
	}
	p := append(points, a)
	for i := range p {
		if i == 0 {
			continue
		}
		p := p[i]
		x2, y2 := p[0], p[1]
		x2Part := (y2 > args.Point.Latitude) || ((x2 > args.Point.Longitude) && (y2 == args.Point.Latitude))
		if x2Part == x1Part {
			x1, y1 = x2, y2
			continue
		}
		mul := (x1-args.Point.Longitude)*(y2-args.Point.Latitude) - (x2-args.Point.Longitude)*(y1-args.Point.Latitude)
		if mul > 0 {
			count += 1
		} else {
			if mul < 0 {
				count -= 1
			}
		}
		x1, y1 = x2, y2
		x1Part = x2Part
	}
	if count == 2 || count == -2 {
		return true
	} else {
		return false
	}
}
