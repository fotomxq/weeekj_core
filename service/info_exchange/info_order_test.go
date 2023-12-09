package ServiceInfoExchange

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	ServiceOrder "gitee.com/weeekj/weeekj_core/v5/service/order"
	ServiceOrderMod "gitee.com/weeekj/weeekj_core/v5/service/order/mod"
	ServiceOrderWait "gitee.com/weeekj/weeekj_core/v5/service/order/wait"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
	"time"
)

var (
	newOrderData ServiceOrder.FieldsOrder
)

func TestInitInfoOrder(t *testing.T) {
	TestInitInfo(t)
	TestCreateInfo(t)
	TestPublishInfo(t)
	TestAuditInfo(t)
	TestPublishInfo(t)
}

func TestCreateInfoOrder(t *testing.T) {
	data, errCode, err := CreateInfoOrder(&ArgsCreateInfoOrder{
		OrgID:              newInfoData.OrgID,
		UserID:             newInfoData.UserID,
		BuyUserID:          newInfoData.UserID,
		CreateFrom:         0,
		Address:            newInfoData.Address,
		Des:                "测试购买信息服务",
		InfoID:             newInfoData.ID,
		PricePay:           false,
		TransportPayAfter:  false,
		ReferrerNationCode: "86",
		ReferrerPhone:      "17777777777",
		Params:             CoreSQLConfig.FieldsConfigsType{},
	})
	if err != nil {
		t.Error(err, ", ", errCode)
		return
	} else {
		t.Log(data)
	}
	//激活订单核验
	go ServiceOrder.Run()
	time.Sleep(time.Second * 5)
	//检查订单
	orderID, errCode, errMsg, err := ServiceOrderWait.CheckOrder(&ServiceOrderWait.ArgsCheckOrder{
		ID:     data.ID,
		OrgID:  data.OrgID,
		UserID: data.UserID,
	})
	if err != nil {
		t.Error(errCode, ", ", errMsg)
		t.Error(err)
		return
	}
	//获取订单信息
	orderData, err := ServiceOrder.GetByID(&ServiceOrder.ArgsGetByID{
		ID:     orderID,
		OrgID:  data.OrgID,
		UserID: data.UserID,
	})
	if err != nil {
		t.Error(err)
		return
	} else {
		t.Log("orderData: ", orderData)
		newOrderData = orderData
	}
}

func TestUpdateInfoOrderPrice(t *testing.T) {
	var priceList ServiceOrderMod.FieldsPrices
	for _, v := range newOrderData.PriceList {
		priceList = append(priceList, ServiceOrderMod.FieldsPrice{
			PriceType: v.PriceType,
			PayID:     v.PayID,
			PayFailed: v.PayFailed,
			IsPay:     v.IsPay,
			Price:     v.Price,
		})
	}
	priceList = append(priceList, ServiceOrderMod.FieldsPrice{
		PriceType: 3,
		PayID:     0,
		PayFailed: "",
		IsPay:     false,
		Price:     30,
	})
	err := UpdateInfoOrderPrice(&ArgsUpdateInfoOrderPrice{
		InfoID:    newInfoData.ID,
		OrgID:     newInfoData.OrgID,
		UserID:    newInfoData.UserID,
		OrgBindID: 0,
		PriceList: priceList,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearInfoOrder(t *testing.T) {
	TestDeleteInfo(t)
	TestClearInfo(t)
}
