package TMSTransport

import (
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newBindToMallData FieldsBindToMall
)

func TestBindToMallInit(t *testing.T) {
	TestTransportInit(t)
}

func TestSetBindToMall(t *testing.T) {
	newUserInfo, _, err := TestOrg.CreateUser(t)
	if err != nil {
		return
	}
	bindData, err := TestOrg.CreateBind(t, newUserInfo)
	if err != nil {
		return
	}
	ToolsTest.ReportData(t, err, bindData)
	newBindData2, err := SetBind(&ArgsSetBind{
		OrgID:          TestOrg.OrgData.ID,
		BindID:         bindData.ID,
		MapAreaID:      newMapArea.ID,
		MoreMapAreaIDs: []int64{},
		Params:         []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, newBindData2)
	t.Log("new bind data, bind id: ", newBindData2.BindID)
	err = SetBindToMall(&ArgsSetBindToMall{
		OrgID:      TestOrg.OrgData.ID,
		BindID:     newBindData.BindID,
		BindMallID: 1,
	})
	ToolsTest.ReportError(t, err)
}

func TestGetBindToMallList(t *testing.T) {
	dataList, dataCount, err := GetBindToMallList(&ArgsGetBindToMallList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:      -1,
		BindID:     -1,
		BindMallID: -1,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil && len(dataList) > 0 {
		newBindToMallData = dataList[0]
	}
	newTransportData, _, err = CreateTransport(&ArgsCreateTransport{
		OrgID:  TestOrg.OrgData.ID,
		BindID: 0,
		InfoID: 0,
		UserID: 0,
		FromAddress: CoreSQLAddress.FieldsAddress{
			Country:    86,
			Province:   10000,
			City:       10000,
			Address:    "测试地址",
			MapType:    0,
			Longitude:  15,
			Latitude:   15,
			Name:       "测试名称",
			NationCode: "86",
			Phone:      "17777777777",
		},
		ToAddress: CoreSQLAddress.FieldsAddress{
			Country:    86,
			Province:   10000,
			City:       10000,
			Address:    "测试地址",
			MapType:    0,
			Longitude:  65,
			Latitude:   65,
			Name:       "测试名称",
			NationCode: "86",
			Phone:      "17777777777",
		},
		OrderID: 0,
		Goods: []FieldsTransportGood{
			{
				System: "mall",
				ID:     1,
				Mark:   "",
				Name:   "测试商品",
				Count:  2,
			},
		},
		Weight:    20,
		Length:    10,
		Width:     20,
		Currency:  0,
		Price:     30,
		PayFinish: false,
		TaskAt:    "",
		Params:    []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, newTransportData)
	dataList2, dataCount2, err := GetLogList(&ArgsGetLogList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:           -1,
		BindID:          -1,
		TransportID:     newTransportData.ID,
		TransportBindID: -1,
		Mark:            "",
		Search:          "",
	})
	ToolsTest.ReportDataList(t, err, dataList2, dataCount2)
	if newTransportData.BindID != newBindToMallData.BindID {
		t.Error("bind not ok, newTransportData.BindID: ", newTransportData.BindID, ", newBindToMallData.BindID: ", newBindToMallData.BindID)
	}
}

func TestDeleteBindToMall(t *testing.T) {
	err := DeleteBindToMall(&ArgsDeleteBindToMall{
		ID:    newBindToMallData.ID,
		OrgID: TestOrg.OrgData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestBindToMallClear(t *testing.T) {
	TestTransportClear(t)
}
