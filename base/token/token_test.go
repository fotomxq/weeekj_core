package BaseToken

import (
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQLFrom "gitee.com/weeekj/weeekj_core/v5/core/sql/from"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	ToolsTest "gitee.com/weeekj/weeekj_core/v5/tools/test"
	"testing"
	"time"
)

var (
	takeData FieldsTokenType
)

func TestInit(t *testing.T) {
	ToolsTest.Init(t)
}

func TestExpire(t *testing.T) {
	//建立过期机制
	//建立新的数据
	newData, errCode, err := Create(&ArgsCreate{
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     606040,
			Mark:   "",
			Name:   "17345628912",
		},
		LoginInfo: CoreSQLFrom.FieldsFrom{
			System: "phone",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
		LimitKeyLen: 16,
		IP:          "0.0.0.0",
		ExpireAt:    CoreFilter.GetNowTime().Add(time.Second * 2).Format(time.RFC3339Nano),
		IsRemember:  false,
	})
	if err != nil {
		t.Error("cannot create expire data, ", errCode, ", ", err)
	} else {
		t.Log("create new to expire data: ", newData)
	}
	//等待1.5秒
	time.Sleep(time.Millisecond * 1500)
	//强制过期处理
	getData, err := GetByFromAndLogin(&ArgsGetByFromAndLogin{
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     606040,
			Mark:   "",
			Name:   "17345628912",
		}, LoginInfo: CoreSQLFrom.FieldsFrom{
			System: "phone",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
	})
	if err != nil {
		t.Error("expire data is not exist, ", "now time: ", CoreFilter.GetNowTime())
	}
	//等待1秒
	time.Sleep(time.Millisecond * 1000)
	//强制过期处理
	//检查该数据是否还存在？
	getData, err = GetByFromAndLogin(&ArgsGetByFromAndLogin{
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     606040,
			Mark:   "",
			Name:   "17345628912",
		}, LoginInfo: CoreSQLFrom.FieldsFrom{
			System: "phone",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
	})
	if err == nil {
		t.Error("expire data is exist, ", "now time: ", CoreFilter.GetNowTime().String(), ", data id and expire time: ", getData.ID, ", ", getData.ExpireAt)
	}
	//重新创建，等待2秒更新，等待1秒测试更新
	newData, errCode, err = Create(&ArgsCreate{
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     606040,
			Mark:   "",
			Name:   "17345628912",
		},
		LoginInfo: CoreSQLFrom.FieldsFrom{
			System: "phone",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
		LimitKeyLen: 16,
		IP:          "0.0.0.0",
		ExpireAt:    CoreFilter.GetNowTime().Add(time.Second * 30).Format(time.RFC3339Nano),
	})
	if err != nil {
		t.Error("cannot create expire data, ", errCode, newData)
	} else {
		t.Log("create new to expire data to wait update: ", newData.ID, ", time: ", newData.CreateAt, ", ", newData.ExpireAt)
	}
	//更新数据过期时间，之后等待1秒测试是否存在，如果不存在则说明更新没生效
	if err == nil {
		//等待2秒
		time.Sleep(time.Millisecond * 2000)
		err = UpdateExpire(&ArgsUpdateExpire{
			ID:       newData.ID,
			ExpireAt: CoreFilter.GetNowTime().Add(time.Minute * 30).Format(time.RFC3339Nano),
		})
		if err != nil {
			t.Error("cannot update expire data, ", newData)
		} else {
			t.Log("update data time: ", CoreFilter.GetNowTime())
		}
		//等待1秒
		time.Sleep(time.Millisecond * 1000)
		//检查该数据是否还存在？
		newData, err = GetByID(&ArgsGetByID{
			ID: newData.ID,
		})
		if err != nil {
			t.Error("expire data is not exist, ", "now time: ", CoreFilter.GetNowTime().String(), ", err: ", err.Error())
		} else {
			t.Log("update expire data and get data: ", newData.ID, ", time: ", newData.CreateAt, ", ", newData.ExpireAt)
		}
	}
	//再等待1秒，使用key更新时间
	time.Sleep(time.Millisecond * 1000)
	newData, errCode, err = Create(&ArgsCreate{
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     606040,
			Mark:   "",
			Name:   "17345628912",
		},
		LoginInfo: CoreSQLFrom.FieldsFrom{
			System: "phone",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
		LimitKeyLen: 16,
		IP:          "0.0.0.0",
		ExpireAt:    CoreFilter.GetNowTime().Add(time.Minute * 3).Format(time.RFC3339Nano),
	})
	if err != nil {
		t.Error("cannot create expire data by key, ", errCode, newData)
	} else {
		t.Log("create new to expire data to wait update by key: ", newData.ID, ", time: ", newData.CreateAt, ", ", newData.ExpireAt)
	}
	//更新数据过期时间，之后等待1秒测试是否存在，如果不存在则说明更新没生效
	if err == nil {
		//等待2秒
		time.Sleep(time.Millisecond * 2000)
		err = UpdateExpireByKey(&ArgsUpdateExpireByKey{
			Key:      newData.Key,
			ExpireAt: CoreFilter.GetNowTime().Add(time.Second * 3).Format(time.RFC3339Nano),
		})
		if err != nil {
			t.Error("cannot update expire data by key, ", newData)
		} else {
			t.Log("update data time by key: ", CoreFilter.GetNowTime())
		}
		//等待1秒
		time.Sleep(time.Millisecond * 1000)
		//检查该数据是否还存在？
		newData, err = GetByKey(&ArgsGetByKey{
			Key: newData.Key,
		})
		if err != nil {
			t.Error("expire data is not exist by key, ", "now time: ", CoreFilter.GetNowTime().String(), ", err: ", err.Error())
		} else {
			t.Log("update expire data and get data by key: ", newData.ID, ", time: ", newData.CreateAt, ", ", newData.ExpireAt)
		}
	}
}

