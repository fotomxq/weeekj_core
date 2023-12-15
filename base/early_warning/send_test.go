package BaseEarlyWarning

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	"testing"
)

func TestInit5(t *testing.T) {
	TestInit(t)
	TestCreateTemplate(t)
	TestCreateTo(t)
	TestSetBind(t)
}

// 发送测试
func TestSendMod(t *testing.T) {
	if err := SendMod(&ArgsSendMod{
		Mark: testTemplateMark, Contents: map[string]string{"[1]": "变量1内容", "[2]": "变量2内容"},
	}); err != nil {
		t.Error(err)
	}
}

func TestGetMessage(t *testing.T) {
	dataList, dataCount, err := GetMessage(&ArgsGetMessage{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		UserID:     0,
		ToID:       0,
		TemplateID: 0,
		NeedIsRead: false,
		IsRead:     false,
		Search:     "",
	})
	if err != nil {
		t.Error(err)
	} else {
		if len(dataList) > 0 {
			sendData = dataList[0]
		}
		t.Log(dataList, dataCount)
	}
}

func TestUpdateSendIsReadByID(t *testing.T) {
	if err := UpdateSendIsRead(&ArgsUpdateSendIsRead{
		ID: sendData.ID,
	}); err != nil {
		t.Error(err)
	}
}

func TestUpdateSendIsReadByMark(t *testing.T) {
	TestDeleteSendByID(t)
	TestSendMod(t)
	TestGetMessage(t)
	if err := UpdateSendIsReadByMark(&ArgsUpdateSendIsReadByMark{
		Mark: testTemplateMark,
	}); err != nil {
		t.Error(err)
	}
}

func TestDeleteSendByID(t *testing.T) {
	if err := DeleteSendByID(&ArgsDeleteSendByID{
		ID: sendData.ID,
	}); err != nil {
		t.Error(err)
	}
}

func TestDeleteSendByMark(t *testing.T) {
	TestSendMod(t)
	if err := DeleteSendByMark(&ArgsDeleteSendByMark{
		Mark: templateData.Mark,
	}); err != nil {
		t.Error(err)
	}
}
func TestClear2(t *testing.T) {
	TestDeleteSendByID(t)
	TestSetUnBind(t)
	TestClear(t)
}
