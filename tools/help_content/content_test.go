package ToolsHelpContent

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newContentData FieldsContent
)

func TestInitContent(t *testing.T) {
	TestInit(t)
}

func TestCreateContent(t *testing.T) {
	var err error
	newContentData, err = CreateContent(&ArgsCreateContent{
		Mark:        "mark01",
		IsPublic:    true,
		SortID:      0,
		Tags:        []int64{},
		Title:       "标题内容",
		CoverFileID: 0,
		Des:         "",
		BindIDs:     []int64{},
		BindMarks:   []string{},
		Params:      []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, newContentData)
}

func TestGetContentList(t *testing.T) {
	dataList, dataCount, err := GetContentList(&ArgsGetContentList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		Mark:     "",
		SortID:   0,
		Tags:     []int64{},
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestGetContentByID(t *testing.T) {
	data, err := GetContentByID(&ArgsGetContentByID{
		ID:         newContentData.ID,
		NeedPublic: false,
		IsPublic:   false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetContentByMark(t *testing.T) {
	data, err := GetContentByMark(&ArgsGetContentByMark{
		Mark:       newContentData.Mark,
		NeedPublic: false,
		IsPublic:   true,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetContentMore(t *testing.T) {
	data, err := GetContentMore(&ArgsGetContentMore{
		IDs:        []int64{newContentData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetContentMoreMap(t *testing.T) {
	data, err := GetContentMoreMap(&ArgsGetContentMore{
		IDs:        []int64{newContentData.ID},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetContentMoreByMark(t *testing.T) {
	data, err := GetContentMoreByMark(&ArgsGetContentMoreByMark{
		Marks:      []string{newContentData.Mark},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetContentMoreByMarkMap(t *testing.T) {
	data, err := GetContentMoreByMarkMap(&ArgsGetContentMoreByMark{
		Marks:      []string{newContentData.Mark},
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateContent(t *testing.T) {
	err := UpdateContent(&ArgsUpdateContent{
		ID:          newContentData.ID,
		Mark:        newContentData.Mark,
		IsPublic:    newContentData.IsPublic,
		SortID:      newContentData.SortID,
		Tags:        newContentData.Tags,
		Title:       newContentData.Title,
		CoverFileID: newContentData.CoverFileID,
		Des:         newContentData.Des,
		BindIDs:     newContentData.BindIDs,
		BindMarks:   newContentData.BindMarks,
		Params:      newContentData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteContent(t *testing.T) {
	err := DeleteContent(&ArgsDeleteContent{
		ID: newContentData.ID,
	})
	ToolsTest.ReportError(t, err)
}
