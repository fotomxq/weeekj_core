package ClassSort

import (
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
)

var (
	sort     Sort
	sortData FieldsSort
	bindID   int64 = 123
)

func TestInit(t *testing.T) {
	ToolsTest.Init(t)
}

func TestTag_Init(t *testing.T) {
	sort.Init("test_sort")
}
func TestTag_Create(t *testing.T) {
	data, err := sort.Create(&ArgsCreate{
		BindID:      bindID,
		Mark:        "",
		ParentID:    0,
		CoverFileID: 0,
		DesFiles:    []int64{},
		Name:        "测试分组A",
		Des:         "测试分组A描述",
		Params:      nil,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
		sortData = data
	}
	//上下级判断处理
	data2, err := sort.Create(&ArgsCreate{
		BindID:      bindID,
		Mark:        "",
		ParentID:    data.ID,
		CoverFileID: 0,
		DesFiles:    []int64{},
		Name:        "测试分组A",
		Des:         "测试分组A描述",
		Params:      nil,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data2)
	}
}

func TestTag_GetList(t *testing.T) {
	dataList, dataCount, err := sort.GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "_id",
			Desc: false,
		},
		BindID:   bindID,
		Mark:     "",
		ParentID: 0,
		Search:   "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(dataList, dataCount)
	}
}

func TestTag_UpdateByID(t *testing.T) {
	if err := sort.UpdateByID(&ArgsUpdateByID{
		ID:          sortData.ID,
		BindID:      bindID,
		Mark:        "",
		ParentID:    0,
		Sort:        0,
		CoverFileID: 0,
		DesFiles:    []int64{},
		Name:        "测试分组A edit",
		Des:         "测试分组A描述 edit",
		Params:      nil,
	}); err != nil {
		t.Error(err)
	} else {
		TestTag_GetList(t)
	}
}

func TestTag_DeleteByID(t *testing.T) {
	if err := sort.DeleteByID(&ArgsDeleteByID{
		ID:     sortData.ID,
		BindID: bindID,
	}); err != nil {
		t.Error(err)
	} else {
		TestTag_GetList(t)
	}
}

func TestClear(t *testing.T) {
}
