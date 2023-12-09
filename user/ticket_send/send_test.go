package UserTicketSend

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	TestOrg "gitee.com/weeekj/weeekj_core/v5/tools/test_org"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
	UserTicket "gitee.com/weeekj/weeekj_core/v5/user/ticket"
	"testing"
)

var (
	newSendData FieldsSend
	newUserID   int64
)

func TestInitSend(t *testing.T) {
	TestInit(t)
}

func TestCreateSend(t *testing.T) {
	newTicketConfig, err := UserTicket.CreateConfig(&UserTicket.ArgsCreateConfig{
		DefaultExpireTime: 60,
		OrgID:             TestOrg.OrgData.ID,
		Mark:              "test",
		Title:             "测试",
		Des:               "测试描述",
		CoverFileID:       0,
		DesFiles:          []int64{},
		ExemptionPrice:    0,
		ExemptionDiscount: 0,
		ExemptionMinPrice: 0,
		UseOrder:          false,
		LimitTimeType:     0,
		LimitCount:        0,
		StyleID:           0,
		Params:            nil,
	})
	if err != nil {
		t.Error(err)
		return
	}
	data, err := CreateSend(&ArgsCreateSend{
		OrgID:               TestOrg.OrgData.ID,
		NeedUserSubConfigID: 0,
		NeedAuto:            true,
		ConfigID:            newTicketConfig.ID,
		PerCount:            1,
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newSendData = data
	}
	userInfo, _, err := UserCore.CreateUser(&UserCore.ArgsCreateUser{
		OrgID:                TestOrg.OrgData.ID,
		Name:                 "测试账户2",
		Password:             "",
		NationCode:           "",
		Phone:                "",
		AllowSkipPhoneVerify: false,
		AllowSkipWaitEmail:   false,
		Email:                "",
		Username:             "",
		Avatar:               0,
		Status:               2,
		Parents:              nil,
		Groups:               nil,
		Infos:                nil,
		Logins:               nil,
		SortID:               0,
		Tags:                 nil,
	})
	if err == nil {
		newUserID = userInfo.ID
	} else {
		t.Error(err)
	}
}

func TestGetSendList(t *testing.T) {
	dataList, dataCount, err := GetSendList(&ArgsGetSendList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:    -1,
		ConfigID: -1,
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil {
		if len(dataList) > 0 {
			if newSendData.ID != dataList[0].ID {
				t.Error("no data")
			}
		}
	}
}

func TestRun(t *testing.T) {
	runSend()
	count, _ := UserTicket.GetTicketCount(&UserTicket.ArgsGetTicketCount{
		ConfigID: newSendData.ConfigID,
		UserID:   newUserID,
	})
	if count < 1 {
		t.Error("less 1")
	}
}

func TestDeleteSend(t *testing.T) {
	err := DeleteSend(&ArgsDeleteSend{
		ID:    newSendData.ID,
		OrgID: -1,
	})
	ToolsTest.ReportError(t, err)
}
