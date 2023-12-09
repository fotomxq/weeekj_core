package MapMathConversion

//ConversionMapType 转化mapType
// 0 / 1 / 2
// WGS-84 / GCJ-02 / BD-09
func ConversionMapType(mapType int) string {
	switch mapType {
	case 0:
		return "WGS-84"
	case 1:
		return "GCJ-02"
	case 2:
		return "BD-09"
	default:
		return "WGS-84"
	}
}
