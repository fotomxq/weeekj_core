package MapAMap

import (
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

func TestInitAddress(t *testing.T) {
	TestInit(t)
}

func TestGetGPSByAddress(t *testing.T) {
	data, err := GetGPSByAddress(&ArgsGetGPSByAddress{
		Address: "山西省太原市小店区府园东居3-301",
		City:    0,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetAddressByGPS(t *testing.T) {
	data, err := GetAddressByGPS(&ArgsGetAddressByGPS{
		Longitude: 112.53353,
		Latitude:  37.84298,
	})
	ToolsTest.ReportData(t, err, data)
}
