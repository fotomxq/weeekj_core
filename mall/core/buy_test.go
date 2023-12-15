package MallCore

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	ServiceOrderWaitFields "github.com/fotomxq/weeekj_core/v5/service/order/wait_fields"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	buyTestAddress = CoreSQLAddress.FieldsAddress{
		Country:    86,
		Province:   10000,
		City:       10000,
		Address:    "测试地址",
		MapType:    0,
		Longitude:  3,
		Latitude:   3,
		Name:       "测试地址名称",
		NationCode: "86",
		Phone:      "17555555555",
	}
)

func TestInitBuy(t *testing.T) {
	TestInitAudit(t)
	TestCreateAudit(t)
	TestUpdateAuditPassing(t)
}

func TestGetProductPrice(t *testing.T) {
	products := []ArgsGetProductPriceProduct{
		{
			ID:       newProductData.ID,
			BuyCount: 1,
		},
	}
	data, errCode, err := GetProductPrice(&ArgsGetProductPrice{
		Products:    products,
		OrgID:       newProductData.OrgID,
		UserID:      TestOrg.UserInfo.ID,
		UserSubID:   0,
		UserTicket:  []int64{},
		UseIntegral: false,
		Address:     buyTestAddress,
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log(data)
	}
}

// 自动化测试会员价格
func TestGetProductPriceUserSub(t *testing.T) {
	products := []ArgsGetProductPriceProduct{
		{
			ID:       newProductData.ID,
			BuyCount: 1,
		},
	}
	data, errCode, err := GetProductPrice(&ArgsGetProductPrice{
		Products:    products,
		OrgID:       newProductData.OrgID,
		UserID:      TestOrg.UserInfo.ID,
		UserSubID:   newUserSubConfig.ID,
		UserTicket:  []int64{},
		UseIntegral: false,
		Address:     buyTestAddress,
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log(data)
		if data.LastPrice != newProductData.UserSubPrice[0].Price {
			t.Error("会员价格异常: ", data.LastPrice, ", 应该价格: ", newProductData.UserSubPrice[0].Price)
		}
	}
	products = []ArgsGetProductPriceProduct{
		{
			ID:       newProductData.ID,
			BuyCount: 3,
		},
	}
	data, errCode, err = GetProductPrice(&ArgsGetProductPrice{
		Products:    products,
		OrgID:       newProductData.OrgID,
		UserID:      TestOrg.UserInfo.ID,
		UserSubID:   newUserSubConfig.ID,
		UserTicket:  []int64{},
		UseIntegral: false,
		Address:     buyTestAddress,
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log(data)
		if data.LastPrice != newProductData.UserSubPrice[0].Price*3 {
			t.Error("会员价格异常: ", data.LastPrice, ", 应该价格: ", newProductData.UserSubPrice[0].Price*3)
		}
	}
}

// 自动化测试票据价格
func TestGetProductPriceUserTicket(t *testing.T) {
	products := []ArgsGetProductPriceProduct{
		{
			ID:       newProductData.ID,
			BuyCount: 1,
		},
	}
	data, errCode, err := GetProductPrice(&ArgsGetProductPrice{
		Products:    products,
		OrgID:       newProductData.OrgID,
		UserID:      TestOrg.UserInfo.ID,
		UserSubID:   0,
		UserTicket:  []int64{newTicketConfig.ID},
		UseIntegral: false,
		Address:     buyTestAddress,
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log(data)
		if data.LastPrice != 0 {
			t.Error("票据价格异常: ", data.LastPrice, ", 应该价格: ", 0)
		}
	}
	//3张等价测试
	products = []ArgsGetProductPriceProduct{
		{
			ID:       newProductData.ID,
			BuyCount: 3,
		},
	}
	data, errCode, err = GetProductPrice(&ArgsGetProductPrice{
		Products:    products,
		OrgID:       newProductData.OrgID,
		UserID:      TestOrg.UserInfo.ID,
		UserSubID:   0,
		UserTicket:  []int64{newTicketConfig.ID},
		UseIntegral: false,
		Address:     buyTestAddress,
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log(data)
		if data.LastPrice != 0 {
			t.Error("票据价格异常: ", data.LastPrice, ", 应该价格: ", 0)
		}
	}
	//超出用户持有测试
	products = []ArgsGetProductPriceProduct{
		{
			ID:       newProductData.ID,
			BuyCount: 4,
		},
	}
	data, errCode, err = GetProductPrice(&ArgsGetProductPrice{
		Products:    products,
		OrgID:       newProductData.OrgID,
		UserID:      TestOrg.UserInfo.ID,
		UserSubID:   0,
		UserTicket:  []int64{newTicketConfig.ID},
		UseIntegral: false,
		Address:     buyTestAddress,
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log(data)
		if data.LastPrice != newProductData.PriceReal {
			t.Error("票据价格异常: ", data.LastPrice, ", 应该价格: ", newProductData.PriceReal, ", ", "用户持有票据: ", data.Tickets[0].UserCount, ", ", "使用张数: ", data.Tickets[0].NeedCount)
			t.Log("折扣信息: ", data.Exemptions)
		}
	}
}

func TestCreateOrder(t *testing.T) {
	products := []ArgsGetProductPriceProduct{
		{
			ID:       newProductData.ID,
			BuyCount: 3,
		},
	}
	data, errCode, err := CreateOrder(&ArgsCreateOrder{
		OrgID:      newProductData.OrgID,
		UserID:     TestOrg.UserInfo.ID,
		CreateFrom: 0,
		Address: CoreSQLAddress.FieldsAddress{
			Country:    86,
			Province:   10000,
			City:       10000,
			Address:    "测试地址",
			MapType:    0,
			Longitude:  3,
			Latitude:   3,
			Name:       "测试地址名称",
			NationCode: "86",
			Phone:      "17555555555",
		},
		TransportTaskAt:       CoreFilter.GetISOByTime(CoreFilter.GetNowTime()),
		Des:                   "测试订单备注",
		Products:              products,
		UserSubID:             0,
		UserTicket:            []int64{},
		UseIntegral:           false,
		PricePay:              false,
		TransportPayAfter:     false,
		SkipProductCountLimit: false,
		ReferrerNationCode:    "",
		ReferrerPhone:         "",
		OtherPriceList:        ServiceOrderWaitFields.FieldsPrices{},
		Params:                []CoreSQLConfig.FieldsConfigType{},
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log(data)
		t.Log("Price: ", data.Price)
		t.Log("PriceTotal: ", data.PriceTotal)
	}
}

func TestClearBuy(t *testing.T) {
	//TestDeleteAudit(t)
	TestClearAudit(t)
}
