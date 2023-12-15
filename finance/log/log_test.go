package FinanceLog

import (
	CoreCurrency "github.com/fotomxq/weeekj_core/v5/core/currency"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	isInit = false
)

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
	}
	isInit = true
}

func TestSetHash(t *testing.T) {
	SetHash("test_hash_123")
}

func TestGetHash(t *testing.T) {
	data := GetHash()
	t.Log(data)
}

func TestCreate(t *testing.T) {
	step := 1
	for {
		if step > 50 {
			break
		}
		priceInt := CoreFilter.GetRandNumber(100, 900000)
		var price = CoreFilter.GetInt64ByFloat64(float64(priceInt) * 1.35)
		if err := Create(&ArgsLogCreate{
			PayID:    123,
			Hash:     "123",
			Key:      "",
			Status:   1,
			Currency: CoreCurrency.GetID("CNY"),
			Price:    price,
			PaymentCreate: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     221,
				Mark:   "",
				Name:   "来源用户昵称",
			},
			PaymentChannel: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     222,
				Mark:   "",
				Name:   "来源用户昵称",
			},
			PaymentFrom: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     223,
				Mark:   "",
				Name:   "来源用户昵称",
			},
			TakeCreate: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     224,
				Mark:   "",
				Name:   "来源用户昵称",
			},
			TakeChannel: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     225,
				Mark:   "",
				Name:   "来源用户昵称",
			},
			TakeFrom: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     226,
				Mark:   "",
				Name:   "来源用户昵称",
			},
			CreateInfo: CoreSQLFrom.FieldsFrom{},
			Des:        "",
		}); err != nil {
			t.Error(err)
		}
		step += 1
	}
}

func TestGetList(t *testing.T) {
	t.Log("test ", "empty")
	data, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "_id",
			Desc: false,
		},
		Key:            "",
		Status:         []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		PaymentCreate:  CoreSQLFrom.FieldsFrom{},
		PaymentChannel: CoreSQLFrom.FieldsFrom{},
		PaymentFrom:    CoreSQLFrom.FieldsFrom{},
		TakeCreate:     CoreSQLFrom.FieldsFrom{},
		TakeChannel:    CoreSQLFrom.FieldsFrom{},
		TakeFrom:       CoreSQLFrom.FieldsFrom{},
		CreateInfo:     CoreSQLFrom.FieldsFrom{},
		TimeBetween:    CoreSQLTime.FieldsCoreTime{},
		IsHistory:      false,
		Search:         "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data, dataCount)
	}
	t.Log("test ", "from user")
	data, dataCount, err = GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "_id",
			Desc: false,
		},
		Key:            "",
		Status:         []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		PaymentCreate:  CoreSQLFrom.FieldsFrom{},
		PaymentChannel: CoreSQLFrom.FieldsFrom{},
		PaymentFrom:    CoreSQLFrom.FieldsFrom{},
		TakeCreate:     CoreSQLFrom.FieldsFrom{},
		TakeChannel:    CoreSQLFrom.FieldsFrom{},
		TakeFrom:       CoreSQLFrom.FieldsFrom{},
		CreateInfo:     CoreSQLFrom.FieldsFrom{},
		TimeBetween:    CoreSQLTime.FieldsCoreTime{},
		IsHistory:      false,
		Search:         "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data, dataCount)
	}
	t.Log("test ", "to finance")
	data, dataCount, err = GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "_id",
			Desc: false,
		},
		Key:            "",
		Status:         []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		PaymentCreate:  CoreSQLFrom.FieldsFrom{},
		PaymentChannel: CoreSQLFrom.FieldsFrom{},
		PaymentFrom:    CoreSQLFrom.FieldsFrom{},
		TakeCreate:     CoreSQLFrom.FieldsFrom{},
		TakeChannel:    CoreSQLFrom.FieldsFrom{},
		TakeFrom:       CoreSQLFrom.FieldsFrom{},
		CreateInfo:     CoreSQLFrom.FieldsFrom{},
		TimeBetween:    CoreSQLTime.FieldsCoreTime{},
		IsHistory:      false,
		Search:         "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data, dataCount)
	}
}
