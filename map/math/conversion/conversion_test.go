package MapMathConversion

import "testing"

func TestWGS84toGCJ02(t *testing.T) {
	//BD: 112.48392,37.860648
	//GCJ02: 37.854640,112.477598
	/** 转化结果
	WGS84: 37.85464, 112.477598
	GCJ02: 37.85464, 112.477598
	BD09: 37.86147011217083, 112.48350635362486
	BD09MC: 4214765.432404487, 4335919424154
	EPSG3857: 4213959.2489626855, NaN
	*/
	//原始数据 16bit：040fb57c,0c11154c
	// bit10: 68 138364, 202 446156
	// bit10/30000: 2271.2788,6748.2052
	// bit10/30000/60: 37.854646,112.470086
	// 公式: 2232.7658’=(22X60+32.7658)X30000=40582974
	// 22x60+71.2788 / 57x60+48.2052
	//实际设备测试经纬度数据：2271.2788,6748.2052
	//手动换算：23.18798,67.80342
	/** 该定位转化数据
	WGS84: 2271.2788, 6748.2052
	GCJ02: 2271.2788, 6748.2052
	BD09: 2271.294303181935, 6748.208175420042
	BD09MC: 252842076.41736862, 5.012590022572102e+26
	EPSG3857: 20037508.342789244, NaN
	*/
	//将GCJ02转为WGS84
	a1, a2 := WGS84toGCJ02(112.470086, 37.854646)
	t.Log(a1, ",", a2)
}

func TestConversion(t *testing.T) {
	a1, err := Conversion(&ArgsConversion{
		SrcType:  "GCJ-02",
		DestType: "WGS-84",
		Data: []ArgsConversionGPS{
			{
				Longitude: 112.51562,
				Latitude:  37.85929,
			},
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(a1)
	}
}
