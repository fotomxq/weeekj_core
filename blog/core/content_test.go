package BlogCore

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
)

var (
	newContentData FieldsContent
)

func TestInitContent(t *testing.T) {
	TestInit(t)
}

func TestCreateContent(t *testing.T) {
	newKey, err := CoreFilter.GetRandStr3(10)
	if err != nil {
		t.Error(err)
		return
	}
	newContentData, err = CreateContent(&ArgsCreateContent{
		OrgID:       1,
		UserID:      0,
		Key:         newKey,
		IsTop:       false,
		SortID:      0,
		Tags:        []int64{1, 3, 4, 5},
		Title:       "测试标题",
		TitleDes:    "",
		CoverFileID: 0,
		Des:         "测试内容姑娘",
		Params:      CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, newContentData)
	newKey2, err := CoreFilter.GetRandStr3(10)
	if err != nil {
		t.Error(err)
		return
	}
	_, err = CreateContent(&ArgsCreateContent{
		OrgID:       1,
		UserID:      0,
		Key:         newKey2,
		IsTop:       false,
		SortID:      0,
		Tags:        []int64{4, 5},
		Title:       "测试标题",
		TitleDes:    "",
		CoverFileID: 0,
		Des:         "测试内容姑娘",
		Params:      CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportError(t, err)
}

func TestGetContentByID(t *testing.T) {
	data, err := GetContentByID(&ArgsGetContentByID{
		ID:         newContentData.ID,
		OrgID:      newContentData.OrgID,
		UserID:     0,
		IsPublish:  false,
		ReadUserID: 0,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetContentByKey(t *testing.T) {
	data, err := GetContentByKey(&ArgsGetContentByKey{
		Key:        newContentData.Key,
		OrgID:      newContentData.OrgID,
		UserID:     0,
		IsPublish:  false,
		ReadUserID: 0,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetContentList(t *testing.T) {
	dataList, dataCount, err := GetContentList(&ArgsGetContentList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:         -1,
		UserID:        -1,
		SortID:        -1,
		Tags:          []int64{},
		TagsOr:        false,
		NeedIsPublish: false,
		IsPublish:     false,
		ParentID:      0,
		NeedIsTop:     false,
		IsTop:         false,
		IsRemove:      false,
		Search:        "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	dataList, dataCount, err = GetContentList(&ArgsGetContentList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:         -1,
		UserID:        -1,
		SortID:        -1,
		Tags:          []int64{},
		TagsOr:        false,
		NeedIsPublish: false,
		IsPublish:     false,
		ParentID:      0,
		NeedIsTop:     false,
		IsTop:         false,
		IsRemove:      false,
		Search:        "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	dataList, dataCount, err = GetContentList(&ArgsGetContentList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:         -1,
		UserID:        -1,
		SortID:        -1,
		Tags:          []int64{1, 3},
		TagsOr:        false,
		NeedIsPublish: false,
		IsPublish:     false,
		ParentID:      0,
		NeedIsTop:     false,
		IsTop:         false,
		IsRemove:      false,
		Search:        "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	for _, v := range dataList {
		isFind := false
		isFind2 := false
		for _, v2 := range v.Tags {
			if v2 == 1 {
				isFind = true
			}
			if v2 == 3 {
				isFind2 = true
			}
		}
		if !isFind || !isFind2 {
			t.Error("tag error, no 1 or 3, tags: ", v.Tags)
		}
	}
	dataList, dataCount, err = GetContentList(&ArgsGetContentList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		OrgID:         -1,
		UserID:        -1,
		SortID:        -1,
		Tags:          []int64{1, 3, 15},
		TagsOr:        true,
		NeedIsPublish: false,
		IsPublish:     false,
		ParentID:      0,
		NeedIsTop:     false,
		IsTop:         false,
		IsRemove:      false,
		Search:        "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	for _, v := range dataList {
		isFind := false
		for _, v2 := range v.Tags {
			if v2 == 1 || v2 == 3 {
				isFind = true
				break
			}
		}
		if !isFind {
			t.Error("tag error, no 1 or 3, tags: ", v.Tags)
		}
	}
}

func TestGetContentMore(t *testing.T) {
	data, err := GetContentMore(&ArgsGetContentMore{
		IDs:        []int64{newContentData.ID},
		OrgID:      newContentData.OrgID,
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestGetContentMoreMap(t *testing.T) {
	data, err := GetContentMoreMap(&ArgsGetContentMore{
		IDs:        []int64{newContentData.ID},
		OrgID:      newContentData.OrgID,
		HaveRemove: false,
	})
	ToolsTest.ReportData(t, err, data)
}

func TestUpdateContent(t *testing.T) {
	var err error
	newContentData, err = UpdateContent(&ArgsUpdateContent{
		ID:          newContentData.ID,
		OrgID:       newContentData.OrgID,
		UserID:      newContentData.UserID,
		Key:         newContentData.Key,
		IsTop:       newContentData.IsTop,
		SortID:      newContentData.SortID,
		Tags:        newContentData.Tags,
		Title:       newContentData.Title,
		TitleDes:    newContentData.TitleDes,
		CoverFileID: newContentData.CoverFileID,
		Des:         newContentData.Des,
		Params:      newContentData.Params,
	})
	ToolsTest.ReportData(t, err, newContentData)
}

func TestUpdatePublish(t *testing.T) {
	err := UpdatePublish(&ArgsUpdatePublish{
		ID:     newContentData.ID,
		OrgID:  newContentData.OrgID,
		UserID: newContentData.UserID,
	})
	ToolsTest.ReportError(t, err)
}

func TestDeleteContent(t *testing.T) {
	err := DeleteContent(&ArgsDeleteContent{
		ID:     newContentData.ID,
		OrgID:  newContentData.OrgID,
		UserID: newContentData.UserID,
	})
	ToolsTest.ReportError(t, err)
}

func TestReturnContent(t *testing.T) {
	err := ReturnContent(&ArgsReturnContent{
		ID:     newContentData.ID,
		OrgID:  newContentData.OrgID,
		UserID: newContentData.UserID,
	})
	ToolsTest.ReportError(t, err)
}
