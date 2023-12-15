package OrgTip

import (
	"testing"
	"time"

	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
)

var (
	createID int64 = 123
)

func TestInitTip(t *testing.T) {
	TestInit(t)
}

// 创建推送
func TestCreate(t *testing.T) {
	if err := Create(&ArgsCreate{
		OrgID: TestOrg.OrgData.ID,
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "mod",
			ID:     createID,
			Mark:   "",
			Name:   "",
		},
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     TestOrg.UserInfo.ID,
			Mark:   "",
			Name:   TestOrg.UserInfo.Name,
		},
		TipAt:       CoreFilter.GetNowTime(),
		Title:       "测试推送消息",
		Content:     "测试推送消息内容",
		Files:       []int64{},
		NeedSMS:     false,
		SMSConfigID: 0,
		SMSParams:   nil,
		Params:      nil,
	}); err != nil {
		t.Error(err)
	}
}

// 删除推送
func TestDeleteByCreateFrom(t *testing.T) {
	if err := DeleteByCreateFrom(&ArgsDeleteByCreateFrom{
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "mod",
			ID:     createID,
			Mark:   "",
			Name:   "",
		},
	}); err != nil {
		t.Error(err)
	}
}

// 进行推送
func TestRun(t *testing.T) {
	TestCreate(t)
	go runSend()
	time.Sleep(time.Second * 3)
}

func TestClearTip(t *testing.T) {
	TestClear(t)
}
