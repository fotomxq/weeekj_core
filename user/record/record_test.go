package UserRecordCore

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
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

func TestCreate(t *testing.T) {
	if err := Create(&ArgsCreate{
		OrgID:       0,
		UserID:      123,
		UserName:    "user123name",
		ContentMark: "create",
		Content:     "创建了人群",
	}); err != nil {
		t.Error(err)
	}
}

func TestGetList(t *testing.T) {
	data, count, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:       0,
		UserID:      0,
		ContentMark: "",
		IsHistory:   false,
		Search:      "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data, count)
	}
}
