package MapMathArea

type ParamsArea struct {
	//锁定ID
	ID int64 `json:"id"`
	//坐标制式
	// WGS-84\GCJ-02\BD-09
	PointType string `json:"pointType"`
	//划区
	Points []ParamsAreaPoint `json:"points"`
}

//ParamsAreaPoint 分区专用坐标点
type ParamsAreaPoint struct {
	//坐标位置
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
