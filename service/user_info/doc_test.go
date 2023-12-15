package ServiceUserInfo

import (
	"errors"
	"fmt"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newDocData FieldsDoc
)

func TestInitDoc(t *testing.T) {
	TestInitInfo(t)
}

func TestCreateDoc(t *testing.T) {
	TestCreateInfo(t)
	data, err := CreateDoc(&ArgsCreateDoc{
		OrgID:      1,
		InfoID:     newInfoData.ID,
		Title:      "测试标题",
		TemplateID: newTemplateData.ID,
		SortID:     0,
		Tags:       []int64{},
		FileData:   "",
		Params:     CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, data)
	if err == nil {
		newDocData = data
	} else {
		err = errors.New(fmt.Sprint("create Doc, ", err))
	}
}

func TestGetDocList(t *testing.T) {
	dataList, dataCount, err := GetDocList(&ArgsGetDocList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:      -1,
		InfoID:     -1,
		TemplateID: -1,
		SortID:     0,
		Tags:       []int64{},
		IsRemove:   false,
		Search:     "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetDocID(t *testing.T) {
	data, err := GetDocID(&ArgsGetDocID{
		ID:    newDocData.ID,
		OrgID: newDocData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetDocMore(t *testing.T) {
	data, err := GetDocMore(&ArgsGetDocMore{
		IDs:        []int64{newDocData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetDocMoreNames(t *testing.T) {
	data, err := GetDocMoreNames(&ArgsGetDocMore{
		IDs:        []int64{newDocData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetOrgDocMore(t *testing.T) {
	data, err := GetOrgDocMore(&ArgsGetOrgDocMore{
		IDs:        []int64{newDocData.ID},
		HaveRemove: false,
		OrgID:      newDocData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetOrgDocMoreNames(t *testing.T) {
	data, err := GetOrgDocMoreNames(&ArgsGetOrgDocMore{
		IDs:        []int64{newDocData.ID},
		HaveRemove: false,
		OrgID:      newDocData.OrgID,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateDoc(t *testing.T) {
	err := UpdateDoc(&ArgsUpdateDoc{
		ID:       newDocData.ID,
		OrgID:    newDocData.OrgID,
		Title:    newDocData.Title,
		SortID:   newDocData.SortID,
		Tags:     newDocData.Tags,
		FileData: "",
		Params:   newDocData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteDoc(t *testing.T) {
	err := DeleteDoc(&ArgsDeleteDoc{
		ID:    newDocData.ID,
		OrgID: newDocData.OrgID,
	})
	ToolsTest.ReportError(t, err)
	TestDeleteInfo(t)
}