func TestCreate(t *testing.T) {
	//创建一个
	newData, errCode, err := Create(&ArgsCreate{
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     606050,
			Mark:   "",
			Name:   "17345628912",
		},
		LoginInfo: CoreSQLFrom.FieldsFrom{
			System: "phone",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
		LimitKeyLen: 16,
		IP:          "0.0.0.0",
		ExpireAt:    CoreFilter.GetNowTime().Add(time.Minute * 3).Format(time.RFC3339Nano),
	})
	if err != nil {
		t.Error("cannot create first data, ", errCode, newData)
	} else {
		t.Log("create new: ", newData)
	}
	//重复创建测试
	takeData, errCode, err = Create(&ArgsCreate{
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     606050,
			Mark:   "",
			Name:   "17345628912",
		},
		LoginInfo: CoreSQLFrom.FieldsFrom{
			System: "phone",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
		LimitKeyLen: 16,
		IP:          "0.0.0.0",
		ExpireAt:    CoreFilter.GetNowTime().Add(time.Minute * 3).Format(time.RFC3339Nano),
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log("create new replace: ", takeData)
	}
	//指定删除测试
	err = DeleteByID(&ArgsDeleteByID{
		ID: takeData.ID,
	})
	if err != nil {
		t.Error("cannot delete by id: ", err)
	}
	//重建后多次删除测试
	takeData, errCode, err = Create(&ArgsCreate{
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     606052,
			Mark:   "",
			Name:   "17345628912",
		},
		LoginInfo: CoreSQLFrom.FieldsFrom{
			System: "phone",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
		LimitKeyLen: 16,
		IP:          "0.0.0.0",
		ExpireAt:    CoreFilter.GetNowTime().Add(time.Minute * 3).Format(time.RFC3339Nano),
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log("create new replace2: ", takeData)
	}
	t.Log(map[string]string{"FromSystem": "user", "FromID": "606052", "LoginSystem": "phone"})
	err = DeleteByFrom(&ArgsDeleteByFrom{
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     606052,
			Mark:   "",
			Name:   "17345628912",
		}, LoginInfo: CoreSQLFrom.FieldsFrom{
			System: "phone",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
	})
	if err != nil {
		t.Error("cannot delete by from: ", err)
	}
	takeData, errCode, err = Create(&ArgsCreate{
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     606050,
			Mark:   "",
			Name:   "17345628912",
		},
		LoginInfo: CoreSQLFrom.FieldsFrom{
			System: "phone",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
		LimitKeyLen: 16,
		IP:          "0.0.0.0",
		ExpireAt:    CoreFilter.GetNowTime().Add(time.Minute * 3).Format(time.RFC3339Nano),
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log("create new replace3: ", takeData)
	}
	//批量创建，重叠测试
	beginInt := 22345678912
	startInt := 0
	maxInt := 10
	for {
		_, errCode, err = Create(&ArgsCreate{
			FromInfo: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     606050,
				Mark:   "",
				Name:   "17345628912",
			},
			LoginInfo: CoreSQLFrom.FieldsFrom{
				System: "phone",
				ID:     0,
				Mark:   "",
				Name:   "",
			},
			LimitKeyLen: 16,
			IP:          "0.0.0.0",
			ExpireAt:    CoreFilter.GetNowTime().Add(time.Minute * 3).Format(time.RFC3339Nano),
		})
		if err != nil {
			t.Error(errCode, err)
		} else {
			t.Log("create new: ", newData)
		}
		beginInt += 1
		startInt += 1
		if startInt >= maxInt {
			break
		}
	}
	//批量添加，后续列表测试
	beginInt = 32345678912
	var startInt2 = 0
	var maxInt2 = 10
	for {
		_, errCode, err = Create(&ArgsCreate{
			FromInfo: CoreSQLFrom.FieldsFrom{
				System: "user",
				ID:     606050,
				Mark:   "",
				Name:   "17345628912",
			},
			LoginInfo: CoreSQLFrom.FieldsFrom{
				System: "phone",
				ID:     0,
				Mark:   "",
				Name:   "",
			},
			LimitKeyLen: 16,
			IP:          "0.0.0.0",
			ExpireAt:    CoreFilter.GetNowTime().Add(time.Minute * 3).Format(time.RFC3339Nano),
		})
		if err != nil {
			t.Error(errCode, err)
		} else {
			t.Log("create new: ", newData)
		}
		beginInt += 1
		startInt2 += 1
		if startInt2 >= maxInt2 {
			break
		}
	}
}

func TestGet(t *testing.T) {
	dataList, count, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		},
	})
	if err != nil {
		t.Error("cannot get list, ", err)
	} else {
		t.Log("get data list, count: ", count, ", data: ", dataList)
	}
	dataList, count, err = GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		}, FromInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
	})
	if err != nil {
		t.Error("cannot get list by from, ", err)
	} else {
		t.Log("get data list, count by from: ", count, ", data count: ", len(dataList))
	}
	dataList, count, err = GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		}, FromInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     takeData.ID,
			Mark:   "",
			Name:   "",
		},
		Search: "",
	})
	if err != nil {
		t.Error("cannot get list by from id, ", err)
	} else {
		t.Log("get data list, count by from id: ", count, ", data count: ", len(dataList))
	}
	dataList, count, err = GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		}, FromInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     takeData.ID,
			Mark:   "",
			Name:   "",
		}, Search: "12345678912",
	})
	if err != nil {
		t.Error("cannot get list by search, ", err)
	} else {
		t.Log("get data list, count by search: ", count, ", data count: ", len(dataList))
	}
	data, err := GetByFromAndLogin(&ArgsGetByFromAndLogin{
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     1,
			Mark:   "",
			Name:   "",
		}, LoginInfo: CoreSQLFrom.FieldsFrom{
			System: "phone",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
	})
	if err == nil {
		t.Error("get token and this is now exist data, ", data)
	} else {
		t.Log("success get data and test failed, ", err)
	}
}

