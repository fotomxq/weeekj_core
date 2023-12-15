package BaseEmail

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	serverID    int64
	sendToEmail = "fotomxq@qq.com"

	newData FieldsEmailType

	isRun = false
)

func TestInit(t *testing.T) {
	if isRun {
		return
	}
	isRun = true
	ToolsTest.Init(t)
	go Run()
	//创建新的
	TestCreateServer(t)
}

func TestSend(t *testing.T) {
	var err error
	//发送消息给某个邮箱
	newData, err = Send(&ArgsSend{
		ServerID: serverID,
		CreateInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     123,
			Mark:   "",
			Name:   "",
		},
		SendAt:  CoreFilter.GetNowTime(),
		ToEmail: sendToEmail,
		Title:   "test",
		Content: "<html><body>test content.</body></html>",
		IsHtml:  true,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("send new mail, ", newData)
	}
	//go Run()
	//time.Sleep(time.Second * 3)
}

// 封闭测试内部的发送系统
func TestSend2(t *testing.T) {
	t.Skip()
	//尝试发送数据
	if err := sendSSLMail(serverGlobData, newData); err != nil {
		t.Error(err)
	}
}

func TestGetEmailByID(t *testing.T) {
	getData, err := GetEmailByID(&ArgsGetEmailByID{
		ID: newData.ID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("get data, ", getData)
	}
}

func TestGetEmailList(t *testing.T) {
	var err error
	//获取列表
	getListData, getListCount, err := GetEmailList(&ArgsGetEmailList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		IsSuccess:     true,
		CreateInfo:    CoreSQLFrom.FieldsFrom{System: "user", ID: 1},
		ToEmailSearch: "",
		Search:        "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("get list, ", getListCount, ", ", getListData)
	}
}

func TestDeleteEmailByID(t *testing.T) {
	if err := DeleteEmailByID(&ArgsDeleteEmailByID{
		ID: newData.ID,
	}); err != nil {
		t.Error(err)
	}
}

func TestClear(t *testing.T) {
	//清理
	if _, err := CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_email", "", nil); err != nil {
		t.Error(err)
	}
	TestClearServer(t)
}
