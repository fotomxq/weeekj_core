package MapTMap

import "math"

const (
	a        = 6378245.0
	ee       = 0.00669342162296594323
	pi       = math.Pi
	xPi      = math.Pi * 3000.0 / 180.0
	deltaLat = 0.01
	deltaLng = 0.01
)

func WGS84ToGcj02(lng, lat float64) (gcj02Lng, gcj02Lat float64) {
	if outOfChina(lat, lng) {
		return lng, lat
	}
	var dLat = transformLat(lng-105.0, lat-35.0)
	var dLng = transformLng(lng-105.0, lat-35.0)
	var radLat = lat / 180.0 * pi
	var magic = math.Sin(radLat)
	magic = 1 - ee*magic*magic
	var sqrtMagic = math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((a * (1 - ee)) / (magic * sqrtMagic) * pi)
	dLng = (dLng * 180.0) / (a / sqrtMagic * math.Cos(radLat) * pi)
	return lng + dLng, lat + dLat
}

func outOfChina(lat, lng float64) bool {
	if lng < 72.004 || lng > 137.8347 {
		return true
	}
	if lat < 0.8293 || lat > 55.8271 {
		return true
	}
	return false
}

func transformLat(x, y float64) float64 {
	var lat = y
	var lng = x
	var ret = -100.0 + 2.0*lng + 3.0*xPi*lat + 0.2*xPi*xPi*lat + 0.1*xPi*lat*lat + 0.2*math.Sqrt(math.Abs(lng))*math.Sin(6.0*xPi*lat) + 0.2*math.Sqrt(math.Abs(lng))*math.Sin(2.0*xPi*lat) + 0.2*math.Sqrt(math.Abs(lat))*math.Sin(lng*xPi) + 0.4*math.Sin(lng*xPi) - 1.25*math.Sin(2.0*lat*xPi)
	ret += deltaLat * math.Sin(lat*pi)
	return lat + ret
}

func transformLng(x, y float64) float64 {
	var lat = y
	var lng = x
	var ret = 300.0 + lng + 2.0*lat + 0.1*lng*lng + 0.1*lng*lat + 0.1*math.Sqrt(math.Abs(lng))*math.Sin(6.0*xPi*lat) + 0.1*math.Sqrt(math.Abs(lng))*math.Sin(2.0*xPi*lat) + 0.1*math.Sqrt(math.Abs(lat))*math.Sin(lng*xPi) + 0.2*math.Sin(lng*xPi) + 0.1*math.Sin(lat*xPi)*math.Sqrt(math.Abs(lng))
	ret += deltaLng * math.Sin(lng*pi)
	return lng + ret
}
