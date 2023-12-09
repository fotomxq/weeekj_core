package FinanceAnalysis

import (
	"testing"

	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
)

func TestInit3(t *testing.T) {
	TestInit(t)
}

func TestAppendData2(t *testing.T) {
	TestAppendData(t)
}

func TestGetAnalysis(t *testing.T) {
	dataList, err := GetAnalysis(&ArgsGetAnalysis{
		TimeBetween: CoreSQLTime.FieldsCoreTime{
			MinTime: CoreFilter.GetNowTime().AddDate(0, 0, -1),
			MaxTime: CoreFilter.GetNowTime().AddDate(0, 0, 1),
		},
		TimeType:       "hour",
		PaymentCreate:  CoreSQLFrom.FieldsFrom{},
		PaymentChannel: CoreSQLFrom.FieldsFrom{},
		PaymentFrom:    CoreSQLFrom.FieldsFrom{},
		TakeCreate:     CoreSQLFrom.FieldsFrom{},
		TakeChannel:    CoreSQLFrom.FieldsFrom{},
		TakeFrom:       CoreSQLFrom.FieldsFrom{},
		Currency:       86,
		IsHistory:      false,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(dataList)
	}
}
