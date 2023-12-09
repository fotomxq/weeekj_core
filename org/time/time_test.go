package OrgTime

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"github.com/golang-module/carbon"
	"testing"
	"time"
)

var (
	isInit             = false
	defaultOrgID int64 = 1
)

func TestTimeInit(t *testing.T) {
	if isInit {
		return
	}
	isInit = true
	ToolsTest.Init(t)
}

func TestCreate(t *testing.T) {
	newData, err := Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试A1",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{1, 2, 3, 4, 5, 6, 7},
			WorkTime: []FieldsWorkTimeTime{
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
		t.Error(err)
	} else {
		t.Log("new data: ", newData)
		err = DeleteByID(&ArgsDeleteByID{
			ID:    newData.ID,
			OrgID: newData.OrgID,
		})
		if err != nil {
			t.Error(err, ", newData.ID: ", newData.ID)
		}
	}
	newData, err = Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试A2",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{1, 2, 3, 4, 5, 6, 7}, WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,
					EndHour:     23,
					EndMinute:   59,
				},
			},
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new data: ", newData)
		err = DeleteByID(&ArgsDeleteByID{
			ID:    newData.ID,
			OrgID: newData.OrgID,
		})
		if err != nil {
			t.Error(err, ", newData.ID: ", newData.ID)
		}
	}
	newData, err = Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试A3",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{1, 2, 3, 6, 7, 8},
			WorkTime: []FieldsWorkTimeTime{
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
	if err == nil {
		t.Error("week check is error")
		err = DeleteByID(&ArgsDeleteByID{
			ID:    newData.ID,
			OrgID: newData.OrgID,
		})
		if err != nil {
			t.Error(err, ", newData.ID: ", newData.ID)
		}
	} else {
		t.Log(newData)
	}
	newData, err = Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试A4",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{},
			WorkTime: []FieldsWorkTimeTime{
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
		t.Error(err)
	} else {
		t.Log(newData)
		err = DeleteByID(&ArgsDeleteByID{
			ID:    newData.ID,
			OrgID: newData.OrgID,
		})
		if err != nil {
			t.Error(err, ", newData.ID: ", newData.ID)
		}
	}
	newData, err = Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试A5",
		Configs: FieldsConfigs{
			Month:     []int{13},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   23,
					EndMinute: 59,
				},
			},
			AllowHoliday: false,
		},
	})
	if err == nil {
		t.Error("month check is error")
		err = DeleteByID(&ArgsDeleteByID{
			ID:    newData.ID,
			OrgID: newData.OrgID,
		})
		if err != nil {
			t.Error(err, ", newData.ID: ", newData.ID)
		}
	} else {
		t.Log(newData)
	}
}

func TestGetAllByUserID(t *testing.T) {
	newData, err := Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试B",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   23,
					EndMinute: 59,
				},
			},
			AllowHoliday: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new data: ", newData)
	}
	//尝试获取该数据
	getData, dataCount, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:  0,
		Search: "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("get all data by user: ", dataCount, ", ", getData)
	}
	getData, dataCount, err = GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: true,
		},
		OrgID:  defaultOrgID,
		Search: "",
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("get all data by user: ", dataCount, ", ", getData)
	}
	getData, dataCount, err = GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "_id",
			Desc: false,
		},
		OrgID:  92387128381,
		Search: "",
	})
	if err == nil {
		if dataCount > 0 {
			t.Error("get not exist data by user, count: ", dataCount, ", data: ", getData)
		} else {
			t.Log("get all data by user and count < 0")
		}
	}
	err = DeleteByID(&ArgsDeleteByID{
		ID:    newData.ID,
		OrgID: newData.OrgID,
	})
	if err != nil {
		t.Error(err, ", newData.ID: ", newData.ID)
	}
}

func TestGetByID(t *testing.T) {
	newData, err := Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试C",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   23,
					EndMinute: 59,
				},
			},
			AllowHoliday: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new data: ", newData)
	}
	//尝试获取该数据
	getData, err := GetOne(&ArgsGetOne{
		ID:    newData.ID,
		OrgID: 0,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("get data by id: ", getData)
	}
	_, err = GetOne(&ArgsGetOne{
		ID:    987666126371823,
		OrgID: 0,
	})
	if err == nil {
		t.Error("get not exist data by id")
	}
	err = DeleteByID(&ArgsDeleteByID{
		ID:    newData.ID,
		OrgID: newData.OrgID,
	})
	if err != nil {
		t.Error(err, ", newData.ID: ", newData.ID)
	}
}

