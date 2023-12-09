package FinanceAnalysis

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	"testing"
)

func TestInit2(t *testing.T) {
	TestInit(t)
}

func TestAppendData(t *testing.T) {
	var err error
	p := 1
	for {
		if p > 10 {
			break
		}
		err = AppendData(&ArgsAppendData{
			PaymentCreate: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     1,
				Mark:   "",
				Name:   "",
			},
			PaymentChannel: CoreSQLFrom.FieldsFrom{
				System: "weixin",
				ID:     0,
				Mark:   "wxx",
				Name:   "",
			},
			PaymentFrom: CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     1,
				Mark:   "",
				Name:   "",
			},
			TakeCreate: CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     1,
				Mark:   "",
				Name:   "",
			},
			TakeChannel: CoreSQLFrom.FieldsFrom{
				System: "",
				ID:     0,
				Mark:   "",
				Name:   "",
			},
			TakeFrom: CoreSQLFrom.FieldsFrom{
				System: "org",
				ID:     1,
				Mark:   "",
				Name:   "",
			},
			Currency: 86,
			Price:    int64(CoreFilter.GetRandNumber(1, 99999)),
		})
		if err != nil {
			t.Error(err)
		} else {
			t.Log("append new data, count: ", p)
		}
		p += 1
	}
}

func TestClear(t *testing.T) {
}
