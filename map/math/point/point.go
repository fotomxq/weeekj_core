package MapMathPoint

import (
	"errors"
	MapMathArgs "gitee.com/weeekj/weeekj_core/v5/map/math/args"
	MapMathConversion "gitee.com/weeekj/weeekj_core/v5/map/math/conversion"
	"math"
)

// ArgsGetDistance 计算两个坐标之间的距离参数
type ArgsGetDistance struct {
	//出发点
	StartPoint MapMathArgs.ParamsPoint `json:"startPoint"`
	//目的地
	EndPoint MapMathArgs.ParamsPoint `json:"endPoint"`
}

// GetDistance 计算两个坐标之间的距离
// @result distanceM 距离米
func GetDistance(args *ArgsGetDistance) (distanceM int64, err error) {
	//换算数据
	var startPoint []MapMathConversion.ArgsConversionGPS
	startPoint, err = MapMathConversion.Conversion(&MapMathConversion.ArgsConversion{
		SrcType:  args.StartPoint.PointType,
		DestType: "WGS-84",
		Data: []MapMathConversion.ArgsConversionGPS{
			{
				Longitude: args.StartPoint.Longitude,
				Latitude:  args.StartPoint.Latitude,
			},
		},
	})
	if err != nil {
		return
	}
	var endPoint []MapMathConversion.ArgsConversionGPS
	endPoint, err = MapMathConversion.Conversion(&MapMathConversion.ArgsConversion{
		SrcType:  args.EndPoint.PointType,
		DestType: "WGS-84",
		Data: []MapMathConversion.ArgsConversionGPS{
			{
				Longitude: args.EndPoint.Longitude,
				Latitude:  args.EndPoint.Latitude,
			},
		},
	})
	if err != nil {
		return
	}
	//计算距离
	if len(startPoint) < 1 || len(endPoint) < 1 {
		err = errors.New("no data")
		return
	}
	distanceM = int64(math.Abs(math.Sqrt(math.Pow(endPoint[0].Longitude-startPoint[0].Longitude, 2) + math.Pow(endPoint[0].Latitude-startPoint[0].Latitude, 2))))
	return
}

// GetDistanceM 计算两个坐标之间的米数
func GetDistanceM(args *ArgsGetDistance) (distance float64, err error) {
	//换算数据
	var startPoint []MapMathConversion.ArgsConversionGPS
	startPoint, err = MapMathConversion.Conversion(&MapMathConversion.ArgsConversion{
		SrcType:  args.StartPoint.PointType,
		DestType: "WGS-84",
		Data: []MapMathConversion.ArgsConversionGPS{
			{
				Longitude: args.StartPoint.Longitude,
				Latitude:  args.StartPoint.Latitude,
			},
		},
	})
	if err != nil {
		return
	}
	var endPoint []MapMathConversion.ArgsConversionGPS
	endPoint, err = MapMathConversion.Conversion(&MapMathConversion.ArgsConversion{
		SrcType:  args.EndPoint.PointType,
		DestType: "WGS-84",
		Data: []MapMathConversion.ArgsConversionGPS{
			{
				Longitude: args.EndPoint.Longitude,
				Latitude:  args.EndPoint.Latitude,
			},
		},
	})
	if err != nil {
		return
	}
	//计算距离
	if len(startPoint) < 1 || len(endPoint) < 1 {
		err = errors.New("no data")
		return
	}
	radius := 6378.137
	rad := math.Pi / 180.0
	endPoint[0].Latitude = endPoint[0].Latitude * rad
	endPoint[0].Longitude = endPoint[0].Longitude * rad
	startPoint[0].Latitude = startPoint[0].Latitude * rad
	startPoint[0].Longitude = startPoint[0].Longitude * rad
	theta := startPoint[0].Longitude - endPoint[0].Longitude
	dist := math.Acos(math.Sin(endPoint[0].Latitude)*math.Sin(startPoint[0].Latitude) + math.Cos(endPoint[0].Latitude)*math.Cos(startPoint[0].Latitude)*math.Cos(theta))
	distance = dist * radius
	distance = math.Abs(distance)
	return
}
