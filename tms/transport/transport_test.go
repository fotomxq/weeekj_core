package TMSTransport

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	newTransportData FieldsTransport
)

func TestTransportInit(t *testing.T) {
	TestBindInit(t)
	TestSetBind(t)
}

func TestCreateTransport(t *testing.T) {
	var err error
	newTransportData, _, err = CreateTransport(&ArgsCreateTransport{
		OrgID:  TestOrg.OrgData.ID,
		BindID: newBindData.BindID,
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
		Price:     30,
		PayFinish: false,
		Params:    []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, newTransportData)
}

func TestGetTransport(t *testing.T) {
	data, err := GetTransport(&ArgsGetTransport{
		ID:    newTransportData.ID,
		OrgID: newTransportData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetTransportList(t *testing.T) {
	dataList, dataCount, err := GetTransportList(&ArgsGetTransportList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:       -1,
		BindID:      -1,
		InfoID:      -1,
		UserID:      -1,
		OrderID:     -1,
		SN:          0,
		SNDay:       0,
		Status:      nil,
		NeedIsPay:   false,
		IsPay:       false,
		PayID:       0,
		IsFinishAt:  false,
		TimeBetween: CoreSQLTime.DataCoreTime{},
		IsRemove:    false,
		IsHistory:   false,
		Search:      "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetTransports(t *testing.T) {
	dataList, err := GetTransports(&ArgsGetTransports{
		IDs:        []int64{newTransportData.ID},
		HaveRemove: false,
		OrgID:      0,
	})
	ToolsTest.ReportData(t, err, dataList)
}

func TestPayTransport(t *testing.T) {
	payData, errCode, err := PayTransport(&ArgsPayTransport{
		ID:     newTransportData.ID,
		OrgID:  newTransportData.OrgID,
		BindID: newTransportData.BindID,
		PaymentChannel: CoreSQLFrom.FieldsFrom{
			System: "cash",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log("pay data: ", payData)
	}
}

func TestPayForceTransport(t *testing.T) {
	err := PayForceTransport(&ArgsPayForceTransport{
		ID:     newTransportData.ID,
		OrgID:  newTransportData.OrgID,
		BindID: newTransportData.BindID,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateTransportComment(t *testing.T) {
	err := UpdateTransportComment(&ArgsUpdateTransportComment{
		ID:     newTransportData.ID,
		InfoID: 0,
		UserID: 0,
		Level:  3,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateTransportGPS(t *testing.T) {
	err := UpdateTransportGPS(&ArgsUpdateTransportGPS{
		ID:        newTransportData.ID,
		OrgID:     newTransportData.OrgID,
		MapType:   0,
		Longitude: 55,
		Latitude:  56,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateTransportPick(t *testing.T) {
	err := UpdateTransportPick(&ArgsUpdateTransportPick{
		ID:     newTransportData.ID,
		OrgID:  newTransportData.OrgID,
		BindID: newTransportData.BindID,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateTransportSend(t *testing.T) {
	err := UpdateTransportSend(&ArgsUpdateTransportSend{
		ID:     newTransportData.ID,
		OrgID:  newTransportData.OrgID,
		BindID: newTransportData.BindID,
	})
	ToolsTest.ReportError(t, err)
}

func TestUpdateTransportFinish(t *testing.T) {
	err := UpdateTransportFinish(&ArgsUpdateTransportFinish{
		ID:     newTransportData.ID,
		OrgID:  newTransportData.OrgID,
		BindID: newTransportData.BindID,
	})
	ToolsTest.ReportError(t, err)
}

func TestGetAnalysisBindGoods(t *testing.T) {
	data, err := GetAnalysisBindGoods(&ArgsGetAnalysisBindGoods{
		OrgID:   newTransportData.OrgID,
		BindIDs: []int64{newTransportData.BindID},
		TimeBetween: CoreSQLTime.DataCoreTime{
			MinTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().SubDays(2).Time),
			MaxTime: CoreFilter.GetISOByTime(CoreFilter.GetNowTimeCarbon().Time),
		},
	})
	ToolsTest.ReportData(t, err, data)
}

func TestDeleteTransport(t *testing.T) {
	err := DeleteTransport(&ArgsDeleteTransport{
		ID:     newTransportData.ID,
		OrgID:  newTransportData.OrgID,
		BindID: newTransportData.BindID,
	})
	ToolsTest.ReportError(t, err)
}

func TestTransportClear(t *testing.T) {
	TestDeleteBind(t)
	TestBindClear(t)
}
