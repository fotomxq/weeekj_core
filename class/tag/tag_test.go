package ClassTag

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	tag     Tag
	tagData FieldsTag
	bindID  int64 = 123
)

func TestInit(t *testing.T) {
	ToolsTest.Init(t)
}

func TestTag_Init(t *testing.T) {
	tag.Init("test_tag")
}

func TestTag_Create(t *testing.T) {
	data, err := tag.Create(&ArgsCreate{
		BindID: bindID,
		Name:   "测试标签A",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
		tagData = data
	}
}

func TestTag_GetList(t *testing.T) {
	dataList, dataCount, err := tag.GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "_id",
			Desc: false,
		},
		BindID: bindID,
		Search: "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(dataList, dataCount)
	}
}

func TestTag_UpdateByID(t *testing.T) {
	if err := tag.UpdateByID(&ArgsUpdateByID{
		ID:     tagData.ID,
		BindID: bindID,
		Name:   "测试标签B",
	}); err != nil {
		t.Error(err)
	}
}

func TestTag_DeleteByID(t *testing.T) {
	if err := tag.DeleteByID(&ArgsDeleteByID{
		ID:     tagData.ID,
		BindID: bindID,
	}); err != nil {
		t.Error(err)
	}
}

func TestClear(t *testing.T) {

}
