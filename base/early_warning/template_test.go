package BaseEarlyWarning

import (
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	"testing"
)

func TestInit2(t *testing.T) {
	TestInit(t)
}

func TestCreateTemplate(t *testing.T) {
	var err error
	templateData, err = CreateTemplate(&ArgsCreateTemplate{
		Mark: testTemplateMark, Name: "测试模版", DefaultExpireTime: "60s", Title: "测试模版", Content: "测试模版内容...[1],[2],...", BindData: []string{"[1]", "[2]"},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGetTemplateByMark(t *testing.T) {
	getData, err := GetTemplateByMark(&ArgsGetTemplateByMark{
		Mark: testTemplateMark,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(getData)
	}
}

func TestGetTemplateList(t *testing.T) {
	getDataList, count, err := GetTemplateList(&ArgsGetTemplateList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("count:", count, ", list: ", getDataList)
	}
}

func TestGetTemplateByID(t *testing.T) {
	getData, err := GetTemplateByID(&ArgsGetTemplateByID{
		ID: templateData.ID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(getData)
	}
}

func TestUpdateTemplate(t *testing.T) {
	if err := UpdateTemplate(&ArgsUpdateTemplate{
		ID: templateData.ID, Mark: testTemplateMark, Name: "测试修改1", DefaultExpireTime: "30s", Title: "测试修改标题1", Content: "测试修改内容...[1]，[2]，。。。", BindData: []string{"[1]", "[2]"},
	}); err == nil {
		getData, err := GetTemplateByID(&ArgsGetTemplateByID{
			ID: templateData.ID,
		})
		if err != nil {
			t.Error(err)
		} else {
			t.Log(getData)
		}
	} else {
		t.Error(err)
	}
}

func TestDeleteTemplate(t *testing.T) {
	if err := DeleteTemplate(&ArgsDeleteTemplate{
		ID: templateData.ID,
	}); err != nil {
		t.Error(err)
	}
}