func TestGetByIDAndUserID(t *testing.T) {
	TestTimeInit(t)
	newData, err := Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试D",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   23,
					EndMinute: 59,
				},
			},
			AllowHoliday: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new data: ", newData)
	}
	//尝试获取该数据
	getData, err := GetOne(&ArgsGetOne{
		ID:    newData.ID,
		OrgID: defaultOrgID,
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("get data by id and user id: ", getData)
	}
	getData2, err := GetOne(&ArgsGetOne{
		ID:    newData.ID,
		OrgID: 987712381312,
	})
	if err == nil {
		t.Error("get not exist data by id and user id, ", getData2)
	}
	err = DeleteByID(&ArgsDeleteByID{
		ID:    newData.ID,
		OrgID: newData.OrgID,
	})
	if err != nil {
		t.Error(err, ", newData.ID: ", newData.ID)
	}
}

func TestUpdateByID(t *testing.T) {
	TestTimeInit(t)
	newData, err := Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试E",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{},
			Week:      []int{1, 2, 3, 4, 5, 6},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   23,
					EndMinute: 59,
				},
			},
			AllowHoliday: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new data: ", newData)
	}
	//修改数据
	err = UpdateByID(&ArgsUpdateByID{
		ID:       newData.ID,
		OrgID:    defaultOrgID,
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试E2",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{},
			Week:      []int{3, 4, 5, 6},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   23,
					EndMinute: 59,
				},
			},
			AllowHoliday: false,
		},
	},
	)
	if err != nil {
		t.Error(err)
	}
	getData, err := GetOne(&ArgsGetOne{
		ID:    newData.ID,
		OrgID: 0,
	})
	if err != nil {
		t.Error(err)
	} else {
		if getData.OrgID != defaultOrgID || getData.Name != "上班测试E2" {
			t.Error("update data is error, ", getData)
		} else {
			t.Log("update data: ", getData)
		}
	}
	err = DeleteByID(&ArgsDeleteByID{
		ID:    newData.ID,
		OrgID: newData.OrgID,
	})
	if err != nil {
		t.Error(err, ", newData.ID: ", newData.ID)
	}
}

func TestUpdateByIDAndUserID(t *testing.T) {
	newData, err := Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试F",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{1, 2, 3, 4, 5, 6},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   23,
					EndMinute: 59,
				},
			},
			AllowHoliday: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new data: ", newData)
	}
	//修改数据
	err = UpdateByID(&ArgsUpdateByID{
		ID:       newData.ID,
		OrgID:    defaultOrgID,
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试F2",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{3, 4, 5, 6},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   23,
					EndMinute: 59,
				},
			},
			AllowHoliday: false,
		},
	})
	if err != nil {
		t.Error(err)
	}
	getData, err := GetOne(&ArgsGetOne{
		ID:    newData.ID,
		OrgID: 0,
	})
	if err != nil {
		t.Error(err)
	} else {
		if getData.OrgID != defaultOrgID || getData.Name != "上班测试F2" {
			t.Error("update data is error, ", getData)
		} else {
			t.Log("update data: ", getData)
		}
	}
	//修改数据2
	var newID int64 = 977712381312
	err = UpdateByID(&ArgsUpdateByID{
		ID:       getData.ID,
		OrgID:    newID,
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试F3",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{3, 4, 5},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   23,
					EndMinute: 59,
				},
			},
			AllowHoliday: false,
		},
	})
	if err == nil {
		t.Error("update is error ", newData.ID, ", OrgID: ", newData.OrgID, ", but select OrganizationalID: ", newID)
		getData, err := GetOne(&ArgsGetOne{
			ID:    newData.ID,
			OrgID: 0,
		})
		if err != nil {
			t.Error(err)
		} else {
			t.Log("update result: ", getData)
		}
	}
	err = DeleteByID(&ArgsDeleteByID{
		ID:    newData.ID,
		OrgID: newData.OrgID,
	})
	if err != nil {
		t.Error(err, ", newData.ID: ", newData.ID)
	}
}