func TestDelete(t *testing.T) {
	//将所有过期数据踢下线
	time.Sleep(time.Second * 3)
	DeleteAll()
	dataList, count, err := GetList(&ArgsGetList{
		Pages: CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  10,
			Sort: "id",
			Desc: false,
		}, FromInfo: CoreSQLFrom.FieldsFrom{
			System: "",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
	})
	if err != nil {
		t.Error(err)
	} else {
		if count > 0 {
			t.Error("delete all expire data but not run, ", dataList)
		}
	}
	//重建数据
	var errCode string
	takeData, errCode, err = Create(&ArgsCreate{
		FromInfo: CoreSQLFrom.FieldsFrom{
			System: "user",
			ID:     606050,
			Mark:   "",
			Name:   "17345628912",
		},
		LoginInfo: CoreSQLFrom.FieldsFrom{
			System: "phone",
			ID:     0,
			Mark:   "",
			Name:   "",
		},
		LimitKeyLen: 16,
		IP:          "0.0.0.0",
		ExpireAt:    CoreFilter.GetNowTime().Add(time.Minute * 3).Format(time.RFC3339Nano),
	})
	if err != nil {
		t.Error(errCode, err)
	} else {
		t.Log("create new 4: ", takeData.ID)
	}
	//强制踢下线某个
	err = DeleteByID(&ArgsDeleteByID{
		ID: takeData.ID,
	})
	if err != nil {
		t.Error("cannot delete id, ", err.Error())
	} else {
		t.Log("delete data by id: ", takeData.ID)
	}
}

func TestDeleteAll(t *testing.T) {
	DeleteAll()
}
