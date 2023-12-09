package MapMathArea

import (
	"errors"
)

//本模块提供分区、点对应的计算工具

//GetAreaCenter 获取中心点
func GetAreaCenter(data *ParamsArea) (float64, float64, error) {
	//没有XY数据，返回失败
	if len(data.Points) < 1 {
		return 0, 0, errors.New("no have x y data")
	}
	//是一个点，直接反馈点数据
	if len(data.Points) == 1 {
		return data.Points[0].Longitude, data.Points[0].Latitude, nil
	}
	//计算中心点并反馈
	xyMin := ParamsAreaPoint{}
	xyMax := ParamsAreaPoint{}
	center := ParamsAreaPoint{}
	for _, v2 := range data.Points {
		//如果更小
		if v2.Longitude+v2.Latitude < xyMin.Longitude+xyMin.Latitude {
			xyMin = v2
		}
		//如果更大
		if v2.Longitude+v2.Latitude > xyMax.Longitude+xyMax.Latitude {
			xyMax = v2
		}
		//计算中心点
		center = ParamsAreaPoint{
			Longitude: xyMax.Longitude - xyMin.Longitude,
			Latitude:  xyMax.Latitude - xyMin.Latitude,
		}
	}
	//反馈
	return center.Longitude, center.Latitude, nil
}

//ArgsGetCircleByPoint 给予一个点坐标和半径，计算圆的范围点数据参数
type ArgsGetCircleByPoint struct {
	//圆心
	Point ParamsAreaPoint
	//半径
	Radius float64
}

//GetCircleByPoint 给予一个点坐标和半径，计算圆的范围点数据
func GetCircleByPoint(args *ArgsGetCircleByPoint) (area []ParamsAreaPoint) {
	area = append(area, ParamsAreaPoint{
		Longitude: args.Point.Longitude - args.Radius,
		Latitude:  args.Point.Latitude - args.Radius,
	})
	area = append(area, ParamsAreaPoint{
		Longitude: args.Point.Longitude - args.Radius,
		Latitude:  args.Point.Latitude + args.Radius,
	})
	area = append(area, ParamsAreaPoint{
		Longitude: args.Point.Longitude + args.Radius,
		Latitude:  args.Point.Latitude + args.Radius,
	})
	area = append(area, ParamsAreaPoint{
		Longitude: args.Point.Longitude + args.Radius,
		Latitude:  args.Point.Latitude - args.Radius,
	})
	return
}
