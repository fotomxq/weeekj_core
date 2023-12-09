package FinanceSafe

import (
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
	"time"
)

var (
	isInit = false
)

//注意，使用本模块测试前，请先执行FinanceLog的单元测试模块，创建日志数据
// 创建之后的数据将作为测试依据，如果条件允许或需要，请创建明显存在某些异常条件的数据

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
	}
	isInit = true
}

func TestRun(t *testing.T) {
	go Run()
	time.Sleep(time.Second * 3)
}

func TestGetList(t *testing.T) {
	data, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "_id",
			Desc: false,
		},
		CreateInfo:    CoreSQLFrom.FieldsFrom{},
		PaymentCreate: CoreSQLFrom.FieldsFrom{},
		PaymentFrom:   CoreSQLFrom.FieldsFrom{},
		TakeCreate:    CoreSQLFrom.FieldsFrom{},
		TakeFrom:      CoreSQLFrom.FieldsFrom{},
		NeedAllowEW:   false,
		AllowEW:       false,
		Code:          "",
		AllowOpen:     false,
		Search:        "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data, dataCount)
	}
}

func TestUpdateDone(t *testing.T) {
	data, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "_id",
			Desc: false,
		},
		CreateInfo:    CoreSQLFrom.FieldsFrom{},
		PaymentCreate: CoreSQLFrom.FieldsFrom{},
		PaymentFrom:   CoreSQLFrom.FieldsFrom{},
		TakeCreate:    CoreSQLFrom.FieldsFrom{},
		TakeFrom:      CoreSQLFrom.FieldsFrom{},
		NeedAllowEW:   false,
		AllowEW:       false,
		Code:          "",
		AllowOpen:     false,
		Search:        "",
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(data, dataCount)
	if len(data) < 1 {
		t.Log("data is empty")
		return
	}
	if err := UpdateDone(&ArgsUpdateDone{
		ID:         data[0].ID,
		TakeCreate: CoreSQLFrom.FieldsFrom{},
		TakeFrom:   CoreSQLFrom.FieldsFrom{},
	}); err != nil {
		t.Error(err)
	}
}