// 检查任务
func TestCheckIsWorkByData(t *testing.T) {
	TestTimeInit(t)
	newData, err := Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{1},
		Binds:    []int64{2},
		Name:     "上班测试J",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{},
			Week:      []int{1, 2, 3, 4, 5, 6, 7, 0},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   23,
					EndMinute: 59,
				},
			},
			AllowHoliday: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new data: ", newData)
	}
	b := checkIsWorkNowDayByData(&newData)
	if !b {
		t.Error("check is error")
	} else {
		t.Log("is work, v: ", newData.ID)
	}
	time.Sleep(time.Second * 3)
	err = DeleteByID(&ArgsDeleteByID{
		ID:    newData.ID,
		OrgID: newData.OrgID,
	})
	if err != nil {
		t.Error(err, ", newData.ID: ", newData.ID)
	}
	newData, err = Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试J2",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{1, 3, 0},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   0,
					EndMinute: 5,
				},
			},
			AllowHoliday: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new data: ", newData)
	}
	b = checkIsWorkByData(&newData)
	if b {
		t.Error("check is error")
	}
	err = DeleteByID(&ArgsDeleteByID{
		ID:    newData.ID,
		OrgID: newData.OrgID,
	})
	if err != nil {
		t.Error(err, ", newData.ID: ", newData.ID)
	}
}

func TestCheckIsWorkByData2(t *testing.T) {
	var (
		b = checkIsWorkByData(&FieldsWorkTime{
			ID:       0,
			CreateAt: time.Time{},
			UpdateAt: time.Time{},
			ExpireAt: time.Time{},
			OrgID:    0,
			Groups:   nil,
			Binds:    nil,
			Name:     "",
			IsWork:   false,
			Configs: FieldsConfigs{
				Month:        []int{},
				MonthDay:     []int{},
				MonthWeek:    []int{},
				Week:         []int{},
				WorkTime:     []FieldsWorkTimeTime{},
				AllowHoliday: true,
			},
			RotConfig: FieldsConfigRot{
				NowKey:  0,
				DiffDay: 1,
				WorkTime: []FieldsWorkTimeTime{

					{
						StartHour:   22,
						StartMinute: 0,
						EndHour:     17,
						EndMinute:   0,
					},
					{
						StartHour:   16,
						StartMinute: 0,
						EndHour:     17,
						EndMinute:   0,
					},
				},
			},
		})
	)
	if !b {
		t.Error("check is error")
	}
}

func TestCheckIsWorkNowDayByData(t *testing.T) {
	var newData, err = Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试I",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{},
			Week:      []int{},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   23,
					EndMinute: 59,
				},
			},
			AllowHoliday: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new data: ", newData)
	}
	b := checkIsWorkNowDayByData(&newData)
	if !b {
		t.Error("check is error")
	}
	err = DeleteByID(&ArgsDeleteByID{
		ID:    newData.ID,
		OrgID: newData.OrgID,
	})
	if err != nil {
		t.Error(err, ", newData.ID: ", newData.ID)
	}
	//检查不在当前时间的情况怎么处理
	newData, err = Create(
		&ArgsCreate{
			ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
			OrgID:    defaultOrgID,
			Groups:   []int64{},
			Binds:    []int64{},
			Name:     "上班测试I2",
			Configs: FieldsConfigs{
				Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
				MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
				MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
				Week:      []int{carbon.Yesterday().DayOfWeek()},
				WorkTime: []FieldsWorkTimeTime{
					{
						StartHour:   0,
						StartMinute: 0,

						EndHour:   23,
						EndMinute: 59,
					},
				},
				AllowHoliday: false,
			},
		},
	)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new data: ", newData)
	}
	b = checkIsWorkNowDayByData(&newData)
	if b {
		t.Error("check is error, week: ", newData.Configs.Week)
	}
	err = DeleteByID(&ArgsDeleteByID{
		ID:    newData.ID,
		OrgID: newData.OrgID,
	})
	if err != nil {
		t.Error(err, ", newData.ID: ", newData.ID)
	}
	//检查只有今天
	newData, err = Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试I3",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{2, 6, 7, CoreFilter.GetNowTimeCarbon().DayOfWeek()},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   23,
					EndMinute: 59,
				},
			},
			AllowHoliday: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new data: ", newData)
	}
	b = checkIsWorkNowDayByData(&newData)
	if !b {
		t.Error("check is error")
	}
	err = DeleteByID(&ArgsDeleteByID{
		ID:    newData.ID,
		OrgID: newData.OrgID,
	})
	if err != nil {
		t.Error(err, ", newData.ID: ", newData.ID)
	}
	//请注意，需修改为当天进行判断处理，否则将失效
	newData, err = Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试I3",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   23,
					EndMinute: 59,
				},
			},
			AllowHoliday: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new data: ", newData)
	}
	b = checkIsWorkNowDayByData(&newData)
	if !b {
		t.Error("check is error")
	}
	err = DeleteByID(&ArgsDeleteByID{
		ID:    newData.ID,
		OrgID: newData.OrgID,
	})
	if err != nil {
		t.Error(err, ", newData.ID: ", newData.ID)
	}
}

