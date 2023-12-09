package BaseEmail

import (
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
	"time"
)

var (
	newTemplateData FieldsTemplate
)

func TestInitTemplate(t *testing.T) {
	TestInitServer(t)
	TestCreateServer(t)
}

func TestCreateTemplate(t *testing.T) {
	err := CreateTemplate(&ArgsCreateTemplate{
		OrgID:     0,
		ServerIDs: []int64{serverGlobData.ID},
		Title:     "测试模版",
		Content:   "模版内容${userID}",
		Params:    nil,
	})
	ToolsTest.ReportError(t, err)
}

func TestGetTemplateList(t *testing.T) {
	dataList, dataCount, err := GetTemplateList(&ArgsGetTemplateList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:    -1,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil {
		newTemplateData = dataList[0]
	}
}

func TestGetTemplate(t *testing.T) {
	data, err := GetTemplate(&ArgsGetTemplate{
		ID:    newTemplateData.ID,
		OrgID: newTemplateData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateTemplate(t *testing.T) {
	err := UpdateTemplate(&ArgsUpdateTemplate{
		ID:        newTemplateData.ID,
		OrgID:     newTemplateData.OrgID,
		ServerIDs: newTemplateData.ServerIDs,
		Title:     newTemplateData.Title,
		Content:   newTemplateData.Content,
		Params:    newTemplateData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestSendTemplate(t *testing.T) {
	err := SendTemplate(&ArgsSendTemplate{
		ID:    newTemplateData.ID,
		OrgID: newTemplateData.OrgID,
		ReplaceData: CoreSQLConfig.FieldsConfigsType{
			{
				Mark: "${userID}",
				Val:  "123123",
			},
		},
		CreateInfo:  CoreSQLFrom.FieldsFrom{},
		SendAt:      time.Time{},
		ToEmailList: []string{"fotomxq@qq.com"},
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteTemplate(t *testing.T) {
	err := DeleteTemplate(&ArgsDeleteTemplate{
		ID:    newTemplateData.ID,
		OrgID: newTemplateData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClearTemplate(t *testing.T) {
	TestClear(t)
}
