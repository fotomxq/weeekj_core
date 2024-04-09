package ClassComment

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	comment     Comment
	commentData FieldsComment
	bindID      int64 = 123
)

func TestInit(t *testing.T) {
	ToolsTest.Init(t)
}

func TestComment_Init(t *testing.T) {
	//comment.Init("test_comment", "test_comment")
	comment.TableName = "test_comment"
}

func TestComment_Create(t *testing.T) {
	var err error
	commentData, err = comment.Create(&ArgsCreate{
		ParentID:  0,
		OrgID:     123,
		UserID:    234,
		BindID:    bindID,
		LevelType: 0,
		Level:     1,
		Title:     "测试标题",
		Des:       "测试描述",
		DesFiles:  []int64{1},
		Params:    []CoreSQLConfig.FieldsConfigType{},
	})
	ToolsTest.ReportData(t, err, commentData)
}

func TestComment_GetList(t *testing.T) {
	dataList, dataCount, err := comment.GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		CommentID: 0,
		ParentID:  0,
		OrgID:     0,
		UserID:    0,
		BindID:    0,
		LevelType: 0,
		LevelMin:  0,
		LevelMax:  0,
		IsRemove:  false,
		Search:    "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

func TestComment_Update(t *testing.T) {
	err := comment.Update(&ArgsUpdate{
		ID:        commentData.ID,
		OrgID:     commentData.OrgID,
		UserID:    commentData.UserID,
		LevelType: commentData.LevelType,
		Level:     commentData.Level,
		Title:     commentData.Title,
		Des:       commentData.Des,
		DesFiles:  commentData.DesFiles,
		Params:    commentData.Params,
	})
	ToolsTest.ReportError(t, err)
}

func TestComment_DeleteByID(t *testing.T) {
	TestComment_Create(t)
	err := comment.DeleteByID(&ArgsDeleteByID{
		ID:     commentData.ID,
		OrgID:  0,
		UserID: 0,
	})
	ToolsTest.ReportError(t, err)
}

func TestComment_DeleteByBind(t *testing.T) {
	TestComment_Create(t)
	err := comment.DeleteByBind(&ArgsDeleteByBind{
		BindID: commentData.BindID,
		OrgID:  0,
	})
	ToolsTest.ReportError(t, err)
}

func TestComment_DeleteByOrg(t *testing.T) {
	TestComment_Create(t)
	err := comment.DeleteByOrg(&ArgsDeleteByOrg{
		OrgID: commentData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}

func TestComment_DeleteByUser(t *testing.T) {
	TestComment_Create(t)
	err := comment.DeleteByUser(&ArgsDeleteByUser{
		UserID: commentData.UserID,
		OrgID:  commentData.OrgID,
	})
	ToolsTest.ReportError(t, err)
}
