package ServiceOrderWait

import (
	"fmt"
	CoreSQLAddress "github.com/fotomxq/weeekj_core/v5/core/sql/address"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	ServiceOrderWaitFields "github.com/fotomxq/weeekj_core/v5/service/order/wait_fields"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
	"time"
)

var (
	isInit  = false
	newData ServiceOrderWaitFields.FieldsWait
)

func TestInit(t *testing.T) {
	if isInit {
		return
	}
	isInit = true
	ToolsTest.ConfigDirAppend = fmt.Sprint("../", ToolsTest.ConfigDirAppend)
	ToolsTest.Init(t)
}

func TestCreateOrder(t *testing.T) {
	data, errCode, err := CreateOrder(&ArgsCreateOrder{
		SystemMark:         "test",
		OrgID:              0,
		UserID:             0,
		CreateFrom:         0,
		AddressFrom:        CoreSQLAddress.FieldsAddress{},
		AddressTo:          CoreSQLAddress.FieldsAddress{},
		Goods:              []ServiceOrderWaitFields.FieldsGood{},
		Exemptions:         []ServiceOrderWaitFields.FieldsExemption{},
		NeedAllowAutoAudit: false,
		AllowAutoAudit:     false,
		TransportAllowAuto: false,
		TransportTaskAt:    time.Time{},
		TransportPayAfter:  false,
		PriceList: []ServiceOrderWaitFields.FieldsPrice{
			{
				PriceType: 0,
				IsPay:     false,
				Price:     15,
			},
		},
		PricePay:    false,
		NeedExPrice: false,
		Currency:    0,
		Des:         "",
		Logs:        []ServiceOrderWaitFields.FieldsLog{},
		Params:      []CoreSQLConfig.FieldsConfigType{},
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		newData = data
	}
}

func TestCheckOrder(t *testing.T) {
	orderID, errCode, errMsg, err := CheckOrder(&ArgsCheckOrder{
		ID:    newData.ID,
		OrgID: 0,
	})
	ToolsTest.ReportError(t, err)
	if err != nil {
		t.Error(errCode, errMsg)
	}
	if err == nil {
		t.Log("orderID: ", orderID)
	}
}
