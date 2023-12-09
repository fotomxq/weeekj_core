package ServiceUserInfo

import (
	"errors"
	"fmt"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newTemplateData FieldsTemplate
)

func TestInitTemplate(t *testing.T) {
	TestInit(t)
}

func TestCreateTemplate(t *testing.T) {
	data, err := CreateTemplate(&ArgsCreateTemplate{
		OrgID:      1,
		Name:       "测试模版",
		Des:        "test des",
		CoverFiles: []int64{},
		SortID:     0,
		Tags:       []int64{},
		FileData:   "",
		Params:     nil,
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newTemplateData = data
	} else {
		err = errors.New(fmt.Sprint("create new template, ", err))
	}
}

func TestGetTemplateList(t *testing.T) {
	dataList, dataCount, err := GetTemplateList(&ArgsGetTemplateList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:    -1,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestUpdateTemplate(t *testing.T) {
	err := UpdateTemplate(&ArgsUpdateTemplate{
		ID:         newTemplateData.ID,
		OrgID:      newTemplateData.OrgID,
		Name:       newTemplateData.Name,
		Des:        newTemplateData.Des,
		CoverFiles: newTemplateData.CoverFiles,
		SortID:     newTemplateData.SortID,
		Tags:       newTemplateData.Tags,
		FileData:   newTemplateData.FileData,
		Params:     newTemplateData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestGetTemplateID(t *testing.T) {
	data, err := GetTemplateID(&ArgsGetTemplateID{
		ID:    newTemplateData.ID,
		OrgID: -1,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetTemplateMore(t *testing.T) {
	data, err := GetTemplateMore(&ArgsGetTemplateMore{
		IDs:        []int64{newTemplateData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestDeleteTemplate(t *testing.T) {
	err := DeleteTemplate(&ArgsDeleteTemplate{
		ID:    newTemplateData.ID,
		OrgID: newTemplateData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}
