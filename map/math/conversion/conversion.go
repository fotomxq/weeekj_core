package MapMathConversion

import (
	"errors"
	"math"
)

//坐标系统转化模块
// WGS-84\GCJ-02\BD-09

// WGS84坐标系：即地球坐标系，国际上通用的坐标系。
// GCJ02坐标系：即火星坐标系，WGS84坐标系经加密后的坐标系。Google Maps，高德在用。
// BD09坐标系：即百度坐标系，GCJ02坐标系经加密后的坐标系。

const (
	XPi    = math.Pi * 3000.0 / 180.0
	OFFSET = 0.00669342162296594323
	AXIS   = 6378245.0
)

type ArgsConversionMapTypeInt struct {
	//源格式
	// 0 / 1 / 2
	// WGS-84 / GCJ-02 / BD-09
	SrcType int `json:"srcType"`
	//目标格式
	// 0 / 1 / 2
	// WGS-84 / GCJ-02 / BD-09
	DestType int `json:"destType"`
	//原数据
	Data []ArgsConversionGPS `json:"data"`
}

func ConversionMapTypeInt(args *ArgsConversionMapTypeInt) (result []ArgsConversionGPS, err error) {
	srcTypeStr := ConversionMapType(args.SrcType)
	destTypeStr := ConversionMapType(args.DestType)
	return Conversion(&ArgsConversion{
		SrcType:  srcTypeStr,
		DestType: destTypeStr,
		Data:     args.Data,
	})
}

//ArgsConversion 快速转化函数处理模块
type ArgsConversion struct {
	//源格式
	// WGS-84\GCJ-02\BD-09
	SrcType string `json:"srcType"`
	//目标格式
	// WGS-84\GCJ-02\BD-09
	DestType string `json:"destType"`
	//原数据
	Data []ArgsConversionGPS `json:"data"`
}

type ArgsConversionGPS struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

func Conversion(args *ArgsConversion) (result []ArgsConversionGPS, err error) {
	switch args.SrcType {
	case "WGS-84":
		switch args.DestType {
		case "WGS-84":
			for _, v := range args.Data {
				result = append(result, ArgsConversionGPS{
					Longitude: v.Longitude,
					Latitude:  v.Latitude,
				})
			}
		case "GCJ-02":
			for _, v := range args.Data {
				log, lat := WGS84toGCJ02(v.Longitude, v.Latitude)
				result = append(result, ArgsConversionGPS{
					Longitude: log,
					Latitude:  lat,
				})
			}
		case "BD-09":
			for _, v := range args.Data {
				log, lat := WGS84toBD09(v.Longitude, v.Latitude)
				result = append(result, ArgsConversionGPS{
					Longitude: log,
					Latitude:  lat,
				})
			}
		default:
			err = errors.New("un know dest type")
		}
	case "GCJ-02":
		switch args.DestType {
		case "WGS-84":
			for _, v := range args.Data {
				log, lat := GCJ02toWGS84(v.Longitude, v.Latitude)
				result = append(result, ArgsConversionGPS{
					Longitude: log,
					Latitude:  lat,
				})
			}
		case "GCJ-02":
			for _, v := range args.Data {
				result = append(result, ArgsConversionGPS{
					Longitude: v.Longitude,
					Latitude:  v.Latitude,
				})
			}
		case "BD-09":
			for _, v := range args.Data {
				log, lat := GCJ02toBD09(v.Longitude, v.Latitude)
				result = append(result, ArgsConversionGPS{
					Longitude: log,
					Latitude:  lat,
				})
			}
		default:
			err = errors.New("un know dest type")
		}
	case "BD-09":
		switch args.DestType {
		case "WGS-84":
			for _, v := range args.Data {
				log, lat := BD09toWGS84(v.Longitude, v.Latitude)
				result = append(result, ArgsConversionGPS{
					Longitude: log,
					Latitude:  lat,
				})
			}
		case "GCJ-02":
			for _, v := range args.Data {
				log, lat := BD09toGCJ02(v.Longitude, v.Latitude)
				result = append(result, ArgsConversionGPS{
					Longitude: log,
					Latitude:  lat,
				})
			}
		case "BD-09":
			for _, v := range args.Data {
				result = append(result, ArgsConversionGPS{
					Longitude: v.Longitude,
					Latitude:  v.Latitude,
				})
			}
		default:
			err = errors.New("un know dest type")
		}
	default:
		err = errors.New("un know src type")
		return
	}
	return
}

