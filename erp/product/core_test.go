package ERPProduct

import (
	BaseBPM "github.com/fotomxq/weeekj_core/v5/base/bpm"
	ClassSort "github.com/fotomxq/weeekj_core/v5/class/sort"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	ServiceCompany "github.com/fotomxq/weeekj_core/v5/service/company"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	TestOrg "github.com/fotomxq/weeekj_core/v5/tools/test_org"
	"testing"
)

var (
	isInit      = false
	newSortData ClassSort.FieldsSort
)

func TestInit(t *testing.T) {
	if !isInit {
		ToolsTest.Init(t)
		TestOrg.LocalInit()
		ServiceCompany.Init()
		BaseBPM.Init()
	}
	isInit = true
	Init()
	if TestOrg.OrgData.ID < 1 {
		TestOrg.LocalCreateOrg(t)
	}
}

// TestSort 测试分类
func TestSortCreate(t *testing.T) {
	var err error
	newSortData, err = Sort.Create(&ClassSort.ArgsCreate{
		BindID:      TestOrg.OrgData.ID,
		Mark:        CoreFilter.GetRandStr4(10),
		ParentID:    0,
		CoverFileID: 0,
		DesFiles:    []int64{},
		Name:        "测试分类",
		Des:         "测试分类描述",
		Params:      CoreSQLConfig.FieldsConfigsType{},
	})
	ToolsTest.ReportData(t, err, newSortData)
}

func TestSortDelete(t *testing.T) {
	err := Sort.DeleteByID(&ClassSort.ArgsDeleteByID{
		ID:     newSortData.ID,
		BindID: TestOrg.OrgData.ID,
	})
	ToolsTest.ReportError(t, err)
}

func TestClear(t *testing.T) {
	TestOrg.LocalClear(t)
}