func TestCheckIsWorkByID(t *testing.T) {
	newData, err := Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试K",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{1, 2, 3, 4, 5, 6, 7, CoreFilter.GetNowTimeCarbon().DayOfWeek()},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   23,
					EndMinute: 59,
				},
			},
			AllowHoliday: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new data: ", newData)
	}
	time.Sleep(time.Second * 2)
	b := checkIsWorkByID(&ArgsCheckIsWorkByID{
		ID: newData.ID,
	})
	if !b {
		t.Error("check is error")
	}
	err = DeleteByID(&ArgsDeleteByID{
		ID:    newData.ID,
		OrgID: newData.OrgID,
	})
	if err != nil {
		t.Error(err, ", newData.ID: ", newData.ID)
	}
	newData, err = Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试K2",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{1, 2, 3, 4, 5, 6, 7},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   0,
					EndMinute: 5,
				},
			},
			AllowHoliday: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new data: ", newData)
	}
	b = checkIsWorkByID(&ArgsCheckIsWorkByID{
		ID: newData.ID,
	})
	if b {
		t.Error("check is error")
	}
	err = DeleteByID(&ArgsDeleteByID{
		ID:    newData.ID,
		OrgID: newData.OrgID,
	})
	if err != nil {
		t.Error(err, ", newData.ID: ", newData.ID)
	}
}

func TestGetWorkDayInMonth(t *testing.T) {
	newData, err := Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试L",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{1, 2, 3, 4, 5, 6, 7},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   23,
					EndMinute: 59,
				},
			},
			AllowHoliday: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new data: ", newData)
	}
	count := getWorkDayInMonth(&newData)
	if count != 7 {
		CoreLog.Error("count is error")
	}
	err = DeleteByID(&ArgsDeleteByID{
		ID:    newData.ID,
		OrgID: newData.OrgID,
	})
	if err != nil {
		t.Error(err, ", newData.ID: ", newData.ID)
	}
}

// 检查某个组是否在工作中？
func TestCheckIsWorkByGroupOrBind(t *testing.T) {
	newData, err := Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{123},
		Binds:    []int64{},
		Name:     "上班测试L",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{1, 2, 3, 4, 5, 6, 7, CoreFilter.GetNowTimeCarbon().DayOfWeek()},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   23,
					EndMinute: 59,
				},
			},
			AllowHoliday: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new data: ", newData)
	}
	time.Sleep(time.Second * 1)
	err = DeleteByID(&ArgsDeleteByID{
		ID:    newData.ID,
		OrgID: newData.OrgID,
	})
	if err != nil {
		t.Error(err, ", newData.ID: ", newData.ID)
	}
}

func TestDeleteByID(t *testing.T) {
	newData, err := Create(&ArgsCreate{
		ExpireAt: CoreFilter.GetNowTime().Add(time.Hour * 1),
		OrgID:    defaultOrgID,
		Groups:   []int64{},
		Binds:    []int64{},
		Name:     "上班测试G",
		Configs: FieldsConfigs{
			Month:     []int{CoreFilter.GetNowTimeCarbon().MonthOfYear()},
			MonthDay:  []int{CoreFilter.GetNowTimeCarbon().DayOfMonth()},
			MonthWeek: []int{CoreFilter.GetNowTimeCarbon().WeekOfMonth() - 1},
			Week:      []int{1, 2, 3, 4, 5, 6, 7},
			WorkTime: []FieldsWorkTimeTime{
				{
					StartHour:   0,
					StartMinute: 0,

					EndHour:   23,
					EndMinute: 59,
				},
			},
			AllowHoliday: false,
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("new data: ", newData)
	}
	err = DeleteByID(&ArgsDeleteByID{
		ID:    newData.ID,
		OrgID: newData.OrgID,
	})
	if err != nil {
		t.Error(err, ", newData.ID: ", newData.ID)
	}
}

// 清理本表所有数据，完成最终测试
func TestDeleteAllData(t *testing.T) {
}
