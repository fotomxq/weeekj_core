package ServiceOrder

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLAddress "gitee.com/weeekj/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	ServiceOrderWait "gitee.com/weeekj/weeekj_core/v5/service/order/wait"
	ServiceOrderWaitFields "gitee.com/weeekj/weeekj_core/v5/service/order/wait_fields"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	"testing"
	"time"
)

var (
	newOrderData FieldsOrder
)

func TestInitCreate(t *testing.T) {
	TestInit(t)
}

func TestCreate(t *testing.T) {
	//随机数量
	randCount := CoreFilter.GetRandNumber(1, 99)
	//创建订单
	data, errCode, err := ServiceOrderWait.CreateOrder(&ServiceOrderWait.ArgsCreateOrder{
		SystemMark: "test",
		OrgID:      TestOrg.OrgData.ID,
		UserID:     TestOrg.UserInfo.ID,
		CreateFrom: 1,
		AddressFrom: CoreSQLAddress.FieldsAddress{
			Country:    86,
			Province:   100,
			City:       10010,
			Address:    "测试地址",
			MapType:    0,
			Longitude:  15,
			Latitude:   15,
			Name:       "",
			NationCode: "86",
			Phone:      "17635705566",
		},
		AddressTo: CoreSQLAddress.FieldsAddress{
			Country:    86,
			Province:   100,
			City:       10010,
			Address:    "太原小店区",
			MapType:    0,
			Longitude:  15,
			Latitude:   15,
			Name:       "",
			NationCode: "86",
			Phone:      "17635705566",
		},
		Goods: []ServiceOrderWaitFields.FieldsGood{
			{
				From: CoreSQLFrom.FieldsFrom{
					System: "mall",
					ID:     123,
					Mark:   "",
					Name:   "",
				},
				Count: int64(randCount),
				Price: int64(randCount * randCount),
				Exemptions: ServiceOrderWaitFields.FieldsExemptions{
					{
						System:   "test",
						ConfigID: 0,
						Name:     "",
						Des:      "",
						Count:    1,
						Price:    5,
					},
				},
			},
		},
		Exemptions:         nil,
		NeedAllowAutoAudit: false,
		AllowAutoAudit:     false,
		TransportAllowAuto: false,
		TransportTaskAt:    time.Time{},
		TransportPayAfter:  true,
		PriceList: []ServiceOrderWaitFields.FieldsPrice{
			{
				PriceType: 0,
				IsPay:     false,
				Price:     25,
			},
		},
		PricePay:           false,
		NeedExPrice:        false,
		Currency:           86,
		Des:                "",
		Logs:               []ServiceOrderWaitFields.FieldsLog{},
		ReferrerNationCode: "",
		ReferrerPhone:      "",
		Params:             []CoreSQLConfig.FieldsConfigType{},
	})
	if err != nil {
		t.Error(errCode, err)
		return
	} else {
		t.Log("price: ", data.Price)
	}
	//启动转化run
	//go runCreate()
	time.Sleep(time.Second * 1)
	//检查状态
	if orderID, errCode, errMsg, err := ServiceOrderWait.CheckOrder(&ServiceOrderWait.ArgsCheckOrder{
		ID:    data.ID,
		OrgID: 0,
	}); err != nil {
		t.Error(err, errCode, errMsg)
	} else {
		t.Log("orderID: ", orderID)
	}
	TestGetList(t)
}

func TestGetByID(t *testing.T) {
	data, err := GetByID(&ArgsGetByID{
		ID:     newOrderData.ID,
		OrgID:  0,
		UserID: 0,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestClearCreate(t *testing.T) {
	TestClear(t)
}
