package OrgCoreCore

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	OrgTime "github.com/fotomxq/weeekj_core/v5/org/time"
	ToolsTest "github.com/fotomxq/weeekj_core/v5/tools/test"
	"testing"
	"time"
)

func TestInit6(t *testing.T) {
	TestInit(t)
	TestCreateUser(t)
}

func TestCreateOrg(t *testing.T) {
	var err error
	var errCode string
	//创建组织
	orgData, errCode, err = CreateOrg(&ArgsCreateOrg{
		UserID:     newUserInfo.ID,
		Key:        "",
		Name:       "组织名称测试",
		Des:        "组织描述测试",
		ParentID:   0,
		ParentFunc: []string{"x1", "x2"},
		OpenFunc:   []string{"all", "x3"},
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log(orgData)
	}
}

func TestCreateWorkTime(t *testing.T) {
	var err error
	workTimeData, err = OrgTime.Create(&OrgTime.ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    orgData.ID,
		Groups:   []int64{orgGroupData.ID},
		Binds:    []int64{},
		Name:     "上班测试A1",
		Configs: OrgTime.FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth()},
			Week:      []int{1, 2, 3, 4, 5, 6, 7, CoreFilter.GetNowTimeCarbon().DayOfWeek()},
			WorkTime: []OrgTime.FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,
					EndHour:     23,
					EndMinute:   59,
				},
			},
			AllowHoliday: false,
		},
	})
	if err != nil {
		t.Error(err, ", CoreFilter.GetNowTimeCarbon().DayOfWeek(): ", CoreFilter.GetNowTimeCarbon().DayOfWeek(), ", CoreFilter.GetNowTimeCarbon().DayOfMonth(): ", CoreFilter.GetNowTimeCarbon().DayOfMonth(), ", CoreFilter.GetNowTimeCarbon().WeekOfMonth(): ", CoreFilter.GetNowTimeCarbon().WeekOfMonth())
	} else {
		t.Log("new workTimeData: ", workTimeData)
	}
}

func TestGetOrgList(t *testing.T) {
	dataList, dataCount, err := GetOrgList(&ArgsGetOrgList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		UserID:     -1,
		ParentID:   0,
		ParentFunc: []string{"x1"},
		OpenFunc:   []string{},
		IsRemove:   false,
		Search:     "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
}

// 子组织获取的验证
func TestGetOrgListByChildOrg(t *testing.T) {
	//创建组织
	newOrgData, errCode, err := CreateOrg(&ArgsCreateOrg{
		UserID:     newUserInfo.ID,
		Name:       "组织名称测试child",
		Des:        "组织描述测试child",
		ParentID:   orgData.ID,
		ParentFunc: []string{},
		OpenFunc:   []string{"all"},
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log(newOrgData)
	}
	dataList, dataCount, err := GetOrgList(&ArgsGetOrgList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		UserID:   0,
		ParentID: orgData.ID,
		IsRemove: false,
		Search:   "",
	})
	ToolsTest.ReportDataList(t, err, dataList, dataCount)
	if err == nil {
		if len(dataList) != 1 {
			t.Error("无法获取子组织: ", dataList)
		} else {
			if dataList[0].ID == orgData.ID {
				t.Error("子组织不匹配")
			}
		}
	}
	if err := DeleteOrg(&ArgsDeleteOrg{
		ID: newOrgData.ID,
	}); err != nil {
		t.Error(err)
	}
}

func TestGetOrgSearch(t *testing.T) {
	dataList, err := GetOrgSearch(&ArgsGetOrgSearch{
		OrgID:  0,
		IDs:    nil,
		Search: "",
	})
	ToolsTest.ReportData(t, err, dataList)
	dataList, err = GetOrgSearch(&ArgsGetOrgSearch{
		OrgID:  0,
		IDs:    []int64{orgData.ID},
		Search: "",
	})
	ToolsTest.ReportData(t, err, dataList)
}

func TestGetOrg(t *testing.T) {
	data, err := GetOrg(&ArgsGetOrg{
		ID: orgData.ID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestGetOrgMore(t *testing.T) {
	data, err := GetOrgMore(&ArgsGetOrgMore{
		IDs: []int64{orgData.ID},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(data)
	}
}

func TestGetOrganizationalList(t *testing.T) {
	dataList, dataCount, err := GetOrgList(&ArgsGetOrgList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
		UserID:   0,
		ParentID: 0,
		IsRemove: false,
		Search:   "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log(dataList, dataCount)
	}
}

func TestUpdateOrg(t *testing.T) {
	if _, err := UpdateOrg(&ArgsUpdateOrg{
		ID:         orgData.ID,
		UserID:     orgData.UserID,
		Key:        orgData.Key,
		Name:       orgData.Name,
		Des:        orgData.Des,
		ParentID:   orgData.ParentID,
		ParentFunc: []string{"all"},
	}); err != nil {
		t.Error(err)
	}
}

func TestDeleteWorkTime(t *testing.T) {
	err := OrgTime.DeleteByID(&OrgTime.ArgsDeleteByID{
		ID:    workTimeData.ID,
		OrgID: orgData.ID,
	})
	if err != nil {
		t.Error(err, ", work time id: ", workTimeData.ID)
	}
}

func TestDeleteOrg(t *testing.T) {
	if err := DeleteOrg(&ArgsDeleteOrg{
		ID: orgData.ID,
	}); err != nil {
		t.Error(err)
	}
}

func TestClear5(t *testing.T) {
	TestClear(t)
}