//BD09toGCJ02 百度坐标系->火星坐标系
func BD09toGCJ02(lon, lat float64) (float64, float64) {
	x := lon - 0.0065
	y := lat - 0.006

	z := math.Sqrt(x*x+y*y) - 0.00002*math.Sin(y*XPi)
	theta := math.Atan2(y, x) - 0.000003*math.Cos(x*XPi)

	gLon := z * math.Cos(theta)
	gLat := z * math.Sin(theta)

	return gLon, gLat
}

//GCJ02toBD09 火星坐标系->百度坐标系
func GCJ02toBD09(lon, lat float64) (float64, float64) {
	z := math.Sqrt(lon*lon+lat*lat) + 0.00002*math.Sin(lat*XPi)
	theta := math.Atan2(lat, lon) + 0.000003*math.Cos(lon*XPi)

	bdLon := z*math.Cos(theta) + 0.0065
	bdLat := z*math.Sin(theta) + 0.006

	return bdLon, bdLat
}

//WGS84toGCJ02 WGS84坐标系->火星坐标系
func WGS84toGCJ02(lon, lat float64) (float64, float64) {
	if isOutOFChina(lon, lat) {
		return lon, lat
	}

	mgLon, mgLat := delta(lon, lat)

	return mgLon, mgLat
}

//GCJ02toWGS84 火星坐标系->WGS84坐标系
func GCJ02toWGS84(lon, lat float64) (float64, float64) {
	if isOutOFChina(lon, lat) {
		return lon, lat
	}

	mgLon, mgLat := delta(lon, lat)

	return lon*2 - mgLon, lat*2 - mgLat
}

//BD09toWGS84 百度坐标系->WGS84坐标系
func BD09toWGS84(lon, lat float64) (float64, float64) {
	lon, lat = BD09toGCJ02(lon, lat)
	return GCJ02toWGS84(lon, lat)
}

//WGS84toBD09 WGS84坐标系->百度坐标系
func WGS84toBD09(lon, lat float64) (float64, float64) {
	lon, lat = WGS84toGCJ02(lon, lat)
	return GCJ02toBD09(lon, lat)
}

func delta(lon, lat float64) (float64, float64) {
	dlat := transformlat(lon-105.0, lat-35.0)
	dlon := transformlng(lon-105.0, lat-35.0)

	radlat := lat / 180.0 * math.Pi
	magic := math.Sin(radlat)
	magic = 1 - OFFSET*magic*magic
	sqrtmagic := math.Sqrt(magic)

	dlat = (dlat * 180.0) / ((AXIS * (1 - OFFSET)) / (magic * sqrtmagic) * math.Pi)
	dlon = (dlon * 180.0) / (AXIS / sqrtmagic * math.Cos(radlat) * math.Pi)

	mgLat := lat + dlat
	mgLon := lon + dlon

	return mgLon, mgLat
}

func transformlat(lon, lat float64) float64 {
	var ret = -100.0 + 2.0*lon + 3.0*lat + 0.2*lat*lat + 0.1*lon*lat + 0.2*math.Sqrt(math.Abs(lon))
	ret += (20.0*math.Sin(6.0*lon*math.Pi) + 20.0*math.Sin(2.0*lon*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lat*math.Pi) + 40.0*math.Sin(lat/3.0*math.Pi)) * 2.0 / 3.0
	ret += (160.0*math.Sin(lat/12.0*math.Pi) + 320*math.Sin(lat*math.Pi/30.0)) * 2.0 / 3.0
	return ret
}

func transformlng(lon, lat float64) float64 {
	var ret = 300.0 + lon + 2.0*lat + 0.1*lon*lon + 0.1*lon*lat + 0.1*math.Sqrt(math.Abs(lon))
	ret += (20.0*math.Sin(6.0*lon*math.Pi) + 20.0*math.Sin(2.0*lon*math.Pi)) * 2.0 / 3.0
	ret += (20.0*math.Sin(lon*math.Pi) + 40.0*math.Sin(lon/3.0*math.Pi)) * 2.0 / 3.0
	ret += (150.0*math.Sin(lon/12.0*math.Pi) + 300.0*math.Sin(lon/30.0*math.Pi)) * 2.0 / 3.0
	return ret
}

func isOutOFChina(lon, lat float64) bool {
	return !(lon > 73.66 && lon < 135.05 && lat > 3.86 && lat < 53.55)
}
