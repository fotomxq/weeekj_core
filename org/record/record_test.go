package OrgRecord

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
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

func TestCreate(t *testing.T) {
	if err := Create(&ArgsCreate{
		OrgID:       123,
		FromMark:    "",
		FromID:      0,
		BindID:      1234,
		ContentMark: "create",
		Content:     "创建了人",
		ChangeData:  nil,
	}); err != nil {
		t.Error(err)
	}
	if err := Create(&ArgsCreate{
		OrgID:       123,
		FromMark:    "",
		FromID:      0,
		BindID:      1234,
		ContentMark: "create2",
		Content:     "创建了人2",
		ChangeData:  nil,
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
		OrgID:        123,
		FromMark:     "",
		FromID:       0,
		BindID:       0,
		ContentMarks: []string{},
		IsHistory:    false,
		Search:       "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data, count)
	}
	data2, count, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:        123,
		FromMark:     "",
		FromID:       0,
		BindID:       0,
		ContentMarks: []string{},
		IsHistory:    false,
		Search:       "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data2, count)
		//检查data1和data2排序问题
		if data[0].ID == data2[0].ID {
			t.Error("排序存在错误")
		} else {
			t.Log("data1:0.id: ", data[0].ID, "; data2:0.id: ", data2[0].ID)
		}
	}
}
