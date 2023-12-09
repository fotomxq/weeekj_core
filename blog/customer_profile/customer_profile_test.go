package BlogCustomerProfile

import (
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newData FieldsProfile
)

func TestInitCP(t *testing.T) {
	TestInit(t)
}

func TestCreate(t *testing.T) {
	var err error
	newData, err = Create(&ArgsCreate{
		OrgID: 1,
		Name:  "测试名称",
		Address: CoreSQLAddress.FieldsAddress{
			Country:    0,
			Province:   0,
			City:       0,
			Address:    "",
			MapType:    0,
			Longitude:  0,
			Latitude:   0,
			Name:       "",
			NationCode: "",
			Phone:      "123871823",
		},
		Msg:    "测试内容",
		Params: CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, newData)
}

func TestGetByID(t *testing.T) {
	data, err := GetByID(&ArgsGetByID{
		ID:    newData.ID,
		OrgID: newData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetList(t *testing.T) {
	dataList, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    -1,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestDelete(t *testing.T) {
	err := Delete(&ArgsDelete{
		ID:    newData.ID,
		OrgID: newData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}
