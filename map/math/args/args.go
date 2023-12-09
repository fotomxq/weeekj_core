package MapMathArgs

type ParamsPoint struct {
	//坐标制式
	// 0 / 1 / 2
	// WGS-84 / GCJ-02 / BD-09
	PointType string `json:"pointType"`
	//坐标位置
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
